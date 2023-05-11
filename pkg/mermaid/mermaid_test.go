package mermaid

import (
	"fmt"
	"testing"
)

func TestMermaid(t *testing.T) {
	content := `graph TD;
    A-->B['name'];
    A-->C["pic"];
    B-->D;
    C-->D;`

	fmt.Println(New().RenderAsPng(content))
}
