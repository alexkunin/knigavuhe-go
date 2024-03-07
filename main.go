package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"knigavuhe/book_parser"
	"log"
	"net/http"
	"time"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func getPageContent(url string) (error, string) {
	myClient := &http.Client{Timeout: 10 * time.Second}

	res, err := myClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "text/html; charset=UTF-8" {
		return errors.New(fmt.Sprintf("Content-Type: \"%s\" is not supported", contentType)), ""
	}

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)), ""
	}
	if err != nil {
		return err, ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	return nil, string(body)
}

func main() {
	bookUrl := "https://knigavuhe.org/book/barliona/"
	fmt.Println("Book URL:", bookUrl)

	err, body := getPageContent(bookUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err, chapters := book_parser.ExtractBookChapters(body); err == nil {
		fmt.Println(prettyPrint(chapters))
	} else {
		log.Fatal(err)
	}

	bookInfo := book_parser.ExtractBookInfo(body)
	fmt.Println(prettyPrint(bookInfo))
}
