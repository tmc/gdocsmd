package convert

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

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
