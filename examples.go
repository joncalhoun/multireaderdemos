package multireaderdemos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func mustOpen(filename string) io.Reader {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func Copy() {
	file, err := pretendToOpenFile("demo.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, file)
}

func pretendToOpenFile(filename string) (io.Reader, error) {
	return strings.NewReader("This is a fake file for demo purposes"), nil
}

type Log struct {
	Level string    `json:"level"`
	Date  time.Time `json:"date"`
	Msg   string    `json:"msg"`
}

func JSONLogsWithoutMultiReader() error {
	// Open the files
	filenames := []string{"day1.log", "day2.log", "day3.log"}
	var files []io.Reader
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	// Read each in sequence
	var logs []Log
	for _, file := range files {
		dec := json.NewDecoder(file)
		for dec.More() {
			var log Log
			err := dec.Decode(&log)
			if err != nil {
				return err
			}
			logs = append(logs, log)
		}
	}

	// Now we can work with all the logs.
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}

func JSONLogsWithMultiReader() error {
	// Open the files
	filenames := []string{"day1.log", "day2.log", "day3.log"}
	var files []io.Reader
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	// Read each in sequence
	var logs []Log
	logReader := io.MultiReader(files...)
	dec := json.NewDecoder(logReader)
	for dec.More() {
		var log Log
		err := dec.Decode(&log)
		if err != nil {
			return err
		}
		logs = append(logs, log)
	}

	// Now we can work with all the logs.
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}

func checkContentTypeV1(r io.ReadSeeker) error {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(buf[:n])
	if contentType != "image/png" {
		return fmt.Errorf("unexpected content type: %s", contentType)
	}
	r.Seek(0, 0)
	return nil
}

func CreateImageV1(r io.ReadSeeker, filename string) error {
	err := checkContentTypeV1(r)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return nil
}

func checkContentType(r io.Reader) ([]byte, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buf[:n])
	if contentType != "image/png" {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}
	return buf[:n], nil
}

func CreateImage(r io.Reader, filename string) error {
	readBytes, err := checkContentType(r)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r = io.MultiReader(bytes.NewReader(readBytes), r)
	_, err = io.Copy(f, r)
	return nil
}
