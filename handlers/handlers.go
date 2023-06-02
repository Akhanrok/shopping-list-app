package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Akhanrok/shopping-list-app/repositories"
	"github.com/Akhanrok/shopping-list-app/services"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		services.RenderTemplate(w, "index.html", nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// Create instances of the repositories
		userRepo := repositories.NewUserRepository(db)

		// Check the credentials in the database
		username, err := userRepo.ValidateCredentials(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if username == "" {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Wrong credentials",
			}
			services.RenderTemplate(w, "login.html", data)
			return
		}

		data := struct {
			Username string
		}{
			Username: username,
		}

		services.RenderTemplate(w, "login-success.html", data)
		return
	}

	services.RenderTemplate(w, "login.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.PostForm.Get("name")
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		if name == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if !services.IsValidEmail(email) {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Password should be at least 8 characters long",
			}
			services.RenderTemplate(w, "register.html", data)
			return
		}

		// Create instances of the repositories
		userRepo := repositories.NewUserRepository(db)

		// Check if the email already exists in the database
		emailExists, err := userRepo.IsEmailExists(email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if emailExists {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Email already exists",
			}
			services.RenderTemplate(w, "register.html", data)
			return
		}

		// Insert the new user into the database
		insertQuery := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
		_, err = db.Exec(insertQuery, name, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the register success page
		http.Redirect(w, r, "/register-success", http.StatusFound)
		return
	}

	services.RenderTemplate(w, "register.html", nil)
}

func LoginSuccessHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		services.RenderTemplate(w, "login-success.html", nil)
	}
}

func RegisterSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		services.RenderTemplate(w, "register-success.html", nil)
	}
}

func CreateListHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listName := r.PostForm.Get("listName")

		// Create instances of the repositories
		listRepo := repositories.NewListRepository(db)

		// Check if the list name already exists in the database
		listExists, err := listRepo.IsListExists(listName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if listExists {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "The list with such name already exists",
			}
			services.RenderTemplate(w, "create-list.html", data)
			return
		}

		// Insert the new list into the database
		insertListQuery := "INSERT INTO lists (user_id, name) VALUES (?, ?)"
		res, err := db.Exec(insertListQuery, 1, listName) // Replace 1 with the appropriate user ID
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listID, _ := res.LastInsertId()

		// Get the product, quantity, and store values from the form
		products := r.PostForm["product[]"]
		quantities := r.PostForm["quantity[]"]
		stores := r.PostForm["store[]"]

		// Insert each product into the database
		insertProductQuery := "INSERT INTO products (list_id, name, quantity, store) VALUES (?, ?, ?, ?)"
		for i := range products {
			_, err = db.Exec(insertProductQuery, listID, products[i], quantities[i], stores[i])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Redirect to the list success page
		http.Redirect(w, r, fmt.Sprintf("/list-success?name=%s", listName), http.StatusFound)
		return
	}

	services.RenderTemplate(w, "create-list.html", nil)
}

func ListSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		listName := r.URL.Query().Get("name")
		data := struct {
			ListName string
		}{
			ListName: listName,
		}
		services.RenderTemplate(w, "list-success.html", data)
	}
}

func ViewListsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		// Get the user ID of the currently authenticated user
		userID := 1 // TODO: user ID retrieval

		// Create instances of the repositories
		listRepo := repositories.NewListRepository(db)

		// Retrieve the list names and items for the user from the database
		lists, err := listRepo.GetListsData(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Lists []repositories.ListData
		}{
			Lists: lists,
		}

		services.RenderTemplate(w, "view-lists.html", data)
	}
}
