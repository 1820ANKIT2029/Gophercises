package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.etcd.io/bbolt"
)

var (
	DefaultURL      = ":8080"
	db              *bbolt.DB
	shortlinkBucket = "shortlink"
)

type ShortenRequest struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func main() {
	// flags
	yamlfilename := flag.String("y", "", "YAML filepath eg: -y URL.yaml")
	jsonfilename := flag.String("j", "", "JSON filepath eg: -j URL.json")
	DBfilename := flag.String("d", "", "Database filepath eg: -d my.db")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/ankit":          "https://github.com/1820ANKIT2029",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	var (
		handler http.HandlerFunc
		err     error
	)

	if *yamlfilename != "" {
		handler, err = YAMLHandler(*yamlfilename, mapHandler)
		if err != nil {
			panic(err)
		}
	} else if *jsonfilename != "" {
		fmt.Println(*jsonfilename)
		handler, err = JSONHandler(*jsonfilename, mapHandler)
		if err != nil {
			panic(err)
		}
	} else if *DBfilename != "" {
		db, err = bbolt.Open(*DBfilename, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(shortlinkBucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}

		handler = DBHandler(mux)
	} else {
		handler = mapHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(DefaultURL, handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/v1/shorten", shorten)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlFile := "index.html"

	content, err := os.ReadFile(htmlFile)
	if err != nil {
		log.Printf("Error reading HTML file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(content))
}

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}

	var reqData ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	path := reqData.Path
	url := reqData.URL

	shortenedURL, err := saveShorten(path, url)
	if err != nil {
		writeJSONError(w, http.StatusConflict, fmt.Sprintf("%v", err))
		return
	}

	// Set status and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{"shortUrl": shortenedURL}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received request - Path: %s, URL: %s, Shortened to: %s\n", path, url, shortenedURL)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
