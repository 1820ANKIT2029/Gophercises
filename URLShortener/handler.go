package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
)

type pathUrlyaml struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type pathUrljson struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlpathname string, fallback http.Handler) (http.HandlerFunc, error) {
	var records []pathUrlyaml

	yml, err := os.Open(yamlpathname)
	if err != nil {
		return nil, err
	}

	yamldecoder := yaml.NewDecoder(yml)
	err = yamldecoder.Decode(&records)
	if err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)
	for _, record := range records {
		pathsToUrls[record.Path] = record.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsonpathname string, fallback http.Handler) (http.HandlerFunc, error) {
	var records []pathUrljson

	jsn, err := os.Open(jsonpathname)
	if err != nil {
		return nil, err
	}

	jsondecoder := json.NewDecoder(jsn)
	err = jsondecoder.Decode(&records)
	if err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)
	for _, record := range records {
		pathsToUrls[record.Path] = record.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func DBHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		dest, err := getURL(path)
		if err == nil {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func saveShorten(path, url string) (string, error) {
	path = "/" + path
	err := saveURL(path, url)
	if err != nil {
		return "", err
	}

	return path, nil
}

func saveURL(path, url string) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(shortlinkBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		value := b.Get([]byte(path))
		if value != nil {
			return fmt.Errorf("value already exist")
		}
		err = b.Put([]byte(path), []byte(url))
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func getURL(path string) (string, error) {
	var value []byte
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(shortlinkBucket))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		value = b.Get([]byte(path))
		return nil
	})

	if err != nil {
		return "", err
	}

	if value == nil {
		return "", fmt.Errorf("no value exist")
	}

	return string(value), nil
}
