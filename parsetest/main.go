package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	b, err := ioutil.ReadFile("../examples/runbook.md")
	if err != nil {
		log.Fatal(err)
	}

	reader := text.NewReader(b)
	gmd := goldmark.New()
	p := gmd.Parser()
	root := p.Parse(reader)

	walker := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// fmt.Println(n.Type())
		// fmt.Println(n.Kind())
		if n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			if heading.Level == 1 {
				headingText := n.Text(b)
				fmt.Printf("%s\n", headingText)
				text := n.FirstChild()
				val := text.Text(b)
				fmt.Printf("%s\n", val)
				fmt.Println("------------------")
			}

			if heading.Level == 2 {
				if heading.NextSibling().Kind() != ast.KindFencedCodeBlock {
					return ast.WalkContinue, nil
				}

				headingText := n.Text(b)
				fmt.Printf("%s\n", headingText)
				codeBlock := heading.NextSibling().(*ast.FencedCodeBlock)
				codeType := codeBlock.Info.Text(b)
				fmt.Printf("%s\n", codeType)
				buf := bytes.NewBufferString("")
				lines := codeBlock.Lines()
				for i := 0; i < lines.Len(); i++ {
					line := lines.At(i)
					buf.Write(line.Value(b))
				}

				fmt.Println(buf.String())

				fmt.Println("------------------")
			}
		}

		return ast.WalkContinue, nil
	}

	err = ast.Walk(root, walker)
	if err != nil {
		log.Fatal(err)
	}
}
