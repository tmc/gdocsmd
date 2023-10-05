// service.go

package convert

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"google.golang.org/api/docs/v1"
)

type DocumentService interface {
	DoBatchUpdate(string, *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error)
}

type RealDocumentService struct {
	*docs.Service
}

func (r *RealDocumentService) DoBatchUpdate(documentId string, request *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error) {
	return r.Documents.BatchUpdate(documentId, request).Do()
}

type MarkdownParser interface {
	Parse(text.Reader, ...parser.ParseOption) ast.Node
}
