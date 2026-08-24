package store

import (
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

// State is the central in-memory plus file-backed registry.
type State struct {
	mu            sync.RWMutex
	root          string
	devices       map[string]*Device
	points        map[string]*Point
	rules         map[string]*Rule
	versions      map[int64]*RuleVersion
	effective     int64
	alarms        map[string]*Alarm
	owners        map[string]*AlarmOwner
	escalations   map[string]*Escalation
	notifications map[string]*Notification
	channels      map[string]*ChannelConfig
	currentRoute  string
	windows       map[string]*MaintenanceWindow
	batches       map[string]*EvalBatch
	cursor        int64
	audit         map[string]*AuditEntry
}

// NewState creates a State rooted at dir and restores persisted registry files.
func NewState(dir string) *State {
	s := &State{
		root:          dir,
		devices:       make(map[string]*Device),
		points:        make(map[string]*Point),
		rules:         make(map[string]*Rule),
		versions:      make(map[int64]*RuleVersion),
		alarms:        make(map[string]*Alarm),
		owners:        make(map[string]*AlarmOwner),
		escalations:   make(map[string]*Escalation),
		notifications: make(map[string]*Notification),
		channels:      make(map[string]*ChannelConfig),
		windows:       make(map[string]*MaintenanceWindow),
		batches:       make(map[string]*EvalBatch),
		audit:         make(map[string]*AuditEntry),
	}
	s.restore()
	return s
}

func (s *State) Root() string {
	return s.root
}

func (s *State) restore() {
	var effective int64
	if err := LoadJSON(filepath.Join(s.root, "effective.json"), &effective); err == nil {
		s.effective = effective
	}
	var cursor int64
	if err := LoadJSON(filepath.Join(s.root, "cursor.json"), &cursor); err == nil {
		s.cursor = cursor
	}
	var route string
	if err := LoadJSON(filepath.Join(s.root, "current-route.json"), &route); err == nil {
		s.currentRoute = route
	}
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *State) PutDevice(d *Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[d.ID] = d
	return SaveJSON(filepath.Join(s.root, "devices", d.ID+".json"), d)
}

func (s *State) Device(id string) (*Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.devices[id]
	return d, ok
}

func (s *State) Devices() []*Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Device, 0, len(s.devices))
	for _, k := range sortedKeys(s.devices) {
		out = append(out, s.devices[k])
	}
	return out
}

func (s *State) PutPoint(p *Point) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.points[p.ID] = p
	return SaveJSON(filepath.Join(s.root, "points", p.ID+".json"), p)
}

func (s *State) Point(id string) (*Point, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.points[id]
	return p, ok
}

func (s *State) PointByName(name string) (*Point, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.points {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

func (s *State) Points() []*Point {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Point, 0, len(s.points))
	for _, k := range sortedKeys(s.points) {
		out = append(out, s.points[k])
	}
	return out
}

func (s *State) PutRule(r *Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[r.ID] = r
	return SaveJSON(filepath.Join(s.root, "rules", r.ID+".json"), r)
}

func (s *State) Rule(id string) (*Rule, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	return r, ok
}

func (s *State) Rules() []*Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Rule, 0, len(s.rules))
	for _, k := range sortedKeys(s.rules) {
		out = append(out, s.rules[k])
	}
	return out
}

func (s *State) PutVersion(v *RuleVersion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.versions[v.Version] = v
	return SaveJSON(filepath.Join(s.root, "versions", strconv.FormatInt(v.Version, 10)+".json"), v)
}

func (s *State) Version(v int64) (*RuleVersion, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ver, ok := s.versions[v]
	return ver, ok
}

func (s *State) Versions() []*RuleVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]int64, 0, len(s.versions))
	for k := range s.versions {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	out := make([]*RuleVersion, 0, len(keys))
	for _, k := range keys {
		out = append(out, s.versions[k])
	}
	return out
}

func (s *State) SetEffective(v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.effective = v
	return SaveJSON(filepath.Join(s.root, "effective.json"), v)
}

func (s *State) Effective() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.effective
}

func (s *State) PutAlarm(a *Alarm) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alarms[a.ID] = a
	return SaveJSON(filepath.Join(s.root, "alarms", a.ID+".json"), a)
}

func (s *State) Alarm(id string) (*Alarm, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.alarms[id]
	return a, ok
}

