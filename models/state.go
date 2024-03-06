package models

import "time"

type State struct {
	Dates          []string
	Emails         []string
	Votes          map[string]int
	PollDuration   time.Duration
	StartPollTimer chan bool
}

func NewState() *State {
	state := State{StartPollTimer: make(chan bool, 1)}
	state.StartPollTimer <- true

	return &state
}

func (state *State) ResetState() {
	state = &State{StartPollTimer: make(chan bool, 1)}
	state.StartPollTimer <- true
}
