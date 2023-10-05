package convert

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/docs/v1"
)

func TestRoundtripMDToGDocToMD(t *testing.T) {
	tests := []struct {
		name     string
		docID    string
		markdown string
	}{
		{
			name:  "basics",
			docID: "1LoxqGRxAVCRDunypVhS_3RaGPf04AANVmf3RK0MA9d8",
			markdown: `
# Untitled Document

# This is a title

This is a paragraph
`,
		},
		{
			name:  "bullets",
			docID: "1LoxqGRxAVCRDunypVhS_3RaGPf04AANVmf3RK0MA9d8",
			markdown: `
# Untitled Document

* This is a bullet
* This is another bullet
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Running test: ", tt.name)
			// Create a TestingDocsService for this subtest
			testingClient := NewTestingDocsService(t)

			// Delete all document contents to start:
			fmt.Println("FIRST CLEAR:")
			clearDocument(t, testingClient, tt.docID)
			// Delete all document contents to start:

			// fmt.Println("SECOND CLEAR:")
			// clearDocument(t, testingClient, tt.docID)
			// time.Sleep(2 * time.Second) // Sleep to avoid rate limiting

			// Convert MD -> GDOC
			parser := NewMarkdownParser()
			gdoc := &docs.Document{
				DocumentId: tt.docID,
			} // Create a new or empty Google Doc representation
			md := strings.TrimSpace(tt.markdown)
			err := MarkdownToDoc(context.Background(), testingClient, parser, gdoc, []byte(md))
			if err != nil {
				t.Fatalf("Failed converting MD -> GDOC: %v", err)
			}

			doc, err := testingClient.realService.Documents.Get(tt.docID).Do()
			if err != nil {
				t.Fatalf("Failed getting document: %v", err)
			}

			// Convert GDOC -> MD
			mdConverter := NewMarkdownConverter()
			convertedMDBytes, err := mdConverter.AsMarkdown(doc)
			if err != nil {
				t.Fatalf("Failed converting GDOC -> MD: %v", err)
			}
			convertedMD := string(convertedMDBytes)

			if cmp.Diff(md, convertedMD) != "" {
				t.Fatalf(cmp.Diff(md, convertedMD))
			}
		})
	}
}

func clearDocument(t *testing.T, testingClient *TestingDocsService, docID string) {
	t.Helper()
	err := clearBullets(t, testingClient, docID)
	if err != nil {
		t.Fatalf("Failed clearing document bullets: %v", err)
	}

	err = clearParagraphs(t, testingClient, docID)
	if err != nil {
		t.Fatalf("Failed clearing document paragraphs: %v", err)
	}
	t.Logf("Cleared document %s", docID)
}

func clearBullets(t *testing.T, testingClient *TestingDocsService, docID string) error {
	doc, err := testingClient.realService.Documents.Get(docID).Do()
	if err != nil {
		t.Fatalf("Failed getting document: %v", err)
	}
	var requests []*docs.Request

	requests = nil
	// delete all bullets:
	// for i := len(doc.Body.Content) - 1; i >= 0; i-- {
	// 	elem := doc.Body.Content[i]
	// 	fmt.Println(jm(elem))
	// 	if elem.Paragraph == nil {
	// 		fmt.Printf("bullet: Skipping non-paragraph element: %v\n", elem)
	// 		continue
	// 	}
	// 	if elem.Paragraph.Bullet == nil {
	// 		fmt.Printf("bullet: Skipping non-bullet element: %v\n", elem)
	// 		continue
	// 	}
	// 	r := &docs.DeleteContentRange
	// 		Range: &docs.Range{
	// 			StartIndex:      elem.StartIndex,
	// 			EndIndex:        elem.EndIndex,
	// 			ForceSendFields: []string{"StartIndex"},
	// 		},
	// 	}
	// 	// if the last one, skip the newline
	// 	// if elem.EndIndex == doc.Body.Content[len(doc.Body.Content)-1].EndIndex {
	// 	// 	r.Range.EndIndex--
	// 	// }
	// 	if r.Range.StartIndex == r.Range.EndIndex {
	// 		fmt.Printf("Skipping empty paragraph: %v\n", elem)
	// 		continue
	// 	}
	// 	// Add a request to delete the paragraph
	// 	requests = append(requests, &docs.Request{
	// 		DeleteContentRange: r,
	// 	})
	// }
	requests = append(requests, &docs.Request{
		DeleteParagraphBullets: &docs.DeleteParagraphBulletsRequest{
			Range: &docs.Range{
				StartIndex: 1,
				EndIndex:   doc.Body.Content[len(doc.Body.Content)-1].EndIndex,
			},
		},
	})

	if len(requests) == 0 {
		t.Logf("No bullets to clear")
		return nil
	}
	_, err = testingClient.realService.Documents.BatchUpdate(docID, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()
	return err
}

func clearParagraphs(t *testing.T, testingClient *TestingDocsService, docID string) error {
	// First get the document
	doc, err := testingClient.realService.Documents.Get(docID).Do()
	if err != nil {
		t.Fatalf("Failed getting document: %v", err)
	}
	var requests []*docs.Request

	requests = nil
	doc, err = testingClient.realService.Documents.Get(docID).Do()
	if err != nil {
		t.Fatalf("Failed getting document: %v", err)
	}
	//spew.Dump(doc.Body.Content)
	// delete all paragraphs, in reverse order
	for i := len(doc.Body.Content) - 1; i >= 0; i-- {
		elem := doc.Body.Content[i]
		if elem.Paragraph == nil {
			fmt.Printf("Skipping non-paragraph element: %v\n", elem)
			continue
		}
		r := &docs.DeleteContentRangeRequest{
			Range: &docs.Range{
				StartIndex:      elem.StartIndex,
				EndIndex:        elem.EndIndex,
				ForceSendFields: []string{"StartIndex"},
			},
		}
		// if the last one, skip the newline
		if elem.EndIndex == doc.Body.Content[len(doc.Body.Content)-1].EndIndex {
			r.Range.EndIndex--
		}
		if r.Range.StartIndex == r.Range.EndIndex {
			fmt.Printf("Skipping empty paragraph: %v\n", elem)
			continue
		}
		// Add a request to delete the paragraph
		requests = append(requests, &docs.Request{
			DeleteContentRange: r,
		})
	}
	//spew.Dump(requests)

	if len(requests) == 0 {
		t.Logf("No paragraphs to clear")
		return nil

	}
	if _, err = testingClient.realService.Documents.BatchUpdate(docID, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do(); err != nil {
		t.Fatalf("Failed clearing document: %v", err)
	}
	return nil
}

func jm(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
