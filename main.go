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

var db *sql.DB

func init() {
	var err error
	connStr := "postgresql://postgres:ldyvojdy8yaj13rz@brawl-star-db-7hbene:5432/postgres?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur de connexion DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB injoignable:", err)
	}
	fmt.Println("✅ Connecte a la base Postgres")

	createTable()
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/users", usersHandler)

	fmt.Println("Serveur lance sur :8000")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "API en ligne"})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listUsers(w)
	case "POST":
		addUser(w, r)
	default:
		http.Error(w, "Methode non supportee", http.StatusMethodNotAllowed)
	}
}

func listUsers(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Age)
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, err := db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", u.Name, u.Age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func createTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		age INT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}
