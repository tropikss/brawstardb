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
	PlayerID   string `json:"player_id"`
	Name string `json:"name"`
}

type Battle struct {
	PlayerID string `json:"player_id"`
	BattleTime string `json:"battle_time"`
	Result    string `json:"result"`
	Mode string `json:"mode"`
	Type  string `json:"type"`
	Map string `json:"map"`
	StarPlayer bool   `json:"star_player"`
	Duration int    `json:"duration"`
	TrophyChange int `json:"trophy_change"`
	Teams []Team `json:"teams"`
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
	http.HandleFunc("/battles", battleHandler)

	fmt.Println("Serveur lance sur :8000")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "API en ligne"})
}

func battleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS battles (
		id SERIAL PRIMARY KEY,
		player_id TEXT,
		battle_time TEXT,
		result TEXT,
		mode TEXT,
		type TEXT,
		map TEXT,
		star_player BOOLEAN,
		duration INT,
		trophy_change INT
	)`)
	switch r.Method {
	case "GET":
		listBattles(w)
	case "POST":
		addBattle(w, r)
	default:
		http.Error(w, "Methode non supportee", http.StatusMethodNotAllowed)
	}
}

func addBattle(w http.ResponseWriter, r *http.Request) {
	// Décoder l'objet JSON reçu
	var b Battle
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Vérifier si ca existe pas déjà
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM battles WHERE player_id=$1 AND battle_time=$2
	)`, b.PlayerID, b.BattleTime).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "already exists"})
		return
	}

	result := b.battle.result
	starplayer := b.battle.star_player.tag == b.PlayerID

	_, err = db.Exec(`INSERT INTO battles 
		(player_id, battle_time, result, mode, type, map, star_player, duration, trophy_change) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		b.playerID, b.battleTime, result, b.battle.mode, b.battle.type, b.event.map, starplayer, b.battle.duration, b.battle.trophyChange)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func listBattles(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, player_id, battle_time, result, mode, type, map, star_player, duration, trophy_change FROM battles")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	var battles []Battle
	for rows.Next() {
		var b Battle
		rows.Scan(&b.ID, &b.PlayerID, &b.BattleTime, &b.Result, &b.Mode, &b.Type, &b.Map, &b.StarPlayer, &b.Duration, &b.TrophyChange)
		battles = append(battles, b)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battles)
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
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", u.Name, u.Age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}