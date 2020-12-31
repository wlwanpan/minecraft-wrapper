package wrapper

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/looplab/fsm"
	"github.com/wlwanpan/minecraft-wrapper/events"
)

const (
	WrapperOffline  = "offline"
	WrapperOnline   = "online"
	WrapperStarting = "starting"
	WrapperStopping = "stopping"
	WrapperSaving   = "saving"
)

var (
	ErrWrapperResponseTimeout = errors.New("response timeout")

	ErrWrapperFailedToParseSeed = errors.New("failed to parse seed")
)

var wrapperFsmEvents = fsm.Events{
	fsm.EventDesc{
		Name: events.Stopping,
		Src:  []string{WrapperOnline},
		Dst:  WrapperStopping,
	},
	fsm.EventDesc{
		Name: events.Stopped,
		Src:  []string{WrapperStopping},
		Dst:  WrapperOffline,
	},
	fsm.EventDesc{
		Name: events.Starting,
		Src:  []string{WrapperOffline},
		Dst:  WrapperStarting,
	},
	fsm.EventDesc{
		Name: events.Started,
		Src:  []string{WrapperStarting},
		Dst:  WrapperOnline,
	},
	fsm.EventDesc{
		Name: events.Saving,
		Src:  []string{WrapperOnline},
		Dst:  WrapperSaving,
	},
	fsm.EventDesc{
		Name: events.Saved,
		Src:  []string{WrapperSaving},
		Dst:  WrapperOnline,
	},
}

type eventsQueue struct {
	q map[string]chan events.GameEvent
}

func newEventsQueue() *eventsQueue {
	return &eventsQueue{
		q: make(map[string]chan events.GameEvent),
	}
}

func (eq *eventsQueue) get(e string) <-chan events.GameEvent {
	_, ok := eq.q[e]
	if !ok {
		eq.q[e] = make(chan events.GameEvent)
	}
	return eq.q[e]
}

func (eq *eventsQueue) push(ev events.GameEvent) {
	e := ev.String()
	_, ok := eq.q[e]
	if !ok {
		eq.q[e] = make(chan events.GameEvent)
	}
	select {
	case eq.q[e] <- ev:
	default:
	}
}

type StateChangeFunc func(*Wrapper, events.Event, events.Event)

type Wrapper struct {
	Version        string
	machine        *fsm.FSM
	console        Console
	parser         LogParser
	clock          *clock
	eq             *eventsQueue
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
		clock:          newClock(),
		eq:             newEventsQueue(),
		gameEventsChan: make(chan events.GameEvent, 10),
	}
	wpr.newFSM()
	return wpr
}

