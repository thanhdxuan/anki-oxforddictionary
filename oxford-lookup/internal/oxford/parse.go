package oxford

import (
	"io"
	"slices"
	"strings"

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
			val, ok := attr(n, "class")
			if !ok || val != "webtop" {
				continue
			}
			for wn := range n.Descendants() {
				val, ok := attr(wn, "class")
				if !ok || !slices.Contains(strings.Fields(val), "pron-uk") {
					continue
				}
				src, ok := attr(wn, "data-src-mp3")
				if !ok || src == "" {
					break
				}
				return Result{AudioUK: src}, nil
			}
			break
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
