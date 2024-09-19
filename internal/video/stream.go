package video

import (
    "net/http"
)

func StreamHandler(w http.ResponseWriter, r *http.Request) {
    videoPath := "uploads/video.mp4"
    http.ServeFile(w, r, videoPath)
}
