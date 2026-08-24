package main

import "testing"

func TestConcatListEntry(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "simple", in: "clip.mp4", want: "file 'clip.mp4'"},
		{name: "space", in: "my clip.mp4", want: "file 'my clip.mp4'"},
		{name: "quote", in: "it's.mp4", want: "file 'it'\\''s.mp4'"},
		{name: "newline", in: "a\nb.mp4", wantErr: true},
		{name: "cr", in: "a\rb.mp4", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := concatListEntry(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("concatListEntry(%q) = %q, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("concatListEntry(%q) = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("concatListEntry(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
