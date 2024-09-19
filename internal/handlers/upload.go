package handlers

import (
	"fmt"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the upload page if it's a GET request
	if r.Method == "GET" {
		http.ServeFile(w, r, "../../static/upload.html") // Relative path
		return
	}

	// Handle file upload if it's a POST request
	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := header.Filename
		err = SaveEncryptedFile(filename, file) // Save the file (implement your encryption here)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully: %s", filename)
		return
	}

	// If the request method is not supported, return an error
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
