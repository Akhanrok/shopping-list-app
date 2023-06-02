package repositories

import (
	"database/sql"
)

type FlashMessage struct {
	Message string
}

type Product struct {
	Product  string
	Quantity int
	Store    string
}

type ListData struct {
	ListName string
	Products []Product
}

var db *sql.DB

type UserRepository interface {
	ValidateCredentials(email, password string) (string, error)
	IsEmailExists(email string) (bool, error)
}

type ListRepository interface {
	IsListExists(listName string) (bool, error)
	GetListsData(userID int) ([]ListData, error)
}

type ProductRepository interface {
	GetProductsData(listID int) ([]Product, error)
}

type userRepository struct {
	db *sql.DB
}

type listRepository struct {
	db *sql.DB
}

type productRepository struct {
	db *sql.DB
}

// Create a new database connection
func NewDatabase(dataSourceName string) (*sql.DB, error) {
	// Initialize the database connection
	database, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	db = database
	return db, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func NewListRepository(db *sql.DB) ListRepository {
	return &listRepository{db}
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *userRepository) ValidateCredentials(email, password string) (string, error) {
	query := "SELECT name FROM users WHERE email = ? AND password = ?"
	var username string
	err := r.db.QueryRow(query, email, password).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (r *userRepository) IsEmailExists(email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *listRepository) IsListExists(listName string) (bool, error) {
	query := "SELECT COUNT(*) FROM lists WHERE name = ?"
	var count int
	err := r.db.QueryRow(query, listName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *listRepository) GetListsData(userID int) ([]ListData, error) {
	query := "SELECT id, name FROM lists WHERE user_id = ?"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []ListData

	for rows.Next() {
		var listID int
		var listName string

		err := rows.Scan(&listID, &listName)
		if err != nil {
			return nil, err
		}

		products, err := NewProductRepository(r.db).GetProductsData(listID)
		if err != nil {
			return nil, err
		}

		listData := ListData{
			ListName: listName,
			Products: products,
		}

		lists = append(lists, listData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *productRepository) GetProductsData(listID int) ([]Product, error) {
	query := "SELECT name, quantity, store FROM products WHERE list_id = ?"
	rows, err := r.db.Query(query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var name string
		var quantity int
		var store string

		err := rows.Scan(&name, &quantity, &store)
		if err != nil {
			return nil, err
		}

		product := Product{
			Product:  name,
			Quantity: quantity,
			Store:    store,
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
