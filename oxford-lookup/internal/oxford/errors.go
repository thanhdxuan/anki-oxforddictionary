package oxford

import (
	"errors"
)

var (
	ErrNotFound = errors.New("word not found")
	ErrNoAudio  = errors.New("no UK audio on page")
)
