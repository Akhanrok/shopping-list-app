package main

import (
	"log"
	"net/http"

	"github.com/Akhanrok/shopping-list-app/handlers"
	"github.com/Akhanrok/shopping-list-apps/repositories"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create a database connection
	var err error
	db, err := repositories.NewDatabase("root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.IndexHandler(w, r)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, db)
	})

	http.HandleFunc("/login-success", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginSuccessHandler(w, r, db)
	})

	http.HandleFunc("/register-success", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterSuccessHandler(w, r)
	})

	http.HandleFunc("/create-list", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateListHandler(w, r, db)
	})

	http.HandleFunc("/list-success", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListSuccessHandler(w, r)
	})

	http.HandleFunc("/view-lists", func(w http.ResponseWriter, r *http.Request) {
		handlers.ViewListsHandler(w, r, db)
	})

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
