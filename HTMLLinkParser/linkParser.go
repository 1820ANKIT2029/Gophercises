package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href, Text string
}

type Links []Link

func (l Link) String() string {
	return fmt.Sprintf("Link{\n\tHref: \"%s\",\n\tText: \"%s\",\n}", l.Href, l.Text)
}

func ParseLinkFromHtml(HTMLData string) (Links, error) {
	node, err := html.Parse(strings.NewReader(HTMLData))
	if err != nil {
		return nil, err
	}

	var links Links

	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			var link Link

			for _, a := range n.Attr {
				if a.Key == "href" {
					link.Href = a.Val
					break
				}
			}

			var text string = ""

			for child := range n.ChildNodes() {
				if child.Type == html.TextNode {
					text += child.Data
				}
				if child.Type == html.ElementNode {
					for innerChild := range child.ChildNodes() {
						if innerChild.Type == html.TextNode {
							text += innerChild.Data
						}
					}
				}
			}

			link.Text = strings.Trim(text, " \n")

			links = append(links, link)
		}
	}

	return links, nil
}

func getHTMLByPath(HTMLPath string) (string, error) {
	htmlData, err := os.ReadFile(HTMLPath)
	if err != nil {
		return "", err
	}

	return string(htmlData), err
}

func getHTMLByURL(URL string) (string, error) {
	res, err := http.Get(URL)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch: status code %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), err
}
