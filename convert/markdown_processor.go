package convert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"google.golang.org/api/docs/v1"
)

var slowdown = true

func MarkdownToDoc(ctx context.Context, docsService DocumentService, parser MarkdownParser, gdoc *docs.Document, mdContent []byte) error {
	doc := parser.Parse(text.NewReader(mdContent))

	updates := []*docs.Request{}
	index := int64(1)
	styleStart := int64(1)
	printLastUpdate := func() {
		return
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
		fmt.Println("performing update", jmar(update))
		time.Sleep(2000 * time.Millisecond)
		_, err := docsService.DoBatchUpdate(gdoc.DocumentId, &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{update},
		})

		if err != nil {
			return fmt.Errorf("unable to perform update: %w", err)
		}
		return nil
	}
	addUpdate := func(update *docs.Request) {
		updates = append(updates, update)
		printLastUpdate()
	}
	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		fmt.Println("node", node.Kind(), entering)
		time.Sleep(150 * time.Millisecond)
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
							EndIndex:   index - 2,
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
			fmt.Println("text", entering)
			if entering {
				return ast.WalkContinue, nil
			}
			textLen := len(n.Segment.Value(mdContent))
			if textLen == 0 {
				return ast.WalkContinue, nil
			}
			addUpdate(&docs.Request{
				InsertText: &docs.InsertTextRequest{
					Text: string(n.Segment.Value(mdContent)),
					Location: &docs.Location{
						Index: int64(index),
					},
				},
			})
			index += int64(len(n.Segment.Value(mdContent)))
		case *ast.List:
			fmt.Println("list", entering)
			if entering {
				//addUpdate(addText("\n"))
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
		case *ast.TextBlock:
			fmt.Println("text block", entering)
			if entering {
				return ast.WalkContinue, nil
			}
		default:
			fmt.Printf("unknown: %s\n", n.Kind())
		}
		return ast.WalkContinue, nil
	})

	if slowdown {
		for _, update := range updates {
			if err := performUpdate(update); err != nil {
				return err
			}
		}
	} else {
		resp, err := docsService.DoBatchUpdate(gdoc.DocumentId, &docs.BatchUpdateDocumentRequest{
			Requests: updates,
		})
		if err != nil {
			return fmt.Errorf("Batch update failed: %w", err)
		}
		if resp.HTTPStatusCode != 200 {
			return fmt.Errorf("Batch update failed: %d - %v", resp.HTTPStatusCode, resp)
		}
	}
	return nil
}

func jmar(v interface{}) string {
	j, _ := json.MarshalIndent(v, "", "  ")
	return string(j)
}
