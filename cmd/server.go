package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wndhydrnt/runbooks/pkg/api"
	"github.com/wndhydrnt/runbooks/pkg/parser"
	"github.com/wndhydrnt/runbooks/pkg/store"
	"github.com/wndhydrnt/runbooks/pkg/ui"
)

var (
	address = flag.String("address", ":8090", "address of the server.")
	uiURL   = flag.String("ui.url", "", "this url will be used to construct the runbook_url annotation. an empty string disables this functionality")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()
	p := parser.NewParser(*uiURL)
	store := store.NewInMemory()
	api.InitRoutesV0(router.PathPrefix("/api").Subrouter(), p, store)
	err := ui.InitRoutes(router, store)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(*address, router)
}
