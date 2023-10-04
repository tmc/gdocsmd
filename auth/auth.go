package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// TokenFileHandler provides methods to handle token file operations
type TokenFileHandler interface {
	ReadToken(file string) (*oauth2.Token, error)
	SaveToken(file string, token *oauth2.Token) error
}

// InputReader provides a method to read input
type InputReader interface {
	ReadInput() (string, error)
}

// DefaultFileHandler is the default file handler
type DefaultFileHandler struct{}

// DefaultInputReader is the default input reader
type DefaultInputReader struct{}

func (DefaultFileHandler) ReadToken(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (DefaultFileHandler) SaveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(token)
}

func (DefaultInputReader) ReadInput() (string, error) {
	var input string
	if _, err := fmt.Scan(&input); err != nil {
		return "", err
	}
	return input, nil
}

func GetClient(filePath string, config *oauth2.Config, fileHandler TokenFileHandler, inputReader InputReader) *http.Client {
	tok, err := fileHandler.ReadToken(filePath)
	if err != nil {
		tok = getTokenFromWeb(config, inputReader)
		err = fileHandler.SaveToken(filePath, tok)
		if err != nil {
			log.Fatalf("Unable to cache OAuth token: %v", err)
		}
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config, inputReader InputReader) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	authCode, err := inputReader.ReadInput()
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
