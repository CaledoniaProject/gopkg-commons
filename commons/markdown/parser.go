package markdown

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

func NewMarkdownHeader(md string) (*MarkdownHeader, error) {
	var (
		header MarkdownHeader
	)

	parts := strings.SplitN(md, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("no valid front matter")
	}

	if err := yaml.Unmarshal([]byte(strings.TrimSpace(parts[1])), &header); err != nil {
		return nil, err
	} else {
		return &header, nil
	}
}
