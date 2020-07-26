package parser

import (
	"bytes"
	"fmt"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/shurcooL/sanitized_anchor_name"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

type Runbook struct {
	Name     string
	Markdown []byte
	Rules    rulefmt.RuleGroup
}

type Parser struct {
	uiURL string
}

func (p *Parser) ParseRunbook(input []byte) (Runbook, error) {
	r := Runbook{
		Markdown: input,
		Rules: rulefmt.RuleGroup{
			Interval: model.Duration(30 * time.Second),
		},
	}

	walker := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			if heading.Level == 1 {
				r.Name = string(n.Text(input))
				r.Rules.Name = string(n.Text(input))
			}

			if heading.Level == 2 {
				if heading.NextSibling().Kind() != ast.KindFencedCodeBlock {
					return ast.WalkContinue, nil
				}

				codeBlock := heading.NextSibling().(*ast.FencedCodeBlock)
				codeType := codeBlock.Info.Text(input)
				if string(codeType) != "yaml" {
					return ast.WalkContinue, nil
				}

				buf := bytes.NewBufferString("")
				lines := codeBlock.Lines()
				for i := 0; i < lines.Len(); i++ {
					line := lines.At(i)
					buf.Write(line.Value(input))
				}

				rule := rulefmt.RuleNode{}
				err := yaml.Unmarshal(buf.Bytes(), &rule)
				if err != nil {
					return ast.WalkStop, err
				}

				_, exists := rule.Annotations["runbook_url"]
				if !exists && p.uiURL != "" {
					headingText := heading.Text(input)
					anchor := sanitized_anchor_name.Create(string(headingText))
					rule.Annotations["runbook_url"] = fmt.Sprintf("%s/runbooks/%s#%s", p.uiURL, r.Name, anchor)
				}

				r.Rules.Rules = append(r.Rules.Rules, rule)
			}
		}

		return ast.WalkContinue, nil
	}

	reader := text.NewReader(input)
	gmd := goldmark.New()
	parser := gmd.Parser()
	root := parser.Parse(reader)
	err := ast.Walk(root, walker)
	if err != nil {
		return Runbook{}, err
	}

	return r, nil
}

func NewParser(uiURL string) *Parser {
	return &Parser{uiURL: uiURL}
}
