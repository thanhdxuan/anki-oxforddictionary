package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"oxford-lookup/internal/oxford"
)

const baseURL = "https://www.oxfordlearnersdictionaries.com"

func main() {
	os.Exit(run())
}

func run() int {

	timeout := flag.Int("timeout", 10, "request timeout (seconds)")
	word := flag.String("word", "", "request word (string)")

	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*(time.Duration(*timeout)))
	defer cancel()

	client := &oxford.Client{BaseURL: baseURL}

	result, err := client.Lookup(ctx, *word)
	exit := 0
	var out any = result
	if err != nil {
		fmt.Fprintf(os.Stderr, "main: %v\n", err)

		switch {
		case errors.Is(err, oxford.ErrNotFound):
			out = &errorOutput{Word: *word, Error: "not_found", Message: err.Error()}
			exit = 1
		case errors.Is(err, oxford.ErrNoAudio):
			out = &errorOutput{Word: *word, Error: "no_audio", Message: err.Error()}
			exit = 1
		default:
			out = &errorOutput{Word: *word, Error: "network", Message: err.Error()}
			exit = 2
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		return 2
	}
	return exit
}
