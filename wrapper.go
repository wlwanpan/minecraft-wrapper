package wrapper

import (
	"io"

	"github.com/looplab/fsm"
)

type Event string

const (
	EmptyEvent   Event = "empty"
	StartedEvent       = "started"
	StoppedEvent       = "stopped"
	StartEvent         = "start"
	StopEvent          = "stop"
)

const (
	ServerOffline  = "offline"
	ServerOnline   = "online"
	ServerStarting = "starting"
	ServerStopping = "stopping"
)

var wrapperFsmEvents = fsm.Events{
	fsm.EventDesc{
		Name: StopEvent,
		Src:  []string{ServerOnline},
		Dst:  ServerStopping,
	},
	fsm.EventDesc{
		Name: StoppedEvent,
		Src:  []string{ServerStopping},
		Dst:  ServerOffline,
	},
	fsm.EventDesc{
		Name: StartEvent,
		Src:  []string{ServerOffline},
		Dst:  ServerStarting,
	},
	fsm.EventDesc{
		Name: StartedEvent,
		Src:  []string{ServerStarting},
		Dst:  ServerOnline,
	},
}

type StateChangeFunc func(Event, Event, *Wrapper)

type Wrapper struct {
	console        Console
	parser         LogParser
	machine        *fsm.FSM
	stateChangeCBs []StateChangeFunc
}

func NewDefaultWrapper(server string, initial, max int) *Wrapper {
	cmd := JavaExecCmd(server, initial, max)
	console := NewConsole(cmd)
	return NewWrapper(console, LogParserFunc)
}

func NewWrapper(c Console, p LogParser) *Wrapper {
	wpr := &Wrapper{
		console: c,
		parser:  p,
	}
	wpr.newFSM()
	return wpr
}

func (w *Wrapper) newFSM() {
	w.machine = fsm.NewFSM(
		ServerOffline,
		wrapperFsmEvents,
		fsm.Callbacks{
			"enter_state": func(ev *fsm.Event) {
				w.triggerStateChangeCBs(Event(ev.Src), Event(ev.Dst))
			},
		},
	)
}

func (w *Wrapper) triggerStateChangeCBs(from, to Event) {
	for _, f := range w.stateChangeCBs {
		f(from, to, w)
	}
}

func (w *Wrapper) processLogEvents() {
	for {
		line, err := w.console.ReadLine()
		if err == io.EOF {
			w.updateState(StoppedEvent)
			return
		}

		event := w.parseLineToEvent(line)
		w.updateState(event)
	}
}

func (w *Wrapper) parseLineToEvent(line string) Event {
	return w.parser(line)
}

func (w *Wrapper) updateState(ev Event) error {
	if ev == EmptyEvent {
		return nil
	}
	return w.machine.Event(string(ev))
}

func (w *Wrapper) RegisterStateChangeCBS(cbs ...StateChangeFunc) {
	w.stateChangeCBs = append(w.stateChangeCBs, cbs...)
}

func (w *Wrapper) State() string {
	return w.machine.Current()
}

func (w *Wrapper) Start() error {
	go w.processLogEvents()
	return w.console.Start()
}

func (w *Wrapper) Stop() error {
	return w.console.WriteCmd("stop")
}

func (w *Wrapper) Kill() error {
	return w.console.Kill()
}
