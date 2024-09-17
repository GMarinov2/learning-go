package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
)

var shortenedUrls map[string]string = make(map[string]string)

func getMD5Hash(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Unrecognized body")
		return
	}

	var body ShortenRequest

	err = json.Unmarshal(requestBody, &body)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	hashed := getMD5Hash(body.Url)
	shortenedUrls[hashed] = body.Url

	response := ShortenResponse{
		Hash: hashed,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	io.Writer.Write(w, jsonData)
}

func redirectUrl(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	if hash == "" {
		w.WriteHeader(404)
		return
	}

	url := shortenedUrls[hash]
	http.Redirect(w, r, url, http.StatusSeeOther)

}

func main() {
	http.HandleFunc("/shorten", shortenUrl)
	http.HandleFunc("/{hash}", redirectUrl)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
