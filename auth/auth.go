package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// TokenFileHandler provides methods to handle token file operations
type TokenFileHandler interface {
	ReadToken(file string) (*oauth2.Token, error)
	SaveToken(file string, token *oauth2.Token) error
}

// DefaultFileHandler is the default file handler
type DefaultFileHandler struct{}

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

func GetClient(filePath string, config *oauth2.Config, fileHandler TokenFileHandler) *http.Client {
	tok, err := fileHandler.ReadToken(filePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		err = fileHandler.SaveToken(filePath, tok)
		if err != nil {
			log.Fatalf("Unable to cache OAuth token: %v", err)
		}
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	callbackCh := make(chan string)

	// Start a temporary HTTP server on a random port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("Could not start HTTP server: %v", err)
	}
	defer listener.Close()

	callbackURL := fmt.Sprintf("http://%s/", listener.Addr().String())
	config.RedirectURL = callbackURL

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found in URL", http.StatusBadRequest)
			return
		}
		callbackCh <- code
		// Show pretty tailwind page
		w.Write([]byte(okToClosePage))
	})

	go http.Serve(listener, nil)

	// Generate OAuth2 URL and open in user's browser
	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	authCode := <-callbackCh
	// Stop the HTTP server gracefully
	listener.Close()

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

const okToClosePage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OK to close</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.2.19/tailwind.min.css">
    <style>.gradient-bg { background: linear-gradient(135deg, #34D399 0%, #059669 100%); }</style>
</head>
<body class="gradient-bg">
    <div class="flex items-center justify-center h-screen">
        <div class="bg-white rounded-xl shadow-2xl p-12">
            <div class="flex items-center justify-center mb-8">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-16 h-16 text-green-500">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                </svg>
                <span class="text-green-500 font-semibold text-3xl ml-4">Authenticated</span>
            </div>
            <h1 class="text-2xl font-semibold text-gray-700 mb-6 text-center">You can close this page now</h1>
            <p class="text-gray-500 text-center">You can close this page now and return to the terminal.</p>
        </div>
    </div>
</body>
</html>
`
