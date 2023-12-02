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
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	contentLength := response.ContentLength

	progressBar := NewDownloadProgressBar(contentLength)
	writer := io.MultiWriter(outFile, progressBar)

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return err
	}
	return nil
}
