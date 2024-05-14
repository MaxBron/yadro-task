package main

import (
	"time"
)

type Table struct {
	start      time.Time
	end        time.Time
	hours      int
	revenue    int
	clientName string
	busy       bool
	duration   time.Time
}
