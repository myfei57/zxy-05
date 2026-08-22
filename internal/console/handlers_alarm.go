package console

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/audit"
	"telemetryguard/internal/escalate"
	"telemetryguard/internal/notify"
)

func (a *API) ListAlarms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Alarms())
}

func (a *API) GetAlarm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	al, ok := a.state.Alarm(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	owner, _ := escalate.OwnerFor(a.state, id)
	writeJSON(w, http.StatusOK, map[string]any{"alarm": al, "owner": owner, "active": alarm.IsActive(al)})
}

func (a *API) AckAlarm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Owner  string `json:"owner"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	al, ok := a.state.Alarm(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	if !alarm.CanAcknowledge(al) {
		writeErr(w, http.StatusConflict, alarm.ErrAlarmNotOpen)
		return
	}
	if err := alarm.Acknowledge(a.state, id, body.Owner, body.Reason, time.Now().UTC()); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "acked"})
}

func (a *API) ResolveAlarm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	al, ok := a.state.Alarm(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	if !alarm.CanResolve(al) {
		writeErr(w, http.StatusConflict, alarm.ErrAlarmNotActive)
		return
	}
	if err := alarm.Resolve(a.state, id, time.Now().UTC()); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "resolved"})
}

func (a *API) CloseAlarm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := alarm.Close(a.state, id, time.Now().UTC()); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "closed"})
}

func (a *API) Stats(w http.ResponseWriter, r *http.Request) {
	devices, points, disabled := pointCounts(a.state)
	writeJSON(w, http.StatusOK, map[string]any{
		"active_alarms":     audit.CountActive(a.state),
		"devices":           devices,
		"points":            points,
		"disabled_devices":  disabled,
		"rules":             len(a.state.Rules()),
		"effective_version": a.state.Effective(),
		"pending_samples":   ingestPendingCount(),
		"cursor":            a.state.Cursor(),
		"failed_notify":     len(notifyFailed(a.state)),
	})
}

func (a *API) ListNotifications(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.SortedNotifications())
}

func (a *API) GetNotification(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	n, ok := a.state.Notification(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (a *API) RetryNotify(w http.ResponseWriter, r *http.Request) {
	count := notify.RetryFailed(a.state, time.Now().UTC())
	writeJSON(w, http.StatusOK, map[string]int{"requeued": count})
}

func (a *API) ListEscalations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Escalations())
}

func (a *API) GetEscalation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e, ok := a.state.Escalation(id)
	if !ok {
		writeErr(w, http.StatusNotFound, errNotFound)
		return
	}
	writeJSON(w, http.StatusOK, e)
}
