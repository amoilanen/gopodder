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

func (r *HttpClient) DownloadFile(url, outputPath string) error {
	var err error
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	currentOutputSize, err := getCurrentOutputSize(outputPath)
	existsNonEmptyOutput := currentOutputSize > 0 && err == nil

	var outFile *os.File
	if existsNonEmptyOutput {
		outFile, err = os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		outFile, err = os.Create(outputPath)
	}
	if err != nil {
		return err
	}
	defer outFile.Close()

	if existsNonEmptyOutput {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentOutputSize))
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if existsNonEmptyOutput {
		if response.StatusCode == http.StatusPartialContent {
			fmt.Printf("Resuming download of file %s, current size %d\n", outputPath, currentOutputSize)
		} else if response.StatusCode == http.StatusOK {
			fmt.Printf("Already fully downloaded %s\n", outputPath)
		} else {
			fmt.Printf("Removing partially downloaded file = %s since server does not support resuming downloads\n", outputPath)
			err := os.Remove(outputPath)
			if err != nil {
				return r.DownloadFile(url, outputPath)
			}
		}
	} else {
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
		}
	}
	contentLength := response.ContentLength

	if contentLength > currentOutputSize {
		progressBar := NewDownloadProgressBar(contentLength)
		progressBar.OnProgress(currentOutputSize)
		writer := io.MultiWriter(outFile, progressBar)
		_, err = io.Copy(writer, response.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
