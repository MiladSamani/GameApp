package main

import (
	"encoding/json"
	"fmt"
	"gameAppProject/repository/mysql"
	"gameAppProject/service/userservice"
	"io"
	"log"
	"net/http"
)

func main() {
	// Create a new ServeMux to handle different routes
	mux := http.NewServeMux()

	// Define endpoints and their corresponding handler functions
	mux.HandleFunc("/health-check", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)

	// Create an HTTP server listening on port 8080 with the defined ServeMux
	server := http.Server{Addr: ":8080", Handler: mux}

	// Start the server and log any errors
	log.Println("server is listening on port 8080...")
	log.Fatal(server.ListenAndServe())
}

// userRegisterHandler handles user registration requests.
func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	// Check if the request method is POST, respond with an error if not
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error": "invalid method"}`)
		return
	}

	// Read the request body
	data, err := io.ReadAll(req.Body)
	if err != nil {
		// Respond with an error if there's an issue reading the request body
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	// Decode JSON request body into a RegisterRequest struct
	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		// Respond with an error if there's an issue decoding JSON
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	// Initialize MySQL repository and user service
	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)

	// Attempt to register the user
	_, err = userSvc.Register(uReq)
	if err != nil {
		// Respond with an error if user registration fails
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	// Respond with a success message if user registration is successful
	writer.Write([]byte(`{"message": "user created"}`))
}

// healthCheckHandler handles health check requests.
func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	// Respond with a JSON message indicating that everything is good
	fmt.Fprintf(writer, `{"message": "everything is good!"}`)
}
