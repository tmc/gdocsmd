package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// Options represents the command line options.
type Options struct {
	GoogleDocID string
	MDFile      string
	TokenFile   string
	Direction   string
	Credentials string
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

	if err := run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, opts Options) error {
	b, err := os.ReadFile(opts.Credentials)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents")
	// config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client := getClient(config)

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Docs client: %w", err)
	}

	doc, err := srv.Documents.Get(opts.GoogleDocID).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from document: %w", err)
	}

	switch opts.Direction {
	case "to-md":
		md, err := asMarkdown(doc)
		if err != nil {
			return fmt.Errorf("unable to marshal md: %w", err)
		}
		if err := os.WriteFile(opts.MDFile, []byte(md), 0644); err != nil {
			return fmt.Errorf("unable to write to md file: %w", err)
		}
	case "to-doc":
		return markdownToDoc(ctx, srv, doc, opts.MDFile)
	default:
		return fmt.Errorf("invalid direction: %s", opts.Direction)
	}
	return nil
}

func markdownToDoc(ctx context.Context, srv *docs.Service, gdoc *docs.Document, mdFile string) error {
	rawMarkdown, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return fmt.Errorf("unable to read md file: %w", err)
	}
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	doc := gm.Parser().Parse(text.NewReader(rawMarkdown))

	updates := []*docs.Request{}

	index := int64(1)

	styleStart := int64(1)
	printLastUpdate := func() {
		j, _ := json.Marshal(updates[len(updates)-1])
		fmt.Println(len(updates), string(j), index)
	}

	addText := func(text string) *docs.Request {
		update := &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Text: text,
				Location: &docs.Location{
					Index: index,
				},
			},
		}
		index += int64(len(text))

		return update
	}

	performUpdate := func(update *docs.Request) error {
		//time.Sleep(5 * time.Millisecond)
		_, err := srv.Documents.BatchUpdate(gdoc.DocumentId, &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{update},
		}).Do()
		if err != nil {
			return fmt.Errorf("unable to perform update: %w", err)
		}
		return nil
	}
	_ = performUpdate
	addUpdate := func(update *docs.Request) {
		updates = append(updates, update)
		printLastUpdate()
	}
	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n := node.(type) {
		case *ast.Document:
		case *ast.Heading:
			fmt.Println("heading", n.Level, entering)
			if entering {
				if index > 1 {
					addUpdate(addText("\n"))
				}
				styleStart = index
			} else {
				addUpdate(&docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType: "HEADING_" + fmt.Sprint(n.Level),
						},
						Range: &docs.Range{
							StartIndex: styleStart,
							EndIndex:   index - 1,
						},
						Fields: "namedStyleType",
					},
				})
			}
		case *ast.Paragraph:
			fmt.Println("paragraph", entering)
			if entering {
				addUpdate(addText("\n"))
				styleStart = index
			} else {
				addUpdate(&docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType: "NORMAL_TEXT",
						},
						Range: &docs.Range{
							StartIndex: styleStart,
							EndIndex:   index - 1,
						},
						Fields: "namedStyleType",
					},
				})
			}
		case *ast.Text:
			//fmt.Println("text", entering)
			if entering {
				return ast.WalkContinue, nil
			}
			textLen := len(n.Segment.Value(rawMarkdown))
			if textLen == 0 {
				return ast.WalkContinue, nil
			}
			addUpdate(&docs.Request{
				InsertText: &docs.InsertTextRequest{
					Text: string(n.Segment.Value(rawMarkdown)),
					Location: &docs.Location{
						Index: int64(index),
					},
				},
			})
			index += int64(len(n.Segment.Value(rawMarkdown)))
		case *ast.List:
			fmt.Println("list", entering)
			if entering {
				addUpdate(addText("\n"))
				styleStart = index
			} else {
				addUpdate(&docs.Request{
					CreateParagraphBullets: &docs.CreateParagraphBulletsRequest{
						BulletPreset: "BULLET_DISC_CIRCLE_SQUARE",
						Range: &docs.Range{
							StartIndex: styleStart,
							EndIndex:   index - 1,
						},
					},
				})
			}
		case *ast.ListItem:
			fmt.Println("list item", entering)
			if entering {
				addUpdate(addText("\n"))
			}
		case *ast.Link:
			fmt.Println("link", entering)
			if entering {
				return ast.WalkContinue, nil
			}
			addUpdate(&docs.Request{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					TextStyle: &docs.TextStyle{
						Link: &docs.Link{
							Url: string(n.Destination),
						},
					},
					Range: &docs.Range{
						StartIndex: index - int64(len(n.Destination)),
						EndIndex:   index - 1,
					},
					Fields: "link",
				},
			})
		default:
			fmt.Printf("unknown: %s\n", n.Kind())
		}
		return ast.WalkContinue, nil
	})

	// add title:
	// updates = append(updates, &docs.Request{
	// 	InsertText: &docs.InsertTextRequest{
	// 		Text: doc.Title,
	// 		Location: &docs.Location{
	// 			Index: 1,
	// 		},
	// 	},
	// })

	// // todo: use markdown parser
	// for i, p := range nonEmptyParagraphs[:50] {
	// 	updates = append(updates, &docs.Request{
	// 		InsertText: &docs.InsertTextRequest{
	// 			Text: p + "\n\n",
	// 			Location: &docs.Location{
	// 				Index: index,
	// 			},
	// 		},
	// 	})
	// 	// increment by number of utf16 code units:
	// 	index += int64(len([]rune(p)))
	// }

	// // send batch update request

	// apply each step with a delay in-between:
	/*
		for i, u := range updates {
			j, _ := json.Marshal(updates[len(updates)-1])
			fmt.Println(i, len(updates), string(j))
			_, err := srv.Documents.BatchUpdate(gdoc.DocumentId, &docs.BatchUpdateDocumentRequest{
				Requests: []*docs.Request{u},
			}).Do()
			if err != nil {
				return fmt.Errorf("unable to batch update: %w", err)
			}
		}
	*/

	resp, err := srv.Documents.BatchUpdate(gdoc.DocumentId, &docs.BatchUpdateDocumentRequest{
		Requests: updates,
	}).Do()
	fmt.Println(resp)

	return err
}