func (s *State) Alarms() []*Alarm {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Alarm, 0, len(s.alarms))
	for _, k := range sortedKeys(s.alarms) {
		out = append(out, s.alarms[k])
	}
	return out
}

func (s *State) PutOwner(o *AlarmOwner) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.owners[o.AlarmID] = o
	return SaveJSON(filepath.Join(s.root, "owners", o.AlarmID+".json"), o)
}

func (s *State) Owner(alarmID string) (*AlarmOwner, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.owners[alarmID]
	return o, ok
}

func (s *State) PutEscalation(e *Escalation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.escalations[e.AlarmID] = e
	return SaveJSON(filepath.Join(s.root, "escalations", e.AlarmID+".json"), e)
}

func (s *State) Escalation(alarmID string) (*Escalation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.escalations[alarmID]
	return e, ok
}

func (s *State) Escalations() []*Escalation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Escalation, 0, len(s.escalations))
	for _, k := range sortedKeys(s.escalations) {
		out = append(out, s.escalations[k])
	}
	return out
}

func (s *State) PutNotification(n *Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[n.ID] = n
	return SaveJSON(filepath.Join(s.root, "notifications", n.ID+".json"), n)
}

func (s *State) Notification(id string) (*Notification, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notifications[id]
	return n, ok
}

func (s *State) Notifications() []*Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Notification, 0, len(s.notifications))
	for _, k := range sortedKeys(s.notifications) {
		out = append(out, s.notifications[k])
	}
	return out
}

func (s *State) PutChannel(c *ChannelConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.channels[c.ID] = c
	return SaveJSON(filepath.Join(s.root, "channels", c.ID+".json"), c)
}

func (s *State) Channel(id string) (*ChannelConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.channels[id]
	return c, ok
}

func (s *State) Channels() []*ChannelConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*ChannelConfig, 0, len(s.channels))
	for _, k := range sortedKeys(s.channels) {
		out = append(out, s.channels[k])
	}
	return out
}

func (s *State) SetCurrentRoute(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentRoute = id
	return SaveJSON(filepath.Join(s.root, "current-route.json"), id)
}

func (s *State) CurrentRoute() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentRoute
}

func (s *State) PutWindow(w *MaintenanceWindow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[w.ID] = w
	return SaveJSON(filepath.Join(s.root, "windows", w.ID+".json"), w)
}

func (s *State) Windows() []*MaintenanceWindow {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*MaintenanceWindow, 0, len(s.windows))
	for _, k := range sortedKeys(s.windows) {
		out = append(out, s.windows[k])
	}
	return out
}

func (s *State) WindowForDevice(deviceID string) (*MaintenanceWindow, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, k := range sortedKeys(s.windows) {
		if s.windows[k].DeviceID == deviceID {
			return s.windows[k], true
		}
	}
	return nil, false
}

func (s *State) PutBatch(b *EvalBatch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches[b.ID] = b
	return SaveJSON(filepath.Join(s.root, "batches", b.ID+".json"), b)
}

func (s *State) Batch(id string) (*EvalBatch, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.batches[id]
	return b, ok
}

func (s *State) Batches() []*EvalBatch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*EvalBatch, 0, len(s.batches))
	for _, k := range sortedKeys(s.batches) {
		out = append(out, s.batches[k])
	}
	return out
}

func (s *State) SetCursor(v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cursor = v
	return SaveJSON(filepath.Join(s.root, "cursor.json"), v)
}

func (s *State) Cursor() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cursor
}

func (s *State) AppendAudit(e *AuditEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audit[e.ID] = e
	return SaveJSON(filepath.Join(s.root, "audit", e.ID+".json"), e)
}

func (s *State) AuditEntries() []*AuditEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*AuditEntry, 0, len(s.audit))
	for _, k := range sortedKeys(s.audit) {
		out = append(out, s.audit[k])
	}
	return out
}

func (s *State) SaveSnapshot(version int64, entries []RuleSnapshotEntry) error {
	return SaveJSON(filepath.Join(s.root, "snapshots", strconv.FormatInt(version, 10)+".json"), entries)
}

func (s *State) LoadSnapshot(version int64) ([]RuleSnapshotEntry, error) {
	var entries []RuleSnapshotEntry
	if err := LoadJSON(filepath.Join(s.root, "snapshots", strconv.FormatInt(version, 10)+".json"), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
