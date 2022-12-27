package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	directory := "./storage"
	extension := "txt"
	url := "https://jsonplaceholder.typicode.com/posts"
	method := "GET"
	client := &http.Client{}
	count := 100

	var wg sync.WaitGroup

	clearDirectory(directory)

	for i := 1; i <= count; i++ {
		wg.Add(1)
		url := fmt.Sprintf("%s/%d", url, i)
		path := fmt.Sprintf("%s/%d.%s", directory, i, extension)
		go func() {
			defer wg.Done()
			data, err := getData(url, method, client)
			if err == nil {
				writeData(path, data)
			}
		}()
	}
	wg.Wait()
}

func writeData(path string, content []byte) {
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		log.Fatal(err)
	}
}

func writeDataWithBuff(path string, content []byte) {
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Sync()
	writer := bufio.NewWriter(file)

	_, err = writer.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}

func clearDirectory(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		pathToFile := fmt.Sprintf("%s/%s", directory, file.Name())
		deleteFile(pathToFile)
	}
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}

func getData(url string, method string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	var emptyArr []byte

	if err != nil {
		return emptyArr, err
	}
	res, err := client.Do(req)

	if err != nil {
		return emptyArr, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyArr, err
	}
	return body, nil
}
