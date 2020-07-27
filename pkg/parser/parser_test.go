package parser

import (
	"fmt"
	"io/ioutil"
	"log"
)

func ExampleParser_ParseRunbook() {
	input, err := ioutil.ReadFile("fixtures/runbook.md")
	if err != nil {
		log.Fatal(err)
	}

	p := &Parser{uiURL: "http://localhost:8090"}
	runbook, err := p.ParseRunbook(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(runbook.Name)
	fmt.Println(runbook.Actions[0].Name)
	fmt.Println(runbook.Actions[0].Type)
	fmt.Printf("%s", runbook.Actions[0].Data)
	fmt.Println(runbook.Alerts[0].Name)
	fmt.Println(runbook.Alerts[0].Type)
	// Output:
	// MyService
	// Restart service
	// bash
	// systemctl restart myservice
	// Instance down
	// yaml
}
