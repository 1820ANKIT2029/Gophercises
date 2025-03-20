package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// Define option struct
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

// Define chapter struct
type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

// Main story map: map of arc name (string) to Chapter
type Story map[string]Chapter

func main() {
	// flags
	storyfilename := flag.String("f", "gopher.json", "Story filename json format eg: -f gopher.json")
	flag.Parse()

	// parsing
	story, err := parseStory(*storyfilename)
	if err != nil {
		panic(err)
	}

	handler := storyHandler(story)

	http.HandleFunc("/", handler)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", nil)
}

func storyHandler(story Story) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arc := strings.Trim(r.URL.Path, "/")
		if arc == "" {
			arc = "intro" // default starting arc
		}
		chapter, ok := story[arc]
		if !ok {
			http.NotFound(w, r)
			return
		}
		tmpl := template.Must(template.ParseFiles("template.html"))
		tmpl.Execute(w, chapter)
	}
}

func parseStory(filename string) (Story, error) {
	jsn, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var story Story
	err = json.Unmarshal(jsn, &story)
	if err != nil {
		return nil, err
	}

	return story, nil
}
