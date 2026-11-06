package payments

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListBillingCyclesUsesServerKeyAndDecodesHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/subscriptions/subscriptions/user/teacher-1/billing-cycles" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "server-key" {
			t.Errorf("key = %q", r.Header.Get("X-API-Key"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"subscription_id":4,"status":"rejected","amount":2500.50,"currency":"ARS","period_start":"2026-10-01T00:00:00","period_end":"2026-11-01T00:00:00"}]`))
	}))
	defer server.Close()

	cycles, err := NewClient(server.URL, "server-key").ListBillingCycles(context.Background(), "teacher-1")
	if err != nil {
		t.Fatalf("ListBillingCycles: %v", err)
	}
	if len(cycles) != 1 || cycles[0].SubscriptionID != 4 || cycles[0].Amount != 2500.50 || cycles[0].PeriodStart.IsZero() {
		t.Fatalf("cycles = %+v", cycles)
	}
}
