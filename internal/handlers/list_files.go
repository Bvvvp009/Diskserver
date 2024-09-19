package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileGroup represents a group of files with a common date
type FileGroup struct {
	Date  string     `json:"date"`
	Files []FileInfo `json:"files"`
}

// FileInfo holds information about a file
type FileInfo struct {
	Name    string `json:"name"`
	IsImage bool   `json:"is_image"`
	IsVideo bool   `json:"is_video"`
	IsDoc   bool   `json:"is_doc"`
}

// ListFilesHandler lists the files and returns them as JSON
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("../../cmd/uploads")
	if err != nil {
		http.Error(w, "Unable to list files", http.StatusInternalServerError)
		return
	}
	fmt.Println(files)

	// Group files by their modification date
	fileGroups := make(map[string][]FileInfo)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		fileName := file.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		modTime := info.ModTime().Format("2006-01-02")

		fileInfo := FileInfo{
			Name:    fileName,
			IsImage: isImage(ext),
			IsVideo: isVideo(ext),
			IsDoc:   isDoc(ext),
		}

		fileGroups[modTime] = append(fileGroups[modTime], fileInfo)
	}

	// Convert the map to a slice and sort it by date
	var sortedGroups []FileGroup
	for date, files := range fileGroups {
		sortedGroups = append(sortedGroups, FileGroup{Date: date, Files: files})
	}
	sort.Slice(sortedGroups, func(i, j int) bool {
		return sortedGroups[i].Date > sortedGroups[j].Date
	})

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the sorted groups into JSON
	if err := json.NewEncoder(w).Encode(sortedGroups); err != nil {
		http.Error(w, "Unable to encode file data to JSON", http.StatusInternalServerError)
		return
	}
}

// Helper functions to check file types
func isImage(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func isVideo(ext string) bool {
	return ext == ".mp4" || ext == ".mov" || ext == ".avi"
}

func isDoc(ext string) bool {
	return ext == ".pdf" || ext == ".docx" || ext == ".txt"
}
