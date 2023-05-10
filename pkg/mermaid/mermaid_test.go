package mermaid

import (
	"context"
	"fmt"
	"testing"

	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/utils"
)

func TestMermaid(t *testing.T) {
	content := `graph TD;
    A-->B['name'];
    A-->C["pic"];
    B-->D;
    C-->D;`
	ctx1 := context.Background()
	re1, _ := NewRenderEngine(ctx1, `mermaid.initialize({'theme': 'base', 'themeVariables': { 'primaryColor': '#1473e6'}});`)
	defer re1.Cancel()
	result_in_bytes, box, err := re1.RenderAsPng(content)

	if err != nil {
		fmt.Println(err)
		return
	}
	if box.Width < 1 || box.Height < 1 {
		fmt.Printf("Render() got empty image = w:%d, h:%d)", box.Width, box.Height)
		return
	}

	utils.ImageSave(result_in_bytes, "D:\\imgAAA\\test.png")

}
