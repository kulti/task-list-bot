package processor

import "strings"

// Processor is a telemgram message processor.
type Processor struct{}

// New creates a new instance of Processor.
func New() *Processor {
	return &Processor{}
}

// Process processes input message.
func (s *Processor) Process(msg string) string {
	if len(msg) == 0 {
		panic("empty message are not allowed")
	}

	if msg[0] == '/' {
		return s.processCommand(msg)
	}
	return msg
}

func (s *Processor) processCommand(msg string) string {
	const msgParts = 2
	msgSplitted := strings.SplitN(msg, " ", msgParts)
	if len(msgSplitted) != msgParts {
		return ""
	}

	cmd := msgSplitted[0]
	params := msgSplitted[1]
	if cmd == "/ns" {
		return s.processNewSprint(params)
	}

	return ""
}

func (s *Processor) processNewSprint(params string) string {
	return params
}
