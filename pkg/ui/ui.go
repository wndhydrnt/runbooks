package ui

import (
	"bytes"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/hoisie/mustache"
	"github.com/wndhydrnt/runbooks/pkg/api"
	"github.com/wndhydrnt/runbooks/pkg/parser"
	"github.com/yuin/goldmark"
	gmparser "github.com/yuin/goldmark/parser"
)

var (
	getRunbookTpl   *mustache.Template
	layoutTpl       *mustache.Template
	listRunbooksTpl *mustache.Template
)

type handler struct {
	gm    goldmark.Markdown
	store api.RunbookStore
}

type runbook struct {
	Name string
}

type getRunbookData struct {
	RunbookHTML string
}

func (h *handler) getRunbook(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	var runbook *parser.Runbook
	storedRunbooks, err := h.store.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, srb := range storedRunbooks {
		if srb.Name == name {
			runbook = srb
		}
	}

	if runbook == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	err = h.gm.Convert(runbook.Markdown, buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload := getRunbookTpl.RenderInLayout(layoutTpl, getRunbookData{RunbookHTML: buf.String()})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

type listRunbooksData struct {
	Runbooks []runbook
}

func (h *handler) listRunbooks(w http.ResponseWriter, r *http.Request) {
	storedRunbooks, err := h.store.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sort.Slice(storedRunbooks, func(i, j int) bool { return storedRunbooks[i].Name < storedRunbooks[j].Name })
	data := listRunbooksData{}
	for _, srb := range storedRunbooks {
		data.Runbooks = append(data.Runbooks, runbook{Name: srb.Name})
	}

	payload := listRunbooksTpl.RenderInLayout(layoutTpl, data)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func InitRoutes(r *mux.Router, rs api.RunbookStore) error {
	var err error
	getRunbookTpl, err = mustache.ParseString(getRunbookTplString)
	if err != nil {
		return err
	}

	listRunbooksTpl, err = mustache.ParseString(listRunbooksTplString)
	if err != nil {
		return err
	}

	layoutTpl, err = mustache.ParseString(layoutTplString)
	if err != nil {
		return err
	}

	gm := goldmark.New(
		goldmark.WithParserOptions(
			gmparser.WithAutoHeadingID(),
		),
	)
	h := &handler{gm: gm, store: rs}
	r.HandleFunc("/runbooks", h.listRunbooks).Name("listRunbooks").Methods("GET")
	r.HandleFunc("/runbooks/{name}", h.getRunbook).Name("getRunbook").Methods("GET")
	r.Handle("/", http.RedirectHandler("/runbooks", http.StatusMovedPermanently))
	return nil
}
