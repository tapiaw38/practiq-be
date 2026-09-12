package ai

import "testing"

func TestContainsTechnicalFallback(t *testing.T) {
	tests := []struct {
		name        string
		curiosities []string
		want        bool
	}{
		{name: "normal curiosities", curiosities: []string{"Los números nos ayudan a entender el mundo"}, want: false},
		{name: "technical fallback", curiosities: []string{"¡Hola! Soy Ana 😊 Disculpa, estoy teniendo algunos problemas técnicos en este momento."}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsTechnicalFallback(tt.curiosities); got != tt.want {
				t.Fatalf("containsTechnicalFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}