func asMarkdown(doc *docs.Document) ([]byte, error) {
	var md []string
	// add title:
	md = append(md, fmt.Sprintf("# %s", doc.Title))
	for _, s := range doc.Body.Content {
		if s.Paragraph != nil {
			md = append(md, paragraphAsMarkdown(s.Paragraph)...)
		} else if s.Table != nil {
			md = append(md, tableAsMarkdown(s.Table)...)
		}
	}
	return []byte(strings.Join(md, "\n")), nil
}

func paragraphAsMarkdown(p *docs.Paragraph) []string {
	var md []string
	prefix := paragraphStylesToPrefix[p.ParagraphStyle.NamedStyleType]
	isList := p.Bullet != nil
	for _, elem := range p.Elements {
		if elem.TextRun != nil {
			text := processTextRuns(elem.TextRun)
			if isList && len(text) > 0 {
				text = "* " + text
			}
			md = append(md, prefix+text)
		}
	}
	if p.Bullet == nil {
		md = append(md, "\n") // add newline after each paragraph
	}
	return md
}

func tableAsMarkdown(table *docs.Table) []string {
	var md []string
	for i, row := range table.TableRows {
		var rowMd []string
		for _, cell := range row.TableCells {
			// process each cell's content
			var cellMd []string
			for _, content := range cell.Content {
				if content.Paragraph != nil {
					cellMd = append(cellMd, strings.Join(paragraphAsMarkdown(content.Paragraph), " "))
					// trim trailing newline:
					cellMd[len(cellMd)-1] = strings.TrimRight(cellMd[len(cellMd)-1], "\n")
				}
			}
			rowMd = append(rowMd, strings.Join(cellMd, " "))
		}
		md = append(md, "| "+strings.Join(rowMd, " | ")+" |")
		if i == 0 {
			// after the first row (header), add a separator row
			md = append(md, "|"+strings.Repeat(" --- |", len(rowMd)))
		}
	}
	md = append(md, "\n") // add newline after each table
	return md
}

func processTextRuns(elements ...*docs.TextRun) string {
	var md []string
	for _, tr := range elements {
		text := tr.Content
		if tr.TextStyle.Bold {
			text = "**" + text + "**"
		}
		if tr.TextStyle.Italic {
			text = "*" + text + "*"
		}
		if tr.TextStyle.Strikethrough {
			text = "~~" + text + "~~"
		}
		if tr.TextStyle.Underline && tr.TextStyle.Link == nil {
			// Markdown doesn't support underline, but it can be represented using HTML
			text = "<u>" + text + "</u>"
		}
		if tr.TextStyle.Link != nil {
			text = "[" + text + "](" + tr.TextStyle.Link.Url + ")"
		}
		md = append(md, strings.TrimSpace(text))
	}
	return strings.Join(md, "")
}

var paragraphStylesToPrefix = map[string]string{
	"NORMAL_TEXT": "",
	"TITLE":       "# ",
	"SUBTITLE":    "## ",
	"HEADING_1":   "# ",
	"HEADING_2":   "## ",
	"HEADING_3":   "### ",
	"HEADING_4":   "#### ",
	"HEADING_5":   "##### ",
	"HEADING_6":   "###### ",
}

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}
