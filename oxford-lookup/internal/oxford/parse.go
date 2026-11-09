package oxford

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Result struct {
	Word      string `json:"word"`
	AudioUK   string `json:"audio_uk"`
	SourceURL string `json:"source_url"`
}

func Parse(r io.Reader) (Result, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return Result{}, err
	}
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Div {
			val, ok := attr(n, "id")

			if !ok || val != "uk_pron" {
				continue
			}

			src, ok := attr(n, "data-src-mp3")
			if !ok || src == "" {
				break
			}
			return Result{AudioUK: src}, nil
		}
	}
	return Result{}, ErrNoAudio
}

func attr(n *html.Node, key string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}
