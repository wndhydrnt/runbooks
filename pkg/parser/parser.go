package parser

import (
	"fmt"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/russross/blackfriday/v2"
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
	var lastErr error
	visitor := func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		}

		switch node.Type {
		case blackfriday.Heading:
			if node.HeadingData.Level == 1 {
				if r.Name != "" {
					return blackfriday.SkipChildren
				}

				r.Name = string(node.FirstChild.Literal)
				r.Rules.Name = string(node.FirstChild.Literal)
				return blackfriday.SkipChildren
			}

			if node.HeadingData.Level == 2 {
				if node.Next != nil && node.Next.Type == blackfriday.CodeBlock && string(node.Next.CodeBlockData.Info) == "yaml" {
					rule := rulefmt.RuleNode{}
					err := yaml.Unmarshal(node.Next.Literal, &rule)
					if err != nil {
						lastErr = err
						return blackfriday.Terminate
					}

					_, exists := rule.Annotations["runbook_url"]
					if !exists && p.uiURL != "" {
						anchor := blackfriday.SanitizedAnchorName(string(node.FirstChild.Literal))
						rule.Annotations["runbook_url"] = fmt.Sprintf("%s/runbooks/%s#%s", p.uiURL, r.Name, anchor)
					}

					r.Rules.Rules = append(r.Rules.Rules, rule)
				}

				return blackfriday.GoToNext
			}
		}

		return blackfriday.GoToNext
	}

	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.FencedCode))
	rootNode := parser.Parse(input)
	rootNode.Walk(visitor)
	if lastErr != nil {
		return Runbook{}, lastErr
	}

	return r, nil
}

func NewParser(uiURL string) *Parser {
	return &Parser{uiURL: uiURL}
}
