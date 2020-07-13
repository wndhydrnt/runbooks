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
	runbook, err := p.ParseRunbook([]byte(input))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(runbook.Name)
	fmt.Println(runbook.Rules.Rules[0].Alert.Value)
	fmt.Println(runbook.Rules.Rules[0].Expr.Value)
	fmt.Println(runbook.Rules.Rules[0].Annotations["runbook_url"])
	// Output:
	// MyService
	// InstanceDown
	// up{job="MyService"} == 0
	// http://localhost:8090/runbooks/MyService#instance-down
}
