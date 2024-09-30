package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the upload page if it's a GET request
	if r.Method == http.MethodPost {
		// Get the folder name from the form value (submitted with the POST request)
		folder := r.URL.Query().Get("folder")

		// Set the base upload directory
		baseUploadDir := "../../cmd/uploads/"
		folderPath := baseUploadDir

		// If folder is provided, join it to the base upload directory
		if folder != "" {
			// Create the folder path
			folderPath = filepath.Join(baseUploadDir, folder)

			// Ensure the folder exists, if not create it
			err := os.MkdirAll(folderPath, os.ModePerm)
			if err != nil {
				http.Error(w, "Error creating folder: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Handle file upload
		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get the filename from the uploaded file
		filename := header.Filename

		fmt.Println("Uploading and saving file:", filename)

		// Create the destination file path
		destFilePath := filepath.Join(folderPath, filename)

		// Create the destination file
		outFile, err := os.Create(destFilePath)
		if err != nil {
			http.Error(w, "Unable to create file for saving", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy the contents of the uploaded file to the destination file
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Respond to the client
		fmt.Fprintf(w, "File uploaded successfully to folder: %s", folder)
		return
	}

	if r.Method == http.MethodDelete {
		// Get the file name from the query parameters or request body
		fileName := r.URL.Query().Get("filename") // Example if the file is passed as a query param

		fmt.Println("this is in the upload delete", fileName)

		if fileName == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}

		// Construct the file path
		filePath := fmt.Sprintf("../uploads/%s", fileName) // Modify the path as per your file structure

		fmt.Println("this is path", filePath)

		// Try to delete the file

		err := os.Remove(filePath)
		if err != nil {
			http.Error(w, "Error deleting file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "File %s deleted successfully", fileName)
	}
	fmt.Println(r.Method)
	if r.Method == "ADD" {
		basePath := "../../cmd/uploads"
		newFolderName := r.URL.Query().Get("foldername")
		folderPath := r.URL.Query().Get("path")

		fmt.Println(newFolderName, folderPath)
		if folderPath != "" {
			basePath = filepath.Join(basePath, folderPath) // Default path
		}

		if newFolderName == "" {
			http.Error(w, "Folder name is required", http.StatusBadRequest)
			return
		}

		// Construct the folder path
		newFolderPath := filepath.Join(basePath, newFolderName)

		// Create the folder
		err := os.Mkdir(newFolderPath, os.ModePerm)
		if err != nil {
			http.Error(w, "Error creating folder: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Folder %s created successfully", newFolderName)
		return
	}

	// If the request method is not supported, return an error
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
