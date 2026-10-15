package bridge

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

// Event is one Server-Sent Event.
type Event struct {
	Type  string // "" means the default "message" type
	ID    string
	Data  string        // data lines joined with "\n"
	Retry time.Duration // the retry: field, when the event block set one
}

// IsMessage reports whether the event carries a JSON-RPC message (default or "message" type).
func (e Event) IsMessage() bool { return e.Type == "" || e.Type == "message" }

// ParseSSE reads an event stream and calls fn for each dispatched event. A block with only
// id: or retry: fields is still dispatched (with empty Data) so callers can track them.
// A trailing event without a terminating blank line is still dispatched at EOF.
func ParseSSE(r io.Reader, fn func(Event) error) error {
	br := bufio.NewReader(r)
	var ev Event
	var data []string
	hasData := false
	dispatch := func() error {
		defer func() { ev, data, hasData = Event{}, nil, false }()
		if !hasData && ev.ID == "" && ev.Retry == 0 {
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
				case "retry":
					if ms, err := strconv.Atoi(value); err == nil && ms >= 0 {
						ev.Retry = time.Duration(ms) * time.Millisecond
					}
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
