package oxford

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	mux := http.NewServeMux()
	ua := "oxford-lookup"
	mux.HandleFunc("GET /definition/english/chromatic", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != ua {
			http.Error(w, "bad UA", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, "./testdata/chromatic.html")
	})

	mux.HandleFunc("GET /definition/english/asgh", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/spellcheck/english/?q=asgh", http.StatusMovedPermanently)
	})

	mux.HandleFunc("GET /spellcheck/english/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./testdata/empty.html")
	})

	mux.HandleFunc("GET /definition/english/broken", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	mockServer := httptest.NewTLSServer(mux)

	defer mockServer.Close()

	var tests = []struct {
		word       string
		wantWord   string
		wantAuURL  string
		wantSrcURL string
		wantErr    error
		wantAnyErr bool
	}{
		{"chromatic", "chromatic", "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/c/chr/chrom/chromatic__gb_1.mp3", mockServer.URL + "/definition/english/chromatic", nil, false},
		{"asgh", "", "", mockServer.URL + "/spellcheck/english/?q=asgh", ErrNotFound, false},
		{"broken", "", "", mockServer.URL + "/definition/english/broken", nil, true},
	}

	client := &Client{HTTP: mockServer.Client(), BaseURL: mockServer.URL, UserAgent: ua}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			src, err := client.Lookup(ctx, tt.word)

			if src.AudioUK != tt.wantAuURL {
				t.Errorf("AudioUK: got %q, want %q", src.AudioUK, tt.wantAuURL)
			}

			if src.Word != tt.wantWord {
				t.Errorf("Word: got %q, want %q", src.Word, tt.wantWord)
			}

			if src.SourceURL != tt.wantSrcURL {
				t.Errorf("SourceURL: got %q, want %q", src.SourceURL, tt.wantSrcURL)
			}

			if tt.wantAnyErr {
				if err == nil || errors.Is(err, ErrNotFound) {
					t.Errorf("got %v, want an error != ErrorNotFound", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}

		})
	}

}
