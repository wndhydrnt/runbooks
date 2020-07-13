package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/wndhydrnt/runbooks/pkg/parser"
	"gopkg.in/yaml.v3"
)

type RunbookStore interface {
	Create(runbook *parser.Runbook) error
	Delete(name string) error
	List() ([]*parser.Runbook, error)
}

type RunbookV0 struct {
	Name     string `json:"name"`
	Markdown []byte `json:"markdown"`
}

type runbookHandler struct {
	parser *parser.Parser
	store  RunbookStore
}

func (h *runbookHandler) createRunbookV0(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rb, err := h.parser.ParseRunbook(b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.store.Create(&rb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type listRunbooksV0Response struct {
	Runbooks []RunbookV0 `json:"runbooks"`
}

func (h *runbookHandler) listRunbooksV0(w http.ResponseWriter, r *http.Request) {
	storedRunbooks, err := h.store.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	runbooks := []RunbookV0{}
	for _, srb := range storedRunbooks {
		runbooks = append(runbooks, RunbookV0{
			Name:     srb.Name,
			Markdown: srb.Markdown,
		})
	}

	payload, err := json.Marshal(runbooks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

type prometheusHandler struct {
	store RunbookStore
}

func (p *prometheusHandler) listPrometheusRulesV0(w http.ResponseWriter, r *http.Request) {
	groups := rulefmt.RuleGroups{Groups: []rulefmt.RuleGroup{}}
	runbooks, err := p.store.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, rb := range runbooks {
		groups.Groups = append(groups.Groups, rb.Rules)
	}

	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	err = enc.Encode(groups)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = enc.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func InitRoutesV0(r *mux.Router, p *parser.Parser, rs RunbookStore) *mux.Router {
	v0 := r.PathPrefix("/v0").Subrouter()

	rh := &runbookHandler{parser: p, store: rs}
	v0.HandleFunc("/runbooks", rh.createRunbookV0).Name("createRunbookV0").Methods("POST")
	v0.HandleFunc("/runbooks/{runbook}", nil).Name("deleteRunbookV0").Methods("DELETE")
	v0.HandleFunc("/runbooks", rh.listRunbooksV0).Name("listRunbooksV0").Methods("GET")

	ph := &prometheusHandler{store: rs}
	v0.HandleFunc("/prometheus/rules", ph.listPrometheusRulesV0).Name("listPrometheusRulesV0").Methods("GET")
	return r
}
