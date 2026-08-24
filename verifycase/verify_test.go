package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/channel"
	"telemetryguard/internal/notify"
	"telemetryguard/internal/store"
)

func TestNotifyUsesFailoverRoute(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	sms, err := channel.RegisterChannel(state, "sms", 1)
	if err != nil {
		t.Fatal(err)
	}
	wecom, err := channel.RegisterChannel(state, "wecom", 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := state.SetCurrentRoute(sms.ID); err != nil {
		t.Fatal(err)
	}
	if err := channel.Failover(state, sms.ID); err != nil {
		t.Fatal(err)
	}
	n, err := notify.Dispatch(state, "a1", 1, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Fatal("no notification dispatched")
	}
	if n.ChannelID != wecom.ID {
		t.Fatalf("notification went to %s, want %s", n.ChannelID, wecom.ID)
	}
}
