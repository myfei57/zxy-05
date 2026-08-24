package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"telemetryguard/internal/store"
)

// API exposes the telemetry alarm control plane over HTTP.
type API struct {
	state *store.State
}

// NewAPI builds the HTTP handler for the control plane.
func NewAPI(state *store.State) http.Handler {
	api := &API{state: state}
	r := chi.NewRouter()
	r.Get("/", api.IndexPage)
	r.Get("/console/alarms", api.AlarmsPage)
	r.Get("/console/rules", api.RulesPage)
	r.Get("/console/points", api.PointsPage)
	r.Get("/console/audit", api.AuditPage)

	r.Get("/api/stats", api.Stats)
	r.Get("/api/alarms", api.ListAlarms)
	r.Get("/api/alarms/{id}", api.GetAlarm)
	r.Post("/api/alarms/{id}/ack", api.AckAlarm)
	r.Post("/api/alarms/{id}/resolve", api.ResolveAlarm)
	r.Post("/api/alarms/{id}/close", api.CloseAlarm)

	r.Get("/api/rules", api.ListRules)
	r.Get("/api/rules/{id}", api.GetRule)
	r.Post("/api/rules", api.CreateRule)
	r.Post("/api/rules/publish", api.PublishRules)
	r.Post("/api/rules/rollback", api.RollbackRules)

	r.Get("/api/points", api.ListPoints)
	r.Post("/api/points", api.RegisterPoint)
	r.Get("/api/points/{id}", api.GetPoint)
	r.Get("/api/points/enabled", api.EnabledPoints)
	r.Post("/api/devices", api.RegisterDevice)
	r.Get("/api/devices/{id}", api.GetDevice)
	r.Post("/api/devices/{id}/disable", api.DisableDevice)
	r.Post("/api/devices/{id}/enable", api.EnableDevice)
	r.Post("/api/points/{id}/sample", api.SamplePoint)

	r.Get("/api/channels", api.ListChannels)
	r.Get("/api/channels/{id}", api.GetChannel)
	r.Post("/api/channels", api.RegisterChannel)
	r.Post("/api/channels/{id}/failover", api.FailoverChannel)
	r.Get("/api/channels/route/current", api.CurrentRoute)

	r.Get("/api/notifications", api.ListNotifications)
	r.Get("/api/notifications/{id}", api.GetNotification)
	r.Post("/api/notify/retry", api.RetryNotify)

	r.Post("/api/windows", api.AddWindow)
	r.Get("/api/windows", api.ListWindows)

	r.Get("/api/batches", api.ListBatches)
	r.Get("/api/batches/{id}", api.GetBatch)
	r.Get("/api/escalations", api.ListEscalations)
	r.Get("/api/escalations/{id}", api.GetEscalation)

	r.Get("/api/audit", api.ListAudit)
	r.Get("/api/audit/{id}", api.GetAuditEntry)
	r.Get("/api/versions/{version}", api.GetVersion)
	return r
}
