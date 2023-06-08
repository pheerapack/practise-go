package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "host=localhost port=5432 user=your_username password=your_password dbname=your_database sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define API endpoint to retrieve users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// Query data from the database
		rows, err := db.Query("SELECT id, name, age FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Fetch the result and store in a slice of User struct
		var users []User
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Age)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the users slice to JSON
		jsonData, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		 }

		// Set the response content type and write the JSON data
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	// Start the server on port 8080
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
