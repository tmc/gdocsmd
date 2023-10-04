package convert

import (
	"fmt"
	"strings"

	"google.golang.org/api/docs/v1"
)

type MarkdownConverter struct {
	StylesToPrefix map[string]string
}

func NewMarkdownConverter() *MarkdownConverter {
	return &MarkdownConverter{
		StylesToPrefix: map[string]string{
			"NORMAL_TEXT": "",
			"TITLE":       "# ",
			"SUBTITLE":    "## ",
			"HEADING_1":   "# ",
			"HEADING_2":   "## ",
			"HEADING_3":   "### ",
			"HEADING_4":   "#### ",
			"HEADING_5":   "##### ",
			"HEADING_6":   "###### ",
		},
	}
}

func (mc *MarkdownConverter) AsMarkdown(doc *docs.Document) ([]byte, error) {
	var md []string
	md = append(md, fmt.Sprintf("# %s", doc.Title))
	for _, s := range doc.Body.Content {
		if s.Paragraph != nil {
			md = append(md, mc.paragraphAsMarkdown(s.Paragraph)...)
		} else if s.Table != nil {
			md = append(md, mc.tableAsMarkdown(s.Table)...)
		}
	}
	return []byte(strings.Join(md, "\n")), nil
}

func (mc *MarkdownConverter) paragraphAsMarkdown(p *docs.Paragraph) []string {
	var md []string
	prefix := mc.StylesToPrefix[p.ParagraphStyle.NamedStyleType]
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
		md = append(md, "\n")
	}
	return md
}

func (mc *MarkdownConverter) tableAsMarkdown(table *docs.Table) []string {
	var md []string
	for i, row := range table.TableRows {
		var rowMd []string
		for _, cell := range row.TableCells {
			// process each cell's content
			var cellMd []string
			for _, content := range cell.Content {
				if content.Paragraph != nil {
					cellMd = append(cellMd, strings.Join(mc.paragraphAsMarkdown(content.Paragraph), " "))
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
