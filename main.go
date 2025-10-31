package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type StarPlayer struct {
	Tag string `json:"tag"`
}

type BattleDetail struct {
	Mode        string     `json:"mode"`
	BattleType  string     `json:"type"`
	Result      string     `json:"result"`
	Duration    int        `json:"duration"`
	TrophyChange int       `json:"trophyChange"`
	StarPlayer  StarPlayer `json:"starPlayer"`
}

type Event struct {
	Map string `json:"map"`
}

type Battle struct {
	PlayerID   string       `json:"playerId"`
	BattleTime string       `json:"battleTime"`
	Battle     BattleDetail `json:"battle"`
	Event      Event        `json:"event"`
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
	var b Battle
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Vérifie si la battle existe déjà
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

	// Vérifie si le joueur est starPlayer
	starplayer := b.Battle.StarPlayer.Tag == b.PlayerID

	// INSERT dans la DB
	_, err = db.Exec(`INSERT INTO battles 
		(player_id, battle_time, result, mode, battle_type, map, star_player, duration, trophy_change) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		b.PlayerID, b.BattleTime, b.Battle.Result, b.Battle.Mode, b.Battle.BattleType,
		b.Event.Map, starplayer, b.Battle.Duration, b.Battle.TrophyChange)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func listBattles(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, player_id, battle_time, result, mode, battle_type, map, star_player, duration, trophy_change FROM battles")
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