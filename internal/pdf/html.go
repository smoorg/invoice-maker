package pdf

import (
	"github.com/gomarkdown/markdown"
)

func ToHTML(md []byte) []byte {
	return markdown.ToHTML(md, nil, nil)
}
