package console

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"telemetryguard/internal/point"
	"telemetryguard/internal/cycle"
)

func (a *API) ListPoints(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"devices": a.state.Devices(),
		"points":  a.state.Points(),
	})
}

func (a *API) GetPoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, ok := a.state.Point(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (a *API) RegisterPoint(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DeviceID string `json:"device_id"`
		Name     string `json:"name"`
		Unit     string `json:"unit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if _, exists := a.state.PointByName(body.Name); exists {
		writeErr(w, http.StatusConflict, errDuplicatePoint)
		return
	}
	pt, err := point.RegisterPoint(a.state, body.DeviceID, body.Name, body.Unit)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, pt)
}

func (a *API) EnabledPoints(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, point.EnabledPoints(a.state))
}

func (a *API) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	d, err := point.RegisterDevice(a.state, body.Name)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, d)
}

func (a *API) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	d, ok := a.state.Device(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (a *API) DisableDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := point.DisableDevice(a.state, id); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "disabled"})
}

func (a *API) EnableDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := point.EnableDevice(a.state, id); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "enabled"})
}

func (a *API) SamplePoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Seq   int64   `json:"seq"`
		Value float64 `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	opened, err := cycle.HandleSample(a.state, id, body.Seq, body.Value, time.Now().UTC())
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"opened": len(opened)})
}

