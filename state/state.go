package state

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
	state.StartPollTimer <- false

	return &state
}

func (state *State) ResetState() {
	state.Dates = nil
	state.Emails = nil
	state.Votes = make(map[string]int)

	state.PollDuration = 0

	close(state.StartPollTimer)
	state.StartPollTimer = make(chan bool, 1)
	state.StartPollTimer <- false
}
