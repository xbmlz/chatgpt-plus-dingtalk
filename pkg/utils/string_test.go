package utils

import (
	"testing"
)

func TestExtractStringBetween(t *testing.T) {

	content := "abc ```mermaid\n graph TD; \n A-->B['name']; \n A-->C[\"pic\"]; \n B-->D; \n C-->D; \n ``` def"

	start := "```mermaid\n"

	end := "```"

	ret := ExtractStringBetween(start, end, content)

	t.Errorf("ret: %s", ret)
}
