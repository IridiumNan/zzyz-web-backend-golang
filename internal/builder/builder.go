package builder

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed style.css
var htmlStyle string

type RawContentType int

const (
	htmlHead   string = "<!DOCTYPE html>\n<html lang=\"cn\">\n"
	tagHTMLEnd string = "</html>"

	tagBodyBegin string = "<body>\n"
	tagBodyEnd   string = "</body>\n"

	tagStyleBegin string = "<style>\n"
	tagStyleEnd   string = "</style>\n"

	tagHeadBegin string = "<head>\n"
	tagHeadEnd   string = "</head>\n"

	metaHTML = `<meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />`

	RawTypeMD   RawContentType = 0
	RawTypeHTML RawContentType = 1
	RawTypeDOCX RawContentType = 2
)

// MDToHTML Use pandoc convert markdown intput into html output file with build-in style
func MDToHTML(srcMDPath string, dstHTMLPath string) (err error) {
	if _, err = os.Stat(srcMDPath); os.IsNotExist(err) {
		return fmt.Errorf("error when check src markdown file: file not exist, err: %w", err)
	}

	cmd := exec.Command("pandoc", "-f", "markdown", "-t", "html", srcMDPath)
	var rawHTMLBuf bytes.Buffer

	cmd.Stdout = &rawHTMLBuf

	// utils.TextLogger.Error("execing md to html convert", "exec_cmd", cmd.String())
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("error when exec command: %s, err: %w", cmd.String(), err)
	}

	htmlOutput := htmlHead + tagHeadBegin + metaHTML + tagStyleBegin + htmlStyle + tagStyleEnd + tagHeadEnd + tagBodyBegin + rawHTMLBuf.String() + tagBodyEnd + tagHTMLEnd

	err = os.WriteFile(dstHTMLPath, []byte(htmlOutput), 0o644)
	if err != nil {
		return fmt.Errorf("error when write htmlOutput into dst file, dst file: %s, err: %w", dstHTMLPath, err)
	}

	return
}
