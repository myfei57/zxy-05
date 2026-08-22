package store

import "time"

// State machine constants shared across the control plane.
const (
	AlarmOpen           = "open"
	AlarmAcknowledged   = "acknowledged"
	AlarmResolved       = "resolved"
	AlarmClosed         = "closed"
	RuleDraft           = "draft"
	RulePublished       = "published"
	BatchPending        = "pending"
	BatchCommitted      = "committed"
	EscalationPending   = "pending"
	EscalationDone      = "done"
	NotificationPending = "pending"
	NotificationSent    = "sent"
	NotificationFailed  = "failed"
	NotificationAcked   = "acked"
)

// Device is a physical industrial device that reports telemetry.
type Device struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
}

// Point is a telemetry point owned by a device.
type Point struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	Enabled   bool      `json:"enabled"`
	LastSeq   int64     `json:"last_seq"`
	LastValue float64   `json:"last_value"`
	LastAt    time.Time `json:"last_at"`
}

// Sample is one normalized telemetry reading.
type Sample struct {
	PointID   string    `json:"point_id"`
	PointName string    `json:"point_name"`
	DeviceID  string    `json:"device_id"`
	Seq       int64     `json:"seq"`
	Value     float64   `json:"value"`
	At        time.Time `json:"at"`
}

// Rule defines a threshold check on a point name.
type Rule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PointName string    `json:"point_name"`
	Op        string    `json:"op"`
	Threshold float64   `json:"threshold"`
	Level     int       `json:"level"`
	State     string    `json:"state"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// RuleSnapshotEntry is one immutable rule inside a published snapshot.
type RuleSnapshotEntry struct {
	RuleID    string  `json:"rule_id"`
	PointName string  `json:"point_name"`
	Op        string  `json:"op"`
	Threshold float64 `json:"threshold"`
	Level     int     `json:"level"`
	State     string  `json:"state"`
}

// RuleVersion is one published or draft version book entry.
type RuleVersion struct {
	Version   int64     `json:"version"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

// Alarm is the alarm lifecycle record.
type Alarm struct {
	ID              string    `json:"id"`
	PointID         string    `json:"point_id"`
	PointName       string    `json:"point_name"`
	DeviceID        string    `json:"device_id"`
	Status          string    `json:"status"`
	OriginalOpenAt  time.Time `json:"original_open_at"`
	LastTriggerAt   time.Time `json:"last_trigger_at"`
	LastValue       float64   `json:"last_value"`
	Owner           string    `json:"owner"`
	Reason          string    `json:"reason"`
	AckedAt         time.Time `json:"acked_at"`
	ResolvedAt      time.Time `json:"resolved_at"`
	EscalationLevel int       `json:"escalation_level"`
}

// AlarmOwner is the durable acknowledgement attribution record.
type AlarmOwner struct {
	AlarmID string    `json:"alarm_id"`
	Owner   string    `json:"owner"`
	Reason  string    `json:"reason"`
	At      time.Time `json:"at"`
}

// Escalation is the SLA timer state for one alarm.
type Escalation struct {
	AlarmID string    `json:"alarm_id"`
	Level   int       `json:"level"`
	DueAt   time.Time `json:"due_at"`
	State   string    `json:"state"`
}

// Notification is one delivery attempt to a channel.
type Notification struct {
	ID        string    `json:"id"`
	AlarmID   string    `json:"alarm_id"`
	ChannelID string    `json:"channel_id"`
	Status    string    `json:"status"`
	Attempts  int       `json:"attempts"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
}

// ChannelConfig is a notification channel with reachability and priority.
type ChannelConfig struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"enabled"`
	Reachable bool      `json:"reachable"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

// MaintenanceWindow suppresses alarms for a device during a time span.
type MaintenanceWindow struct {
	ID       string    `json:"id"`
	DeviceID string    `json:"device_id"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

// EvalBatch is one scheduled evaluation batch.
type EvalBatch struct {
	ID          string    `json:"id"`
	Cursor      int64     `json:"cursor"`
	State       string    `json:"state"`
	PointCount  int       `json:"point_count"`
	StartedAt   time.Time `json:"started_at"`
	CommittedAt time.Time `json:"committed_at"`
}

// AuditEntry records one control-plane operation.
type AuditEntry struct {
	ID     string    `json:"id"`
	Action string    `json:"action"`
	Target string    `json:"target"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}
