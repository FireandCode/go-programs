package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jxskiss/base62"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type URL struct {
	ID    int64   `gorm:"primaryKey" json:"id"`
	LongURL string `json:"long_url"`
}

// Status represents the status of a task.
type Status uint

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

var DB *gorm.DB

func InitDB() {
	
	dsn := "root:user@1234@tcp(127.0.0.1:3306)/golangDB?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	// Auto-migrate will create/update the tables if needed.
	DB.AutoMigrate(&URL{})
}

// writeError sends a JSON error response with the given status and message.
func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Proceed to the next handler.
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.RequestURI, duration)
	})
}

type Response struct {
	HashValue string `json:"hash_value"`
}

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	var URL URL
	json.NewDecoder(r.Body).Decode(&URL)


	DB.Create(&URL)

	//base62 encode
	var response Response
	
	response.HashValue = string(base62.Encode([]byte(strconv.Itoa(int(URL.ID)))))

	json.NewEncoder(w).Encode(response)
}

func GetLongURL(w http.ResponseWriter, r *http.Request) {
	var URL URL

	vars := mux.Vars(r)
	hashValue := vars["hashValue"]


	id,_ :=base62.Decode([]byte(hashValue))

	strId := string(id)

	URL.ID, _ = strconv.ParseInt(strId, 10, 64)
	
	log.Println(URL)
	DB.Find(&URL)

	log.Println(URL)
	json.NewEncoder(w).Encode(URL.LongURL)
}

func main() {
	// Initialize the database connection.
	InitDB()

	// Create a new Gorilla Mux router.
	r := mux.NewRouter()

	// Register middleware.
	r.Use(loggingMiddleware)

	// Define user-related routes.
	r.HandleFunc("/convert_to_short_curl", GetShortURL).Methods("POST")
	r.HandleFunc("/{hashValue}", GetLongURL).Methods("GET")
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
