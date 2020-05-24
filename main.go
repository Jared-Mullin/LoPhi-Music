package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var (
	accessToken  string
	refreshToken string
)

const (
	state = "state"
)

func main() {
	router := chi.NewRouter()
	spotifyConf := setupSpotifyClient()

	//Middleware Stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/spotify/auth", func(w http.ResponseWriter, r *http.Request) {
		url := spotifyConf.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	router.Get("/spotify/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != state {
			log.Println("Invalid OAuth2 State")
		} else {
			token, err := spotifyConf.Exchange(oauth2.NoContext, r.FormValue("code"))
			if err != nil {
				log.Println(err)
			} else {
				accessToken = token.AccessToken
				refreshToken = token.RefreshToken
			}
		}
	})

	router.Get("/spotify/genres", func(w http.ResponseWriter, r *http.Request) {

	})

	http.ListenAndServe(":4200", router)
}

func setupSpotifyClient() *oauth2.Config {
	cID := os.Getenv("SPOTIFY_CLIENT_ID")
	cSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	conf := &oauth2.Config{
		ClientID:     cID,
		ClientSecret: cSecret,
		RedirectURL:  "http://localhost:4200/spotify/callback",
		Scopes:       []string{"user-top-read"},
		Endpoint:     spotify.Endpoint,
	}
	return conf
}
