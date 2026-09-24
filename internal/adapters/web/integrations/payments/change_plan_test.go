package payments

import (
	"encoding/json"
	"testing"
)

// Moving down charges nothing, and the screen says so rather than claiming a
// payment that never happened.
func TestPlanChangeReportsWhatWasCharged(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want float64
	}{
		{"moving up charges the difference", `{"subscription_id":10,"plan_id":2,"status":"authorized","charged":50.0}`, 50},
		{"moving down charges nothing", `{"subscription_id":10,"plan_id":1,"status":"authorized","charged":0}`, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var change PlanChange
			if err := json.Unmarshal([]byte(tc.body), &change); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if change.Charged != tc.want {
				t.Fatalf("charged = %v, want %v", change.Charged, tc.want)
			}
		})
	}
}

// A paused teacher keeps the month they bought, so the entitlement is active
// and the status is what tells the screen to offer resuming.
func TestAPausedEntitlementIsActiveAndSaysSo(t *testing.T) {
	body := `{"user_id":"t","active":true,"status":"paused","access_until":"2026-10-24T02:49:15"}`
	var entitlement Entitlement
	if err := json.Unmarshal([]byte(body), &entitlement); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !entitlement.Active {
		t.Fatal("paid time must still count as entitled")
	}
	if entitlement.Status != "paused" {
		t.Fatalf("status = %q, want paused", entitlement.Status)
	}
}
