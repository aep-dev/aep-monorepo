package cmd

import (
	"testing"

	"github.com/aep-dev/aep-e2e-validator/pkg/validator"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name    string
		raw     []string
		want    []validator.Header
		wantErr bool
	}{
		{
			name: "single header",
			raw:  []string{"Authorization=Bearer tok"},
			want: []validator.Header{{Key: "Authorization", Value: "Bearer tok"}},
		},
		{
			name: "multiple headers",
			raw:  []string{"A=1", "B=2"},
			want: []validator.Header{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}},
		},
		{
			name: "value with equals",
			raw:  []string{"Authorization=Bearer=tok=extra"},
			want: []validator.Header{{Key: "Authorization", Value: "Bearer=tok=extra"}},
		},
		{
			name: "whitespace trimmed",
			raw:  []string{"  Key  =  Value  "},
			want: []validator.Header{{Key: "Key", Value: "Value"}},
		},
		{
			name:    "missing colon",
			raw:     []string{"bad"},
			wantErr: true,
		},
		{
			name: "empty input",
			raw:  []string{},
			want: []validator.Header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeaders(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseHeaders() len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseHeaders()[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
