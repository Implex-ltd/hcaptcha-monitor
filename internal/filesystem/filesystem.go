package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

const BasePath = "../assets"

func ReadFile(filePath string) ([]string, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", BasePath, filePath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func AppendFile(filePath, line string) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", BasePath, filePath), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, line)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(url string, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s/%s", BasePath, filePath))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