func (w *Wrapper) newFSM() {
	w.machine = fsm.NewFSM(
		WrapperOffline,
		wrapperFsmEvents,
		fsm.Callbacks{
			"enter_state": func(ev *fsm.Event) {
				srcEvent := events.NewStateEvent(ev.Src)
				dstEvent := events.NewStateEvent(ev.Dst)
				go w.triggerStateChangeCBs(srcEvent, dstEvent)
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
			w.updateState(event.(events.StateEvent))
		case events.TypeCmd:
			w.handleCmdEvent(event.(events.GameEvent))
		case events.TypeGame:
			select {
			case w.gameEventsChan <- event.(events.GameEvent):
			default:
			}
		default:
		}
	}
}

func (w *Wrapper) parseLineToEvent(line string) (events.Event, events.EventType) {
	return w.parser(line, w.clock.Tick)
}

func (w *Wrapper) updateState(ev events.StateEvent) error {
	return w.machine.Event(ev.String())
}

func (w *Wrapper) handleCmdEvent(ev events.GameEvent) {
	if ev.Is(events.TimeIsEvent) {
		w.clock.syncTick(ev.Tick)
		return
	}
	if ev.Is(events.VersionEvent) {
		w.Version = ev.Data["version"]
		return
	}
	w.eq.push(ev)
}

func (w *Wrapper) processClock() {
	w.clock.start()
	for {
		<-w.clock.requestSync()
		w.clock.resetLastSync()
		w.console.WriteCmd("time query daytime")
	}
}

func (w *Wrapper) processCmdToEvent(cmd, e string, timeout time.Duration) (events.GameEvent, error) {
	evChan := w.eq.get(e)
	if err := w.console.WriteCmd(cmd); err != nil {
		return events.NilGameEvent, err
	}

	select {
	case <-time.After(timeout):
		return events.NilGameEvent, ErrWrapperResponseTimeout
	case ev := <-evChan:
		errMessage, ok := ev.Data["error_message"]
		if ok {
			// If the game event carries an 'error_message' in its Data field,
			// wrap and propagate the error message as an error.
			return events.NilGameEvent, errors.New(errMessage)
		}
		return ev, nil
	}
}

// GameEvents returns a receive-only channel of game related event. For example:
// - Player joined, left, died, was banned.
// - Game updates like game mode changes.
// - Player sends messages...
func (w *Wrapper) GameEvents() <-chan events.GameEvent {
	return w.gameEventsChan
}

// RegisterStateChangeCBs allow you to register a callback func
// that is called on each state changes to your minecraft server
// For example: server goes from 'offline' to 'starting'.
func (w *Wrapper) RegisterStateChangeCBs(cbs ...StateChangeFunc) {
	w.stateChangeCBs = append(w.stateChangeCBs, cbs...)
}

func (w *Wrapper) Ban(player, reason string) error {
	cmd := strings.Join([]string{"ban", player, reason}, " ")
	return w.console.WriteCmd(cmd)
}

// DataGet returns the Go struct representation of an 'entity' or 'block' or 'storage'.
// The data is originally stored in a NBT format.
func (w *Wrapper) DataGet(t, id string) (*DataGetOutput, error) {
	cmd := fmt.Sprintf("data get %s %s", t, id)
	ev, err := w.processCmdToEvent(cmd, events.DataGet, 3*time.Second)
	if err != nil {
		return nil, err
	}
	rawData := []byte(ev.Data["data_raw"])
	resp := &DataGetOutput{}
	err = DecodeSNBT(rawData, resp)
	return resp, err
}

// DefaultGameMode sets the default game mode for new players joining.
func (w *Wrapper) DefaultGameMode(mode GameMode) error {
	cmd := fmt.Sprintf("defaultgamemode %s", mode)
	return w.console.WriteCmd(cmd)
}

// SaveAll marks all chunks and player data to be saved to the data storage device.
// When flush is true, the marked data are saved immediately.
func (w *Wrapper) SaveAll(flush bool) error {
	cmd := "save-all"
	if flush {
		cmd += " flush"
	}
	return w.console.WriteCmd(cmd)
}

// Say sends the given message in the minecraft in-game chat.
func (w *Wrapper) Say(msg string) error {
	return w.console.WriteCmd("say " + msg)
}

// Seed returns the world seed.
func (w *Wrapper) Seed() (int, error) {
	ev, err := w.processCmdToEvent("seed", events.Seed, 1*time.Second)
	if err != nil {
		return 0, err
	}
	rawData := []byte(ev.Data["data_raw"])
	resp := []int{}
	err = DecodeSNBT(rawData, &resp)
	if len(resp) < 1 {
		return 0, ErrWrapperFailedToParseSeed
	}
	return resp[0], err
}

// Start will initialize the minecraft java process and start
// orchestrating the wrapper machine.
func (w *Wrapper) Start() error {
	go w.processLogEvents()
	go w.processClock()
	return w.console.Start()
}

// State returns the current state of the server, it can be one of:
// 'offline', 'online', 'starting' or 'stopping'.
func (w *Wrapper) State() string {
	return w.machine.Current()
}

// Stop pipes a 'stop' command to the minecraft java process.
func (w *Wrapper) Stop() error {
	return w.console.WriteCmd("stop")
}

// Kill the java process, use with caution since it will not trigger a save game.
func (w *Wrapper) Kill() error {
	return w.console.Kill()
}

// Tick returns the current minecraft game tick, which runs at a fixed rate
// of 20 ticks per second, src: https://minecraft.gamepedia.com/Tick.
func (w *Wrapper) Tick() int {
	return w.clock.Tick
}
