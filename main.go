package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

func main() {
	var targetURL string = os.Args[1]
	var content string = fetchContent(targetURL)
	//fmt.Println(content)
	putContentFile(content, "./buffer.html")
	buffer, err := readBuffer("./buffer.html")

	if err != nil {
		log.Fatal(err)
	}
	parseHTML(buffer)
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
	fmt.Printf("wrote %d bytes\n\n", n)
	fmt.Printf("end fetching content : %s \n\n", filename)
	f.Sync()
}

func readBuffer(fileName string) (string, error) {

	result, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func parseHTML(content string) (data []string) {
	reader := strings.NewReader(content)
	tokenizer := html.NewTokenizer(reader)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return
			}
			fmt.Printf("Error: %v", tokenizer.Err())
			return
		}

		tag, hasAttr := tokenizer.TagName()

		if string(tag) == "" {
			continue
		}

		fmt.Printf("Tag: %v", string(tag))

		var attributes string = ""

		if hasAttr {
			for {
				attrKey, attrValue, moreAttr := tokenizer.TagAttr()
				// if string(attrKey) == "" {
				//     break
				// }
				// fmt.Printf("Attr: %v\n", string(attrKey))
				// fmt.Printf("Attr: %v\n", string(attrValue))
				// fmt.Printf("Attr: %v\n", moreAttr)

				attributes += " " + string(attrKey) + "=" + string(attrValue) + ";"

				if !moreAttr {
					break
				}
			}

			fmt.Print(" Attr:" + attributes)

		}
		fmt.Print("\n")
	}
}

func fetchContent(url string) string {

	fmt.Printf("start fetching content : %s \n\n", url)

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
