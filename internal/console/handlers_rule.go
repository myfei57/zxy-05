package console

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"telemetryguard/internal/rule"
)

func (a *API) ListRules(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"rules":    a.state.Rules(),
		"versions": a.state.Versions(),
		"effective": a.state.Effective(),
	})
}

func (a *API) GetRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rl, ok := a.state.Rule(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, rl)
}

func (a *API) CreateRule(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name      string  `json:"name"`
		PointName string  `json:"point_name"`
		Op        string  `json:"op"`
		Threshold float64 `json:"threshold"`
		Level     int     `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	rl, err := rule.CreateRule(a.state, body.Name, body.PointName, body.Op, body.Threshold, body.Level)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, rl)
}

func (a *API) PublishRules(w http.ResponseWriter, r *http.Request) {
	version, err := rule.PublishVersion(a.state)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"version": version})
}

func (a *API) RollbackRules(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Version int64 `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := rule.RollbackTo(a.state, body.Version); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"effective": rule.EffectiveVersion(a.state)})
}

func (a *API) GetVersion(w http.ResponseWriter, r *http.Request) {
	version, err := strconv.ParseInt(chi.URLParam(r, "version"), 10, 64)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	v, ok := a.state.Version(version)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"version":  v,
		"complete": rule.SnapshotComplete(a.state, version),
	})
}

