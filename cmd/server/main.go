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
	// Serve secure (HTTPS) routes
	secureMux := http.NewServeMux()
	secureMux.HandleFunc("/", auth.AuthMiddleware(handlers.HomeHandler))
	secureMux.HandleFunc("/login", handlers.LoginHandler)
	secureMux.HandleFunc("/upload", auth.AuthMiddleware(handlers.UploadHandler))
	secureMux.HandleFunc("/stream", auth.AuthMiddleware(video.StreamHandler))
	secureMux.HandleFunc("/api/list-files", auth.AuthMiddleware(handlers.ListFilesHandler))
	secureMux.HandleFunc("/file/", auth.AuthMiddleware(handlers.ServeFileHandler))

	// Serve static files from ../uploads
	fs := http.FileServer(http.Dir("../uploads")) //Set disk name or folder that want to be served.
	secureMux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Port configuration
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8082" // Default HTTP port
	}
	httpsPort := os.Getenv("HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = "8080" // Default HTTPS port
	}

	// TLS certificate paths
	certFile := "../../cert.pem"
	keyFile := "../../key.pem"

	// Start HTTP server with redirection to HTTPS
	go func() {
		log.Printf("Starting HTTP server with redirection on http://localhost:%s", httpPort)
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		if err := http.ListenAndServe(":"+httpPort, httpMux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Check if the cert.pem and key.pem files exist for HTTPS
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatal("TLS certificate file not found:", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatal("TLS key file not found:", keyFile)
	}

	// Start HTTPS server
	log.Printf("Starting HTTPS server on https://localhost:%s", httpsPort)
	if err := http.ListenAndServeTLS(":"+httpsPort, certFile, keyFile, secureMux); err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
