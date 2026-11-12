package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"oxford-lookup/internal/oxford"
	"time"
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "main: %v\n", err)

		switch {
		case errors.Is(err, oxford.ErrNotFound):
			o := &errorOutput{Word: *word, Error: "not_found", Message: err.Error()}
			json.NewEncoder(os.Stdout).Encode(o)
			exit = 1
		case errors.Is(err, oxford.ErrNoAudio):
			o := &errorOutput{Word: *word, Error: "no_audio", Message: err.Error()}
			json.NewEncoder(os.Stdout).Encode(o)
			exit = 1
		default:
			o := &errorOutput{Word: *word, Error: "network", Message: err.Error()}
			json.NewEncoder(os.Stdout).Encode(o)
			exit = 2
		}
	} else {
		json.NewEncoder(os.Stdout).Encode(result)
	}
	return exit
}
