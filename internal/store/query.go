package store

import "sort"

// ActiveAlarms returns alarms whose live status is open or acknowledged.
func (s *State) ActiveAlarms() []*Alarm {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Alarm, 0)
	for _, k := range sortedKeys(s.alarms) {
		a := s.alarms[k]
		if a.Status == AlarmOpen || a.Status == AlarmAcknowledged {
			out = append(out, a)
		}
	}
	return out
}

// AlarmByPoint returns the latest alarm record for a point.
func (s *State) AlarmByPoint(pointID string) (*Alarm, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var latest *Alarm
	for _, a := range s.alarms {
		if a.PointID != pointID {
			continue
		}
		if latest == nil || a.LastTriggerAt.After(latest.LastTriggerAt) {
			latest = a
		}
	}
	return latest, latest != nil
}

// PublishedRules returns rules that belong to a published version.
func (s *State) PublishedRules() []*Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Rule, 0)
	for _, k := range sortedKeys(s.rules) {
		r := s.rules[k]
		if r.State == RulePublished {
			out = append(out, r)
		}
	}
	return out
}

// DraftRules returns rules still in draft state.
func (s *State) DraftRules() []*Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Rule, 0)
	for _, k := range sortedKeys(s.rules) {
		r := s.rules[k]
		if r.State == RuleDraft {
			out = append(out, r)
		}
	}
	return out
}

// SortedNotifications returns notifications ordered by creation time.
func (s *State) SortedNotifications() []*Notification {
	all := s.Notifications()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.Before(all[j].CreatedAt) })
	return all
}

// ChannelByPriority returns channels ordered by priority ascending.
func (s *State) ChannelByPriority() []*ChannelConfig {
	all := s.Channels()
	sort.Slice(all, func(i, j int) bool {
		if all[i].Priority == all[j].Priority {
			return all[i].Name < all[j].Name
		}
		return all[i].Priority < all[j].Priority
	})
	return all
}

// FailedNotifications returns notifications that are still pending retry.
func (s *State) FailedNotifications() []*Notification {
	out := make([]*Notification, 0)
	for _, n := range s.Notifications() {
		if n.Status == NotificationFailed || n.Status == NotificationPending {
			out = append(out, n)
		}
	}
	return out
}
