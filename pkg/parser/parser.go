package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var (
	actionsHeading = "actions"
	alertsHeading  = "alerts"
)

type Runbook struct {
	Actions  []Action
	Alerts   []Alert
	Name     string
	Markdown []byte
}

type Action struct {
	Data []byte
	Name string
	Type string
}

type Alert struct {
	Data []byte
	Name string
	Type string
}

type Parser struct {
	uiURL string
}

func (p *Parser) ParseRunbook(input []byte) (Runbook, error) {
	r := Runbook{
		Markdown: input,
	}

	reader := text.NewReader(input)
	gmd := goldmark.New()
	parser := gmd.Parser()
	root := parser.Parse(reader)
	runbookName, err := extractRunbookName(input, root)
	if err != nil {
		return Runbook{}, fmt.Errorf("parse runbook: %v", err)
	}

	r.Name = runbookName
	parseActions := false
	parseAlerts := false
	headingName := ""
	walker := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			if heading.Level == 2 {
				headingText := strings.ToLower(string(heading.Text(input)))
				switch headingText {
				case actionsHeading:
					parseActions = true
					parseAlerts = false
				case alertsHeading:
					parseActions = false
					parseAlerts = true
				default:
					parseActions = false
					parseAlerts = false
				}
			}
		}

		if n.Kind() == ast.KindHeading {
			if !parseActions && !parseAlerts {
				return ast.WalkContinue, nil
			}

			heading := n.(*ast.Heading)
			if heading.Level == 3 {
				headingName = string(heading.Text(input))
			}
		}

		if n.Kind() == ast.KindFencedCodeBlock {
			if !parseActions && !parseAlerts {
				return ast.WalkContinue, nil
			}

			codeBlock := n.(*ast.FencedCodeBlock)
			codeType := codeBlock.Info.Text(input)
			buf := &bytes.Buffer{}
			lines := codeBlock.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				buf.Write(line.Value(input))
			}

			if parseActions {
				r.Actions = append(r.Actions, Action{
					Data: buf.Bytes(),
					Name: headingName,
					Type: string(codeType),
				})
			}

			if parseAlerts {
				r.Alerts = append(r.Alerts, Alert{
					Data: buf.Bytes(),
					Name: headingName,
					Type: string(codeType),
				})
			}
		}

		return ast.WalkContinue, nil
	}

	err = ast.Walk(root, walker)
	if err != nil {
		return Runbook{}, err
	}

	return r, nil
}

func extractRunbookName(input []byte, root ast.Node) (string, error) {
	firstChild := root.FirstChild()
	if firstChild.Kind() != ast.KindHeading {
		return "", fmt.Errorf("runbook does not start with h1 heading")
	}

	return string(firstChild.Text(input)), nil
}

func NewParser(uiURL string) *Parser {
	return &Parser{uiURL: uiURL}
}
