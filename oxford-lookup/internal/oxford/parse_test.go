package oxford

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

var tests = []struct {
	name    string
	file    string
	want    string
	wantErr error
}{
	{"has audio", "observe.html", "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/o/obs/obser/observe__gb_2.mp3", nil},
	{"no audio attr", "observe_noaudio.html", "", ErrNoAudio},
	{"empty page", "empty.html", "", ErrNoAudio},
}

func TestParser(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filepath := fmt.Sprintf("testdata/%s", tt.file)
			r, err := os.Open(filepath)

			if err != nil {
				t.Fatalf("Error when opening file: %s", filepath)
			}

			src, err := Parse(r)

			if src.AudioUK != tt.want {
				t.Errorf("got %q, want %q", src.AudioUK, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %q, want %q", err, tt.wantErr)
			}

			defer r.Close()
		})
	}
}
