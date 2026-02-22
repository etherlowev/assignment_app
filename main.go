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
)

type body struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
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

	for i := range amountOfFiles {
		fmt.Printf("Creating %v/%v\n", i+1, amountOfFiles)
		reqBody := body{Author: "tools", Title: strconv.FormatUint(i+1, 10)}
		jsonBody, err := json.Marshal(reqBody)
		checkError(err)

		response, err := http.Post("http://127.0.0.1:8080/api/document/create",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)

		checkError(err)
		checkError(response.Body.Close())
	}
	os.Exit(0)
}
