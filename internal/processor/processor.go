package processor

import (
	"strconv"
	"strings"
	"time"

	"github.com/cabify/timex"
)

// Processor is a telemgram message processor.
type Processor struct{}

// New creates a new instance of Processor.
func New() *Processor {
	return &Processor{}
}

// Process processes input message.
func (p *Processor) Process(msg string) string {
	if len(msg) == 0 {
		panic("empty message are not allowed")
	}

	if msg[0] == '/' {
		return p.processCommand(msg)
	}
	return msg
}

func (p *Processor) processCommand(msg string) string {
	const msgParts = 2
	msgSplitted := strings.SplitN(msg, " ", msgParts)
	if len(msgSplitted) != msgParts {
		return ""
	}

	cmd := msgSplitted[0]
	params := msgSplitted[1]
	if cmd == "/ns" {
		return p.processNewSprint(params)
	}

	return ""
}

func (p *Processor) processNewSprint(params string) string {
	_, _, ok := p.parseBeginEnd(params)
	if !ok {
		return "Invalid format of new sprint. Should be `DD.MM - DD.MM` (e.g. `01.12 - 07.12`)"
	}

	return params
}

func (p *Processor) parseBeginEnd(s string) (time.Time, time.Time, bool) {
	const partsCount = 2
	parts := strings.Split(s, "-")
	if len(parts) != partsCount {
		return time.Time{}, time.Time{}, false
	}

	beginStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	begin, ok := p.parseNewSprintDate(beginStr)
	if !ok {
		return time.Time{}, time.Time{}, false
	}

	end, ok := p.parseNewSprintDate(endStr)
	if !ok {
		return time.Time{}, time.Time{}, false
	}

	return begin, end, true
}

func (p *Processor) parseNewSprintDate(s string) (time.Time, bool) {
	const dateParts = 2
	parts := strings.Split(s, ".")
	if len(parts) != dateParts {
		return time.Time{}, false
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, false
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, false
	}

	year := timex.Now().Year()

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), true
}
