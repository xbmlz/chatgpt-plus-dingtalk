package mermaid

import "encoding/base64"

type Mermaid struct {
}

func New() *Mermaid {
	return &Mermaid{}
}

func (m *Mermaid) RenderAsPng(content string) string {
	bytes := []byte(content)
	encodeToString := base64.StdEncoding.EncodeToString(bytes)
	return "https://mermaid.ink/img/" + encodeToString
}
