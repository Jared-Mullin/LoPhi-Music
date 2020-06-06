package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	srv := httptest.NewServer(LoPhiRouter())
	client := srv.Client()
	t.Run("/spotify/auth redirects the user to Spotify", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/auth", nil)
		if err != nil {
			log.Println(err)
		}
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		expect := "https://accounts.spotify.com/"
		receive := res.Request.URL.String()
		if !strings.HasPrefix(receive, expect) {
			t.Errorf("Incorrect URL! Redirected to %q, Expected %q", receive, expect)
		}
	})
	t.Run("/spotify/callback requires state and code URL parameters", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/callback", nil)
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(400)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	t.Run("/spotify/artists is forbidden unless JWT is present", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/artists", nil)
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(401)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	t.Run("/spotify/tracks is forbidden unless JWT is present", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/tracks", nil)
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(401)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	t.Run("/spotify/genres is forbidden unless JWT is present", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/genres", nil)
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(401)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	t.Run("/spotify/artists should not accept a bad JWT", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/artists", nil)
		req.Header.Add("Authorization", "Bearer eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(401)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	t.Run("/spotify/tracks should not accept a bad JWT", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/tracks", nil)
		req.Header.Add("Authorization", "Bearer eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(401)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	t.Run("/spotify/genres should not accept a bad JWT", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/genres", nil)
		req.Header.Add("Authorization", "Bearer eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
		if err != nil {
			t.Errorf("Error Creating Request\n%q", err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error Doing Request\n%q", err)
		}
		expect := strconv.Itoa(401)
		receive := strconv.Itoa(res.StatusCode)
		if receive != expect {
			t.Errorf("Expected %q, Received %q", expect, receive)
		}
	})
	srv.Close()
}
