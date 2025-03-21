package main

import (
	"fmt"
)

func main() {
	htmlData, err := getHTMLByPath("ex3.html")
	if err != nil {
		panic(err)
	}

	links, err := ParseLinkFromHtml(htmlData)
	if err != nil {
		panic(err)
	}

	for _, link := range links {
		fmt.Println(link)
	}

	htmlData, err = getHTMLByURL("https://pkg.go.dev/golang.org/x/net@v0.37.0/html")
	if err != nil {
		panic(err)
	}

	links, err = ParseLinkFromHtml(htmlData)
	if err != nil {
		panic(err)
	}

	for _, link := range links {
		fmt.Println(link)
	}

}
