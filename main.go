package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/yosssi/gohtml"
)

func main() {
	var content string = fetchContent("https://www.google.com")
	fmt.Println(content)
	putContentFile(content, "./buffer.html")
}

func putContentFile(content string, filename string) {

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()
}

func fetchContent(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyText = gohtml.FormatBytes(bodyText)

	return string(bodyText)
}
