package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	API_TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6ImJlOTgzM2JlLTdmMTQtNDYxMy05NTNiLTg1MzgxYzgzZTJiMCIsImlhdCI6MTc2MTgyMzY4MSwic3ViIjoiZGV2ZWxvcGVyLzliNjI3Zjg2LTZiMjctNTFjZS1jNWQ4LWMwYWNkMGM3MWI0NyIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiODIuNjYuMTEwLjk3Il0sInR5cGUiOiJjbGllbnQifV19.huIuBsE_m3P1YNZyYJkQ0ZoeyIFKbW7G6epjz4EyvLhYB2rjTpE-49zWK163F0n7wn6btpRYpQdC4vAqA_d7JQ"
	BASE_URL  = "https://api.brawlstars.com/v1"
)

var playersTag = []string{"2LGRYGVP"}

func getPlayer(tag string) (map[string]interface{}, error) {
	tagEscaped := url.PathEscape(tag)
	req, err := http.NewRequest("GET", BASE_URL+"/players/"+tagEscaped+"/battlelog", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Erreur %d: %s", resp.StatusCode, string(body))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func postBattle(url string, battleData map[string]interface{}) error {
	body, err := json.Marshal(battleData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Erreur %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func main() {
	for _, playerTag := range playersTag {
		data, err := getPlayer("#" + playerTag)
		if err != nil {
			fmt.Println("Failed to get player:", err)
			continue
		}
		items, ok := data["items"].([]interface{})
		if !ok {
			fmt.Println("No items found for player:", playerTag)
			continue
		}
		for _, item := range items {
			battle, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			battle["playerId"] = playerTag
			if err := postBattle("http://localhost:8000/battles", battle); err != nil {
				fmt.Println("Failed to post battle:", err)
			} else {
				fmt.Println("Battle posted successfully")
			}
		}
	}
}
