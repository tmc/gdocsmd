// service.go

package convert

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"google.golang.org/api/docs/v1"
)

type DocumentService interface {
	BatchUpdate(string, *docs.BatchUpdateDocumentRequest) *docs.DocumentsBatchUpdateCall
}

type MarkdownParser interface {
	Parse(text.Reader, ...parser.ParseOption) ast.Node
}
