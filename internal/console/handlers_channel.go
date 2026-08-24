package console

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"telemetryguard/internal/channel"
	"telemetryguard/internal/schedule"
)

func (a *API) ListChannels(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"channels": a.state.Channels(),
		"route":    channel.CachedRoute(a.state),
	})
}

func (a *API) GetChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, ok := a.state.Channel(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (a *API) RegisterChannel(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Priority int    `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	c, err := channel.RegisterChannel(a.state, body.Name, body.Priority)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (a *API) FailoverChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := channel.Failover(a.state, id); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	route, _ := channel.CurrentRoute(a.state)
	writeJSON(w, http.StatusOK, map[string]string{"route": route})
}

func (a *API) CurrentRoute(w http.ResponseWriter, r *http.Request) {
	route, err := channel.CurrentRoute(a.state)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"route": route})
}

func (a *API) AddWindow(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DeviceID string `json:"device_id"`
		Start    string `json:"start"`
		End      string `json:"end"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	start, err1 := parseTime(body.Start)
	end, err2 := parseTime(body.End)
	if err1 != nil || err2 != nil {
		writeErr(w, http.StatusBadRequest, errBadTime)
		return
	}
	win, err := schedule.AddWindow(a.state, body.DeviceID, start, end)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, win)
}

func (a *API) ListWindows(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Windows())
}

func (a *API) ListBatches(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Batches())
}

func (a *API) GetBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	b, ok := a.state.Batch(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

