package bridge

import (
	"bufio"
	"io"
	"strings"
)

// Event is one Server-Sent Event.
type Event struct {
	Type string // "" means the default "message" type
	ID   string
	Data string // data lines joined with "\n"
}

// IsMessage reports whether the event carries a JSON-RPC message (default or "message" type).
func (e Event) IsMessage() bool { return e.Type == "" || e.Type == "message" }

// ParseSSE reads an event stream and calls fn for each dispatched event.
// A trailing event without a terminating blank line is still dispatched at EOF.
func ParseSSE(r io.Reader, fn func(Event) error) error {
	br := bufio.NewReader(r)
	var ev Event
	var data []string
	hasData := false
	dispatch := func() error {
		defer func() { ev, data, hasData = Event{}, nil, false }()
		if !hasData {
			return nil
		}
		ev.Data = strings.Join(data, "\n")
		return fn(ev)
	}
	for {
		line, err := br.ReadString('\n')
		if len(line) > 0 || err == nil {
			line = strings.TrimRight(line, "\r\n")
			if line == "" {
				if derr := dispatch(); derr != nil {
					return derr
				}
			} else if !strings.HasPrefix(line, ":") {
				field, value, _ := strings.Cut(line, ":")
				value = strings.TrimPrefix(value, " ")
				switch field {
				case "event":
					ev.Type = value
				case "data":
					data = append(data, value)
					hasData = true
				case "id":
					ev.ID = value
				}
			}
		}
		if err == io.EOF {
			return dispatch()
		}
		if err != nil {
			return err
		}
	}
}
