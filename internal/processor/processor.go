package processor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cabify/timex"
	"go.uber.org/zap"

	"github.com/kulti/task-list-bot/internal/models"
)

type store interface {
	CurrentSprint() (models.TaskList, error)
	CreateNewSprint(begin, end time.Time) error
	CreateTask(text string, points int) error
	DoneTask(id int) (string, error)
	BurnTaskPoints(id int, burnt int) (string, error)
}

// Processor is a telemgram message processor.
type Processor struct {
	store store
}

// New creates a new instance of Processor.
func New(store store) *Processor {
	return &Processor{
		store: store,
	}
}

// Process processes input message.
func (p *Processor) Process(msg string) string {
	if len(msg) == 0 {
		panic("empty message are not allowed")
	}

	if msg[0] == '/' {
		return p.processCommand(msg)
	}
	return p.processNewTask(msg)
}

func (p *Processor) processNewTask(msg string) string {
	msgSplitted := strings.Split(msg, " ")

	task := msg
	points := 1
	if len(msgSplitted) > 1 {
		parsedPoints, err := strconv.Atoi(msgSplitted[0])
		if err == nil {
			task = strings.Join(msgSplitted[1:], " ")
			points = parsedPoints
		} else {
			parsedPoints, err := strconv.Atoi(msgSplitted[len(msgSplitted)-1])
			if err == nil {
				task = strings.Join(msgSplitted[:len(msgSplitted)-1], " ")
				points = parsedPoints
			}
		}
	}

	logger := zap.L().With(zap.String("msg", msg))
	err := p.store.CreateTask(task, points)
	if err != nil {
		logger.Warn("failed to create new task", zap.Error(err))
		return "Oops! Failed to create new task. Try later."
	}

	return p.fullTaskList(logger)
}

func (p *Processor) processCommand(msg string) string {
	const msgParts = 2
	msgSplitted := strings.SplitN(msg, " ", msgParts)
	if len(msgSplitted) != msgParts {
		return ""
	}

	cmd := msgSplitted[0]
	params := msgSplitted[1]
	logger := zap.L().With(zap.String("cmd", cmd))

	switch cmd {
	case "/ns":
		return p.processNewSprint(logger, params)
	case "/d":
		return p.processDoneTask(logger, params)
	}

	return ""
}

func (p *Processor) processNewSprint(logger *zap.Logger, params string) string {
	begin, end, ok := p.parseBeginEnd(params)
	if !ok {
		return "Invalid format of new sprint. Should be `DD.MM - DD.MM` (e.g. `01.12 - 07.12`)"
	}

	err := p.store.CreateNewSprint(begin, end)
	if err != nil {
		logger.Warn("failed to create new sprint", zap.Error(err))
		return "Oops! Failed to create new sprint. Try later."
	}
	return p.fullTaskList(logger)
}

func (p *Processor) processDoneTask(logger *zap.Logger, params string) string {
	doneParams := strings.Split(params, " ")

	const doneParamsCount = 2
	if len(doneParams) > doneParamsCount {
		return "Too much params for done command. Should be `/d id [burnt]`."
	}

	id, err := strconv.Atoi(doneParams[0])
	if err != nil {
		return "Task id should be a number."
	}

	var task string
	if len(doneParams) == 1 {
		task, err = p.store.DoneTask(id)
		if err != nil {
			logger.Warn("failed to done task", zap.Error(err))
			return "Oops! Failed to done task. Try later."
		}
	} else {
		burnt, err := strconv.Atoi(doneParams[1])
		if err != nil {
			return "Burnt points should be a number."
		}

		task, err = p.store.BurnTaskPoints(id, burnt)
		if err != nil {
			logger.Warn("failed to burn task points", zap.Error(err))
			return "Oops! Failed to burn task points. Try later."
		}
	}

	taskList := p.fullTaskList(logger)
	return fmt.Sprintf("The task %q is marked as done.\n\n%s", task, taskList)
}

func (p *Processor) fullTaskList(logger *zap.Logger) string {
	taskList, err := p.store.CurrentSprint()
	if err != nil {
		logger.Warn("failed to read current sprint", zap.Error(err))
		return "Oops! Failed to load current task list. Try later."
	}

	if taskList.Title == "" {
		return "no sprint"
	}

	b := strings.Builder{}
	b.WriteString(taskList.Title)
	b.WriteByte('\n')

	b.WriteString("Total ")
	b.WriteString(taskList.Points.String())
	b.WriteByte('\n')
	b.WriteByte('\n')

	for _, t := range taskList.Tasks {
		b.WriteString(strconv.Itoa(t.ID))
		b.WriteByte(' ')
		b.WriteString(t.Points.String())
		b.WriteByte(' ')
		b.WriteString(t.Text)
		b.WriteByte('\n')
	}
	return b.String()
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
