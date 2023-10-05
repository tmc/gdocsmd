package convert

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.comt/tmc/gdocsmd/auth"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func SimplifyMarkdown(md string) string {
	// Handle multiple spaces, newlines, etc.
	// Add more simplifications as needed.
	md = strings.ReplaceAll(md, "\n\n", "\n")
	md = strings.TrimSpace(md)
	return md
}

func CompareMarkdowns(md1, md2 string) bool {
	return cmp.Equal(SimplifyMarkdown(md1), SimplifyMarkdown(md2))
}

type TestingDocsService struct {
	realService *docs.Service
	t           *testing.T
}

var _ DocumentService = (*TestingDocsService)(nil)

func NewTestingDocsService(t *testing.T) *TestingDocsService {
	t.Helper()
	s := &TestingDocsService{
		t: t,
	}
	if *updateGolden {
		s.realService = NewRealDocsService(t)
	}
	return s
}

func (tds *TestingDocsService) DoBatchUpdate(documentId string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error) {
	tds.t.Helper()
	// Ensure the testdata directory exists
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		os.Mkdir("testdata", 0755)
	}
	// Filename for the golden response based on the test name
	filename := fmt.Sprintf("testdata/golden_response_%s.json", strings.ReplaceAll(tds.t.Name(), "/", "_"))

	// If not updating golden, attempt to use the saved response
	if !*updateGolden {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			tds.t.Fatalf("Failed to read golden response: %v", err)
		}

		var resp docs.BatchUpdateDocumentResponse
		err = json.Unmarshal(data, &resp)
		if err != nil {
			tds.t.Fatalf("Failed to unmarshal golden response: %v", err)
		}

		return &resp, nil
	}

	// Otherwise, make the actual API call
	resp, err := tds.realService.Documents.BatchUpdate(documentId, req).Do()
	if err != nil {
		return nil, err
	}

	// If updating golden files, save the new response
	data, err := resp.MarshalJSON()
	if err != nil {
		tds.t.Fatalf("Failed to marshal response: %v", err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		tds.t.Fatalf("Failed to save golden response: %v", err)
	}

	return resp, nil
}

func NewRealDocsService(t *testing.T) *docs.Service {
	t.Helper()
	srv, err := newRealDocsService()
	if err != nil {
		t.Fatalf("Failed to create real Docs service: %v", err)
	}
	return srv
}

func newRealDocsService() (*docs.Service, error) {
	b, err := os.ReadFile("../client-secret.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, docs.DocumentsScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client := auth.GetClient("../token.json", config, auth.DefaultFileHandler{})

	srv, err := docs.NewService(context.Background(), option.WithHTTPClient(client), option.WithScopes(docs.DocumentsScope))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Docs client: %w", err)
	}
	return srv, nil
}
