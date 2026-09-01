package test

import (
	"fmt"
	"testing"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/builder"
)

func TestMDToHTML(t *testing.T) {
	filePath := "assets/index.md"

	dstPath := "output/index.html"
	err := builder.MDToHTML(filePath, dstPath)
	if err != nil {
		t.Fatal("fail to convert markdown to html, err", err)
	}

	fmt.Println("check the output: ", dstPath)
}
