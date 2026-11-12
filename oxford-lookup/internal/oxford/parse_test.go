package oxford

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {

	var tests = []struct {
		name    string
		file    string
		want    string
		wantErr error
	}{
		{"has audio", "chromatic.html", "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/c/chr/chrom/chromatic__gb_1.mp3", nil},
		{"no audio attr", "chromatic_noaudio.html", "", ErrNoAudio},
		{"empty page", "empty.html", "", ErrNoAudio},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.file)
			r, err := os.Open(filePath)

			if err != nil {
				t.Fatalf("Error when opening file: %s", filePath)
			}
			defer r.Close()

			src, err := Parse(r)

			if src.AudioUK != tt.want {
				t.Errorf("got %q, want %q", src.AudioUK, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %q, want %q", err, tt.wantErr)
			}

		})
	}
}
