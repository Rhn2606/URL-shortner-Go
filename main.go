package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OgURL        string    `json:"OgURL"`
	ShortURL     string    `json:"shortURL"`
	CreationDate time.Time `json:"date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OgURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OgURL))
	fmt.Println("hasher :", hasher)
	data := hasher.Sum(nil)
	fmt.Println("Hasher data", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Shortened URL", hash[:8])
	return hash[:8]
}

func createNewURL(OgURL string) string {
	shortURL := generateShortURL(OgURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OgURL:        OgURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getOriginalURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"URL"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	shortURL := createNewURL(data.URL)
	// fmt.Fprintf(w, shortURL)
	response := struct {
		Url string `json:"shortURL"`
	}{Url: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getOriginalURL(id)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OgURL, http.StatusFound)
}

func main() {
	fmt.Println("Starting URL shortner ...")
	//OgURL := "https://google.com"
	//generateShortURL(OgURL)

	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURL)

	fmt.Println("server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error on starting server", err)
	}

}
