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
	address       = flag.String("address", ":8090", "address of the server.")
	uiURL         = flag.String("ui.url", "", "this url will be used to construct the runbook_url annotation. an empty string disables this functionality")
	storeType     = flag.String("store.type", "memory", "the type of store to use (valid: file, memory)")
	storeFilePath = flag.String("store.file.path", "*.md", "a glob pattern to use to discover runbooks")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()
	p := parser.NewParser(*uiURL)

	var s api.RunbookStore
	switch *storeType {
	case "file":
		var err error
		s, err = store.NewFile(*storeFilePath, p)
		if err != nil {
			log.Fatalf("create file store: %s", err)
		}
	case "memory":
		s = store.NewInMemory()
	default:
		log.Fatalf("unknown store of type %s", *storeType)
	}

	api.InitRoutesV0(router.PathPrefix("/api").Subrouter(), p, s)
	err := ui.InitRoutes(router, s)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(*address, router)
}
