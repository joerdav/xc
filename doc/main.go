package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func customFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(r.URL.Path)
		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				// Serve the custom 404 page
				http.ServeFile(w, r, "public/404.html")
				return
			}
			// For other errors, serve the default 500 page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Serve the file if found
		http.FileServer(fs).ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      customFileServer(http.Dir("public")),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
