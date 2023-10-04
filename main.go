package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.comt/tmc/gdocsmd/auth"
	"github.comt/tmc/gdocsmd/convert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type Options struct {
	GoogleDocID string
	MDFile      string
	TokenFile   string
	Direction   string
	Credentials string
}

type App struct {
	Client *docs.Service
}

func (a *App) Run(ctx context.Context, opts Options) error {
	doc, err := a.Client.Documents.Get(opts.GoogleDocID).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from document: %w", err)
	}

	switch opts.Direction {
	case "to-md":
		md, err := convert.NewMarkdownConverter().AsMarkdown(doc)
		if err != nil {
			return fmt.Errorf("unable to marshal md: %w", err)
		}
		if err := os.WriteFile(opts.MDFile, []byte(md), 0644); err != nil {
			return fmt.Errorf("unable to write to md file: %w", err)
		}
	case "to-doc":
		c, err := ioutil.ReadFile(opts.MDFile)
		if err != nil {
			return fmt.Errorf("unable to read md file: %w", err)
		}
		return convert.MarkdownToDoc(ctx, a.Client, doc, convert.NewMarkdownParser(), c)
	default:
		return fmt.Errorf("invalid direction: %s", opts.Direction)
	}
	return nil
}

func NewApp(ctx context.Context, opts Options) (*App, error) {
	b, err := os.ReadFile(opts.Credentials)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client := auth.GetClient(opts.TokenFile, config, auth.DefaultFileHandler{}, auth.DefaultInputReader{})

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Docs client: %w", err)
	}

	return &App{Client: srv}, nil
}

func main() {
	flagGoogleDocID := flag.String("doc", "", "Google Doc ID")
	flagMDFile := flag.String("md", "", "Markdown file")
	flagTokenFile := flag.String("token", "token.json", "Token file")
	flagDirection := flag.String("direction", "to-md", "Direction of conversion (to-md or to-doc)")
	flagCredentials := flag.String("credentials", "credentials.json", "Credentials file")

	flag.Parse()

	opts := Options{
		GoogleDocID: *flagGoogleDocID,
		MDFile:      *flagMDFile,
		TokenFile:   *flagTokenFile,
		Direction:   *flagDirection,
		Credentials: *flagCredentials,
	}

	ctx := context.Background()
	app, err := NewApp(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}
