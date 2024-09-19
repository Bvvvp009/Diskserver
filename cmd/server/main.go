package main

import (
	"diskserver/internal/auth"
	"diskserver/internal/handlers"
	"diskserver/internal/video"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", auth.AuthMiddleware(handlers.HomeHandler))
	http.HandleFunc("/login", handlers.LoginHandler)

	http.HandleFunc("/upload", auth.AuthMiddleware(handlers.UploadHandler))
	http.HandleFunc("/stream", auth.AuthMiddleware(video.StreamHandler))
	http.HandleFunc("/api/list-files", auth.AuthMiddleware(handlers.ListFilesHandler))
	http.HandleFunc("/file/", auth.AuthMiddleware(handlers.ServeFileHandler))

	// Serve static files
	fs := http.FileServer(http.Dir("../uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
