package processor

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

	return msg
}
