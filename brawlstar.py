import requests

API_TOKEN = "ton_token_api"
BASE_URL = "https://api.brawlstars.com/v1"
API_TOKEN = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6ImJlOTgzM2JlLTdmMTQtNDYxMy05NTNiLTg1MzgxYzgzZTJiMCIsImlhdCI6MTc2MTgyMzY4MSwic3ViIjoiZGV2ZWxvcGVyLzliNjI3Zjg2LTZiMjctNTFjZS1jNWQ4LWMwYWNkMGM3MWI0NyIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiODIuNjYuMTEwLjk3Il0sInR5cGUiOiJjbGllbnQifV19.huIuBsE_m3P1YNZyYJkQ0ZoeyIFKbW7G6epjz4EyvLhYB2rjTpE-49zWK163F0n7wn6btpRYpQdC4vAqA_d7JQ'

headers = {
    "Authorization": f"Bearer {API_TOKEN}"
}

players_name = ['Mathis']
players_tag = ['2LGRYGVP']


def post_battle(url, battle_data):
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, json=battle_data, headers=headers)
    
    if response.status_code == 200:
        return response.json()
    else:
        raise Exception(f"Erreur {response.status_code}: {response.text}")

def get_player(tag):
    tag = tag.replace("#", "%23")
    url = f"{BASE_URL}/players/{tag}/battlelog"
    r = requests.get(url, headers=headers)
    if r.ok:
        return r.json()
    else:
        print(f"Erreur {r.status_code}: {r.text}")
        return None

if __name__ == "__main__":
    for player_tag in players_tag:
        data = get_player(f"#{player_tag}")
        if data:
            for battle in data['items']:
                battle['playerId'] = player_tag
                try:
                    response = post_battle("http://localhost:8000/battles", battle)
                    print("Battle posted successfully:", response)
                except Exception as e:
                    print("Failed to post battle:", e)
