package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// Artist Spotify Response Structure
type Artist struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Genres []string `json:"genres"`
	Href   string   `json:"href"`
	ID     string   `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Type       string `json:"type"`
	URI        string `json:"uri"`
}

var (
	accessToken  string
	refreshToken string
)

const (
	state = "state"
)

func main() {
	router := chi.NewRouter()
	spotifyConf := setupSpotifyConf()

	//Middleware Stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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

	router.Get("/spotify/artists", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(accessToken)
		client := http.Client{}
		req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/artists", nil)
		if err != nil {
			log.Println("Error in Creating Request")
			log.Println(err)
		} else {
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res, err := client.Do(req)
			if err != nil {
				log.Println("Error in Performing Request")
				log.Println(err)
			} else {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				}
				w.Write(body)
			}
		}
	})

	router.Get("/spotify/tracks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(accessToken)
		client := http.Client{}
		req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks", nil)
		if err != nil {
			log.Println("Error in Creating Request")
			log.Println(err)
		} else {
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res, err := client.Do(req)
			if err != nil {
				log.Println("Error in Performing Request")
				log.Println(err)
			} else {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				}
				w.Write(body)
			}
		}
	})

	http.ListenAndServe(":4200", router)
}

func setupSpotifyConf() *oauth2.Config {
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
