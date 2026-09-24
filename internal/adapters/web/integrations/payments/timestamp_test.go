package payments

import (
	"encoding/json"
	"testing"
	"time"
)

// The payments service serialises naive datetimes. time.Time rejects them for
// want of an offset, and the failure is not a wrong date: the whole decode
// fails, the entitlement reads as a payments outage, and a teacher who paid is
// shown the free plan.
func TestEntitlementSurvivesATimestampWithoutAZone(t *testing.T) {
	// Copied from what the service actually answered.
	body := `{"user_id":"teacher-1","active":true,"subscription_id":10,"plan_id":1,` +
		`"access_until":"2026-10-24T02:49:15","metadata":{"max_students":5,"name":"Inicial"}}`

	var entitlement Entitlement
	if err := json.Unmarshal([]byte(body), &entitlement); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !entitlement.Active {
		t.Fatal("a paid entitlement must read as active")
	}
	want := time.Date(2026, 10, 24, 2, 49, 15, 0, time.UTC)
	if entitlement.AccessUntil == nil || !entitlement.AccessUntil.Equal(want) {
		t.Fatalf("access_until = %v, want %v", entitlement.AccessUntil, want)
	}
}

func TestATimestampThatNamesItsZoneIsBelieved(t *testing.T) {
	var ts Timestamp
	if err := json.Unmarshal([]byte(`"2026-10-24T02:49:15-03:00"`), &ts); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !ts.Equal(time.Date(2026, 10, 24, 5, 49, 15, 0, time.UTC)) {
		t.Fatalf("offset ignored: got %v", ts.UTC())
	}
}

func TestANullTimestampIsNotAnError(t *testing.T) {
	var entitlement Entitlement
	if err := json.Unmarshal([]byte(`{"active":false,"access_until":null}`), &entitlement); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
