package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/joeshaw/envdecode"
	"golang.org/x/oauth2"
)

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}

	conn = netc
	return netc, nil
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}

	if reader != nil {
		reader.Close()
	}
}

var (
	authClient *oauth2.Config
	creds      *oauth2.Token
)

func setupXAuth() {
	var x struct { //stores env variables
		ConsumerKey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AccessToken    string `env:"SP_TWITTER_ACCESSTOKEN,required"`
		AccessSecret   string `env:"SP_TWITTER_ACCESSSECRET,required"`
	}
	if err := envdecode.Decode(&x); err != nil {
		log.Fatalln(err)
	}

	creds = &oauth2.Token{
		AccessToken:  x.AccessToken,
		RefreshToken: x.AccessSecret,
	}

	authClient = &oauth2.Config{
		ClientID:     x.ConsumerKey,
		ClientSecret: x.ConsumerSecret,
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupXAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			},
		}
	})

	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form- urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
	
	return httpClient.Do(req)
}
