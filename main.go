package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type Album struct {
	AlbumType        string   `json:"album_type"`
	Artists          []Artist `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	ExternalUrls     struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name                 string `json:"name"`
	ReleaseDate          string `json:"release_date"`
	ReleaseDatePrecision string `json:"release_date_precision"`
	TotalTracks          int    `json:"total_tracks"`
	Type                 string `json:"type"`
	URI                  string `json:"uri"`
}

type Track struct {
	Album            Album    `json:"album"`
	Artists          []Artist `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationMs       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIds      struct {
		Isrc string `json:"isrc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	IsLocal     bool   `json:"is_local"`
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	PreviewURL  string `json:"preview_url"`
	TrackNumber int    `json:"track_number"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}
type Items struct {
	Href  string   `json:"href"`
	Items []Artist `json:"items"`
}

type User struct {
	DisplayName  string `json:"display_name"`
	ExternalUrls struct {
		Spotify string `json:"spotify" bson:"spotify"`
	} `json:"external_urls" bson:"externalurls"`
	SpotifyID string `json:"id" bson:"spotifyid"`
	Images    []struct {
		Height int    `json:"height" bson:"height"`
		URL    string `json:"url" bson:"url"`
		Width  int    `json:"width" bson:"width"`
	} `json:"images" bson:"images"`
	AccessToken  string    `json:"access_token" bson:"accesstoken"`
	RefreshToken string    `json:"refresh_token" bson:"refreshtoken"`
	TokenExpiry  time.Time `bson:"tokenexpiry"`
	TokenType    string    `bson:"tokentype"`
}

var (
	mongoClient, mongoContext = createMongoClient()
	spotifyConf               = setupSpotifyConf()
)

func main() {
	r := chi.NewRouter()
	tokenAuth := setupJWTAuth()
	defer mongoClient.Disconnect(mongoContext)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/spotify", func(router chi.Router) {
		router.Get("/auth", func(w http.ResponseWriter, req *http.Request) {
			state := generateStateCookie(w)
			url := spotifyConf.AuthCodeURL(state)
			http.Redirect(w, req, url, http.StatusTemporaryRedirect)
		})

		router.Get("/callback", func(w http.ResponseWriter, req *http.Request) {
			state, err := req.Cookie("oauthstate")
			if err != nil {
				log.Println(err)
			} else {
				if req.FormValue("state") != state.Value {
					log.Println("Invalid OAuth2 State")
				} else {
					token, err := spotifyConf.Exchange(oauth2.NoContext, req.FormValue("code"))
					if err != nil {
						log.Println(err)
					} else {
						body, err := spotifyRequest(token, "https://api.spotify.com/v1/me/")
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							log.Println(err)
						} else {
							var user User
							json.Unmarshal(body, &user)
							user.AccessToken = token.AccessToken
							user.RefreshToken = token.RefreshToken
							user.TokenExpiry = token.Expiry
							user.TokenType = token.TokenType
							userCollection := mongoClient.Database("test").Collection("users")
							ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
							bUser, _ := bson.Marshal(user)
							exists, err := userCollection.CountDocuments(ctx, bson.M{"spotifyid": user.SpotifyID})
							if err != nil {
								log.Println("Error in Querying users Collection")
								log.Println(err)
							} else if exists == 1 {
								log.Println("User Already Exists")
							} else {
								_, err := userCollection.InsertOne(ctx, bUser)
								if err != nil {
									log.Println("Error in Performing Request")
									log.Println(err)
								} else {
									expiry := time.Now().Add(3 * 24 * time.Hour)
									_, tokenString, err := tokenAuth.Encode(jwt.MapClaims{"id": user.SpotifyID})
									if err != nil {
										log.Println("Error Creating Token")
									}
									http.SetCookie(w, &http.Cookie{Name: "token", Value: tokenString, Path: "/", Domain: "lophi.dev", Expires: expiry})
								}
							}
						}
					}
				}
			}
		})

		router.Group(func(spotifyRouter chi.Router) {
			spotifyRouter.Use(jwtauth.Verifier(tokenAuth))
			spotifyRouter.Use(jwtauth.Authenticator)

			spotifyRouter.Get("/artists", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				id := claims["id"].(string)
				token, err := getToken(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				}
				body, err := spotifyRequest(token, "https://api.spotify.com/v1/me/top/artists")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				}
				w.Write(body)
			})

			spotifyRouter.Get("/tracks", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				id := claims["id"].(string)
				token, err := getToken(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				}
				body, err := spotifyRequest(token, "https://api.spotify.com/v1/me/top/tracks")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				} else {
					w.Write(body)
				}
			})

			spotifyRouter.Get("/genres", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				id := claims["id"].(string)
				token, err := getToken(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				}
				body, err := spotifyRequest(token, "https://api.spotify.com/v1/me/top/artists")
				genres := make(map[string]int)
				var itemWrapper Items
				json.Unmarshal(body, &itemWrapper)
				for _, artist := range itemWrapper.Items {
					for _, genre := range artist.Genres {
						genres[genre] = genres[genre] + 1
					}
				}
				response, err := json.Marshal(genres)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
				} else {
					w.Write(response)
				}
			})
		})
	})
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "client/dist"))
	fileServer(r, "/", filesDir)

	http.ListenAndServe(":8080", r)
}

func setupSpotifyConf() *oauth2.Config {
	cID := os.Getenv("SPOTIFY_CLIENT_ID")
	cSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	conf := &oauth2.Config{
		ClientID:     cID,
		ClientSecret: cSecret,
		RedirectURL:  "https://lophi.dev/spotify/callback",
		Scopes:       []string{"user-read-private", "user-read-email", "user-top-read"},
		Endpoint:     spotify.Endpoint,
	}
	return conf
}

func setupJWTAuth() *jwtauth.JWTAuth {
	auth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	return auth
}

func createMongoClient() (*mongo.Client, context.Context) {
	uri := "mongodb+srv://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_URL") + "/test?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx
}

func fileHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}

func generateStateCookie(w http.ResponseWriter) string {
	expiry := time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiry}
	http.SetCookie(w, &cookie)
	return state
}

func spotifyRequest(token *oauth2.Token, url string) ([]byte, error) {
	client := spotifyConf.Client(oauth2.NoContext, token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error in Creating Request")
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in Performing Request")
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode == 400 {
		log.Println("Bad Request Syntax")
		log.Println(string(body))
		return body, err
	}
	if res.StatusCode == 401 {
		log.Println("Unauthorized Request")
		log.Println(string(body))
		return body, err
	}
	if err != nil {
		return body, err
	}
	return body, err

}
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func getToken(usrID string) (*oauth2.Token, error) {
	userCollection := mongoClient.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	res := userCollection.FindOne(ctx, bson.M{"spotifyid": usrID})
	if res != nil {
		token := new(oauth2.Token)
		var user User
		err := res.Decode(&user)
		if err != nil {
			log.Println(err)
			return token, err
		}
		token.AccessToken = user.AccessToken
		token.RefreshToken = user.RefreshToken
		token.Expiry = user.TokenExpiry
		token.TokenType = user.TokenType
		return token, err
	} else {
		return nil, errors.New("User Not Found")
	}
}
