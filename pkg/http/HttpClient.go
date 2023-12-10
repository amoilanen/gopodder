package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type HttpClient struct {
}

func (r *HttpClient) GetBytesFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (r *HttpClient) DownloadFile(url, outputPath string) error {
	currentOutputSize, err := getCurrentOutputSize(outputPath)
	existsNonEmptyOutput := currentOutputSize > 0 && err == nil

	if existsNonEmptyOutput {
		return resumeDownloadToFile(url, outputPath, currentOutputSize)
	} else {
		return downloadToFileFromScratch(url, outputPath)
	}
}

func getCurrentOutputSize(outputPath string) (int64, error) {
	if fileInfo, err := os.Stat(outputPath); err == nil {
		size := fileInfo.Size()
		return size, nil
	} else if os.IsNotExist(err) {
		return 0, nil
	} else {
		return 0, fmt.Errorf("Error trying to check if the file %s is empty %v", outputPath, err)
	}
}

func resumeDownloadToFile(url, outputPath string, currentOutputSize int64) error {
	outFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentOutputSize))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	needToDownload := true
	contentLength := response.ContentLength
	fullDownloadSize := contentLength
	if response.StatusCode == http.StatusPartialContent {
		fmt.Printf("Resuming download of file %s, current size %d\n", outputPath, currentOutputSize)
		fullDownloadSize = contentLength + currentOutputSize
	} else if response.StatusCode == http.StatusOK {
		fmt.Printf("Already fully downloaded %s\n", outputPath)
		needToDownload = false
	} else {
		fmt.Printf("Removing partially downloaded file = %s since server does not support resuming downloads\n", outputPath)
		err := os.Remove(outputPath)
		if err != nil {
			return downloadToFileFromScratch(url, outputPath)
		}
	}

	if needToDownload {
		return downloadWithProgress(fullDownloadSize, currentOutputSize, response, outFile)
	}
	return nil
}

func downloadToFileFromScratch(url, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}
	return downloadWithProgress(response.ContentLength, 0, response, outFile)
}

func downloadWithProgress(downloadSize int64, initialProgress int64, response *http.Response, output *os.File) error {
	progressBar := NewDownloadProgressBar(downloadSize)
	progressBar.OnProgress(initialProgress)
	writer := io.MultiWriter(output, progressBar)
	_, err := io.Copy(writer, response.Body)
	if err != nil {
		return err
	}
	return nil
}
