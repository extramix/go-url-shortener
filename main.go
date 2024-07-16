package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	urls map[string]string // key: short url, value: long url
}

const baseURL = "http://localhost:8080/shorts"
const shortenRoute = "/shorten"
const shortsRoute = "/shorts/"

func main() {
	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc(shortenRoute, shortener.HandleShorten)
	http.HandleFunc(shortsRoute, shortener.HandleRedirect)

	fmt.Println("Server is running on port 8080")
	fmt.Println("To shorten a URL, send a POST request to http://localhost:8080/shorten")
	fmt.Println("Example: curl -X POST -d 'url=https://google.com' http://localhost:8080/shorten")
	fmt.Println("To redirect to the original URL, go to http://localhost:8080/shorts/<shortened-url>")
	http.ListenAndServe(":8080", nil)

}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 5

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	shortenedURL := baseURL + shortKey

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: %s</p>
        <p>Shortened URL: <a href="%s">%s</a></p>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
        </form>
    `, originalURL, shortenedURL, shortenedURL)
	//print responseHTML
	fmt.Fprint(w, responseHTML)

}

func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/shorts/"):]
	originalURL, ok := us.urls[shortKey]
	if !ok {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}
