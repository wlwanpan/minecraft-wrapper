package wrapper

import (
	"fmt"
	"io"

	"github.com/looplab/fsm"
	"github.com/wlwanpan/minecraft-wrapper/events"
)

const (
	ServerOffline  = "offline"
	ServerOnline   = "online"
	ServerStarting = "starting"
	ServerStopping = "stopping"
	ServerSaving   = "saving"
)

var wrapperFsmEvents = fsm.Events{
	fsm.EventDesc{
		Name: events.Stopping,
		Src:  []string{ServerOnline},
		Dst:  ServerStopping,
	},
	fsm.EventDesc{
		Name: events.Stopped,
		Src:  []string{ServerStopping},
		Dst:  ServerOffline,
	},
	fsm.EventDesc{
		Name: events.Starting,
		Src:  []string{ServerOffline},
		Dst:  ServerStarting,
	},
	fsm.EventDesc{
		Name: events.Started,
		Src:  []string{ServerStarting},
		Dst:  ServerOnline,
	},
	fsm.EventDesc{
		Name: events.Saving,
		Src:  []string{ServerOnline},
		Dst:  ServerSaving,
	},
	fsm.EventDesc{
		Name: events.Saved,
		Src:  []string{ServerSaving},
		Dst:  ServerOnline,
	},
}

type StateChangeFunc func(*Wrapper, events.Event, events.Event)

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
		f(w, from, to)
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

func (w *Wrapper) parseLineToEvent(line string) (events.Event, events.EventType) {
	return w.parser(line, w.clock.Tick)
}

func (w *Wrapper) updateState(ev events.Event) error {
	return w.machine.Event(ev.String())
}

func (w *Wrapper) handleGameEvent(ev events.GameEvent) {
	if ev.Is(events.TimeIsEvent) {
		w.clock.syncTick(ev.Tick)
		return
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

// GameEvents returns a read channel with any game events like:
// - Player joined
// - Player left
// - Player sent a message and so on.
func (w *Wrapper) GameEvents() <-chan events.GameEvent {
	return w.gameEventsChan
}

// RegisterStateChangeCBs allow you to register a callback func
// that is called on each state changes to your minecraft server
// For example: server goes from 'offline' to 'starting'.
func (w *Wrapper) RegisterStateChangeCBs(cbs ...StateChangeFunc) {
	w.stateChangeCBs = append(w.stateChangeCBs, cbs...)
}

// State returns the current state of the server, it can be one of:
// 'offline', 'online', 'starting' or 'stopping'.
func (w *Wrapper) State() string {
	return w.machine.Current()
}

// Tick returns the current minecraft game tick, which runs at a fixed rate
// of 20 ticks per second, src: https://minecraft.gamepedia.com/Tick.
func (w *Wrapper) Tick() int {
	return w.clock.Tick
}

// Start will initialize the minecraft java process and start
// orchestrating the wrapper machine.
func (w *Wrapper) Start() error {
	go w.processLogEvents()
	go w.processClock()
	return w.console.Start()
}

// Stop pipes a 'stop' command to the minecraft java process.
func (w *Wrapper) Stop() error {
	return w.console.WriteCmd("stop")
}

// Kill the java process, use with caution since it will not trigger a save game.
func (w *Wrapper) Kill() error {
	return w.console.Kill()
}

// Save triggers a save game.
func (w *Wrapper) Save() error {
	return w.console.WriteCmd("save-all")
}

func (w *Wrapper) DataGet(t, id string) error {
	cmd := fmt.Sprintf("data get %s %s", t, id)
	return w.console.WriteCmd(cmd)
}
