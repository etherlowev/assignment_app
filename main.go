package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type body struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

func checkError(e error) {
	if e != nil {
		log.Fatal("Failed to create document: ", e)
	}
}

func createDocument(num uint64, total uint64, wg *sync.WaitGroup, counter *atomic.Uint64) {
	reqBody := body{Author: "tools", Title: strconv.FormatUint(num, 10)}
	jsonBody, err := json.Marshal(reqBody)
	checkError(err)

	response, err := http.Post("http://127.0.0.1:8080/api/document/create",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	checkError(err)
	checkError(response.Body.Close())
	fmt.Printf("Created document %v/%v\n", counter.Add(1), total)
	wg.Done()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file-path>")
		os.Exit(1)
	}
	path := os.Args[1]

	file, err := os.Open(path)
	checkError(err)

	defer func(file *os.File) {
		checkError(file.Close())
	}(file)

	buf := bufio.NewReader(file)

	line, _, err := buf.ReadLine()
	checkError(err)

	amountOfFiles, err := strconv.ParseUint(string(line), 10, 64)
	checkError(err)

	var wg sync.WaitGroup
	var counter atomic.Uint64
	for i := range amountOfFiles {
		wg.Add(1)
		go createDocument(i+1, amountOfFiles, &wg, &counter)
	}

	wg.Wait()

	os.Exit(0)
}
