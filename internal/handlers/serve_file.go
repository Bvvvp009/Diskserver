package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
)

func ServeFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/file/")
	filePath := filepath.Join("../../cmd/uploads", filename)

	data, err := ReadEncryptedFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "inline; filename="+filename)
	http.ServeContent(w, r, filename, GetFileModTime(filePath), data)
}
