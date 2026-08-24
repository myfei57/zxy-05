package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

)

func (a *API) ListAudit(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.AuditEntries())
}

func (a *API) GetAuditEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for _, e := range a.state.AuditEntries() {
		if e.ID == id {
			writeJSON(w, http.StatusOK, e)
			return
		}
	}
	writeErr(w, http.StatusNotFound, errNotFound)
}

