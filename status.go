package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TaskStatus int

const (
	Todo TaskStatus = iota
	Doing
	Done
)

func (s TaskStatus) String() string {
	switch s {
	case Todo:
		return "todo"
	case Doing:
		return "doing"
	case Done:
		return "done"
	}
	return "unknown"
}

func (s TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON must be a *pointer receiver* to ensure that the indirect from the
// parsed value can be set on the unmarshaling object. This means that the
// ParseSuit function must return a *value* and not a pointer.
func (s *TaskStatus) UnmarshalJSON(data []byte) (err error) {
	var suits string
	if err := json.Unmarshal(data, &suits); err != nil {
		return err
	}
	if *s, err = ParseTaskStatus(suits); err != nil {
		return err
	}
	return nil
}

func ParseTaskStatus(s string) (TaskStatus, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "todo" {
		return Todo, nil
	} else if s == "doing" {
		return Doing, nil
	} else if s == "done" {
		return Done, nil
	}
	return TaskStatus(0), fmt.Errorf("%q is not a valid task status", s)
}
