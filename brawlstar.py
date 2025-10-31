import requests

API_TOKEN = "ton_token_api"
BASE_URL = "https://api.brawlstars.com/v1"
API_TOKEN = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6ImJlOTgzM2JlLTdmMTQtNDYxMy05NTNiLTg1MzgxYzgzZTJiMCIsImlhdCI6MTc2MTgyMzY4MSwic3ViIjoiZGV2ZWxvcGVyLzliNjI3Zjg2LTZiMjctNTFjZS1jNWQ4LWMwYWNkMGM3MWI0NyIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiODIuNjYuMTEwLjk3Il0sInR5cGUiOiJjbGllbnQifV19.huIuBsE_m3P1YNZyYJkQ0ZoeyIFKbW7G6epjz4EyvLhYB2rjTpE-49zWK163F0n7wn6btpRYpQdC4vAqA_d7JQ'

headers = {
    "Authorization": f"Bearer {API_TOKEN}"
}

player_tag = '2LGRYGVP'

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
    data = get_player(f"#{player_tag}")
    if data:
        print(data)
