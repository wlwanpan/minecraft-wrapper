package wrapper

import (
	"io"

	"github.com/looplab/fsm"
	"github.com/wlwanpan/minecraft-wrapper/events"
)

const (
	ServerOffline  = "offline"
	ServerOnline   = "online"
	ServerStarting = "starting"
	ServerStopping = "stopping"
)

var wrapperFsmEvents = fsm.Events{
	fsm.EventDesc{
		Name: events.Stop,
		Src:  []string{ServerOnline},
		Dst:  ServerStopping,
	},
	fsm.EventDesc{
		Name: events.Stopped,
		Src:  []string{ServerStopping},
		Dst:  ServerOffline,
	},
	fsm.EventDesc{
		Name: events.Start,
		Src:  []string{ServerOffline},
		Dst:  ServerStarting,
	},
	fsm.EventDesc{
		Name: events.Started,
		Src:  []string{ServerStarting},
		Dst:  ServerOnline,
	},
}

type StateChangeFunc func(events.Event, events.Event)

type Wrapper struct {
	machine        *fsm.FSM
	console        Console
	parser         LogParser
	clock          *Clock
	gameEventsChan chan (events.GameEvent)
	stateChangeCBs []StateChangeFunc
}

func NewDefaultWrapper(server string, initial, max int) *Wrapper {
	cmd := JavaExecCmd(server, initial, max)
	console := NewConsole(cmd)
	return NewWrapper(console, LogParserFunc)
}

func NewWrapper(c Console, p LogParser) *Wrapper {
	wpr := &Wrapper{
		console:        c,
		parser:         p,
		clock:          NewClock(),
		gameEventsChan: make(chan events.GameEvent, 10),
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
				srcEvent := events.NewStateEvent(ev.Src)
				dstEvent := events.NewStateEvent(ev.Dst)
				w.triggerStateChangeCBs(srcEvent, dstEvent)
			},
		},
	)
}

func (w *Wrapper) triggerStateChangeCBs(from, to events.Event) {
	for _, f := range w.stateChangeCBs {
		f(from, to)
	}
}

func (w *Wrapper) processLogEvents() {
	for {
		line, err := w.console.ReadLine()
		if err == io.EOF {
			w.updateState(events.StoppedEvent)
			return
		}

		event, t := w.parseLineToEvent(line)
		switch t {
		case events.TypeState:
			w.updateState(event)
		case events.TypeGame:
			w.handleGameEvent(event.(events.GameEvent))
		default:
		}
	}
}

func (w *Wrapper) parseLineToEvent(line string) (events.Event, int) {
	return w.parser(line, w.clock.Tick)
}

func (w *Wrapper) updateState(ev events.Event) error {
	return w.machine.Event(ev.String())
}

func (w *Wrapper) handleGameEvent(ev events.GameEvent) {
	if ev.Is(events.TimeIsEvent) {
		w.clock.syncTick(ev.Tick)
	}
	select {
	case w.gameEventsChan <- ev:
	default:
	}
}

func (w *Wrapper) processClock() {
	for {
		select {
		case <-w.clock.requestSync():
			w.clock.resetLastSync()
			w.console.WriteCmd("time query daytime")
		}
	}
}

func (w *Wrapper) GameEvents() <-chan events.GameEvent {
	return w.gameEventsChan
}

func (w *Wrapper) RegisterStateChangeCBs(cbs ...StateChangeFunc) {
	w.stateChangeCBs = append(w.stateChangeCBs, cbs...)
}

func (w *Wrapper) State() string {
	return w.machine.Current()
}

func (w *Wrapper) Start() error {
	go w.processLogEvents()
	go w.processClock()
	return w.console.Start()
}

func (w *Wrapper) Stop() error {
	return w.console.WriteCmd("stop")
}

func (w *Wrapper) Kill() error {
	return w.console.Kill()
}
