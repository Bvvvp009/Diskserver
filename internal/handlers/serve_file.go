package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ServeFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/file/")
	filePath := filepath.Join("../../cmd/uploads", filename)

	// fmt.Println(filename, filePath)

	data, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	defer data.Close()

	w.Header().Set("Content-Disposition", "inline; filename="+filename)
	http.ServeContent(w, r, filename, GetFileModTime(filePath), data)
}
