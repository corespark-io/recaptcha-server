package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/mitchs-dev/library-go/loggingFormatter"
	log "github.com/sirupsen/logrus"
)

var (
	// Environment variables
	recaptchaSecretKey   string
	timezone             string
	port                 string
	loadedFromDotenv     bool
	logLevel             string
	recaptchaAPIEndpoint string
	frontendURL          string

	// Default values
	defaultTimezone    = "UTC"
	defaultPort        = "8080"
	defaultLogLevel    = "info"
	defaultFrontendURL = "*" // Allow all origins by default
)

// Default reCAPTCHA API endpoint
const (
	defaultRecaptchaAPIEndpoint = "https://www.google.com/recaptcha/api/siteverify"
)

func init() {
	// Check if there is a .env file and load it if it exists
	if _, err := os.Stat(".env"); err == nil {
		loadedFromDotenv = true
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Failed to load .env file: ", err)
		}
	}

	// Set the timezone if RECAPTCHA_TIMEZONE environment variable is set
	timezone = os.Getenv("RECAPTCHA_TIMEZONE")
	if timezone == "" {
		timezone = defaultTimezone // Default to UTC if not set
	}
	// Set log format
	log.SetFormatter(&loggingFormatter.JSONFormatter{
		Timezone: timezone,
		Prefix:   "recaptcha-server-",
	})

	// Set log level
	logLevel = os.Getenv("RECAPTCHA_LOG_LEVEL")
	if logLevel == "" {
		logLevel = defaultLogLevel // Default to 'info' if not set
	}

	// Set output to standard output
	log.SetOutput(os.Stdout)

	// Ensure the 'RECAPTCHA_SECRET_KEY' environment variable is set
	recaptchaSecretKey = os.Getenv("RECAPTCHA_SECRET_KEY")
	if recaptchaSecretKey == "" {
		log.Fatal("RECAPTCHA_SECRET_KEY environment variable is not set")
	}

	// Set the port from the environment variable or default to 8080
	port = os.Getenv("RECAPTCHA_PORT")
	if port == "" {
		port = defaultPort // Default to 8080 if not set
	}

	// Set the reCAPTCHA API endpoint
	recaptchaAPIEndpoint = os.Getenv("RECAPTCHA_API_ENDPOINT")
	if recaptchaAPIEndpoint == "" {
		recaptchaAPIEndpoint = defaultRecaptchaAPIEndpoint
	}

	// Set the frontend URL for CORS
	frontendURL = os.Getenv("RECAPTCHA_FRONTEND")
	if frontendURL == "" {
		frontendURL = defaultFrontendURL
	}

	log.Info("CORS configured to allow requests from: ", frontendURL)
	if frontendURL == "*" {
		log.Warn("CORS is set to allow all origins. This is not recommended for production environments. You should set the RECAPTCHA_FRONTEND environment variable to a specific URL.")
	}

	// Show .env warning if loaded
	if loadedFromDotenv {
		log.Warn("You should not use .env files in production. This should only be used for development purposes.")
	}
}

// enableCORS is a middleware that sets CORS headers
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if this is a CORS request
		if origin != "" {
			// Check if the origin is allowed
			if frontendURL != "*" && origin != frontendURL {
				log.Warn("Rejected CORS request from unauthorized origin: ", origin)
				// Still set CORS headers to ensure browser gets proper response
			}

			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", frontendURL)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
func main() {
	log.Info("Starting recaptcha server on port: ", port)

	// Handle verify endpoint with CORS support
	http.HandleFunc("/verify", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Warn("Received non-POST request to /verify endpoint")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var requestData struct {
			Token string `json:"token"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			log.Warnf("Failed to parse JSON request: %v", err)
			http.Error(w, "Failed to parse JSON request", http.StatusBadRequest)
			return
		}

		if requestData.Token == "" {
			log.Warn("Missing recaptcha token in request")
			http.Error(w, "Missing recaptcha token", http.StatusBadRequest)
			return
		}

		// Verify token with Google's recaptcha API using the configurable endpoint
		resp, err := http.PostForm(recaptchaAPIEndpoint, url.Values{
			"secret":   {recaptchaSecretKey},
			"response": {requestData.Token},
			"remoteip": {r.RemoteAddr},
		})
		if err != nil {
			log.Error("Failed to verify token: ", err)
			http.Error(w, "Failed to verify token", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Parse and return the verification result
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Error("Failed to decode response: ", err)
			http.Error(w, "Failed to decode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	log.Info("Server listening on port ", port)
	log.Info("Requests can be made to the /verify endpoint with a POST request containing a JSON body with a 'token' field")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
