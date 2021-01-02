package wrapper

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/looplab/fsm"
	"github.com/wlwanpan/minecraft-wrapper/events"
	"github.com/wlwanpan/minecraft-wrapper/snbt"
)

const (
	WrapperOffline  = "offline"
	WrapperOnline   = "online"
	WrapperStarting = "starting"
	WrapperStopping = "stopping"
	WrapperSaving   = "saving"
)

var (
	// ErrWrapperResponseTimeout is returned when a command fails to receive
	// its respective event from the server logs within some timeframe. Hence
	// no output could be decoded for the command.
	ErrWrapperResponseTimeout = errors.New("response timeout")
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

type StateChangeFunc func(*Wrapper, events.Event, events.Event)

// Wrapper is the minecraft-wrapper core struct, representing an instance
// of a minecraft server (JE). It is used to manage and interact with the
// java process by proxying its stdin and stdout via the Console interface.
type Wrapper struct {
	// Version is the minecraft server version being wrapped.
	// The Version is detected and set from the log line:
	// "Starting minecraft server version [X.X.X]""
	Version        string
	machine        *fsm.FSM
	console        Console
	parser         LogParser
	clock          *clock
	eq             *eventsQueue
	gameEventsChan chan (events.GameEvent)
	startedChan    chan bool
}

func NewDefaultWrapper(server string, initial, max int) *Wrapper {
	cmd := javaExecCmd(server, initial, max)
	console := newConsole(cmd)
	return NewWrapper(console, logParserFunc)
}

func NewWrapper(c Console, p LogParser) *Wrapper {
	wpr := &Wrapper{
		console:        c,
		parser:         p,
		clock:          newClock(),
		eq:             newEventsQueue(),
		gameEventsChan: make(chan events.GameEvent, 10),
		startedChan:    make(chan bool, 1),
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
				if ev.Src == WrapperStarting && ev.Dst == WrapperOnline {
					w.startedChan <- true
				}
			},
		},
	)
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

func (w *Wrapper) processCmdToEventArr(cmd, e string, timeout time.Duration) ([]events.GameEvent, error) {
	evChan := w.eq.get(e)
	if err := w.console.WriteCmd(cmd); err != nil {
		return nil, err
	}

	expectedEventsCount := 1
	events := []events.GameEvent{}
	for {
		select {
		case ev := <-evChan:
			entryType := ev.Data["entry_type"]
			if entryType == "header" {
				c, ok := ev.Data["entry_count"]
				if !ok {
					return events, nil
				}
				expectedEventsCount, _ = strconv.Atoi(c)
				break
			}

			events = append(events, ev)
			if len(events) >= expectedEventsCount {
				return events, nil
			}
		case <-time.After(timeout):
			return events, ErrWrapperResponseTimeout
		}
	}
}

// GameEvents returns a receive-only channel of game related event. For example:
// - Player joined, left, died, was banned.
// - Game updates like game mode changes.
// - Player sends messages...
func (w *Wrapper) GameEvents() <-chan events.GameEvent {
	return w.gameEventsChan
}

func (w *Wrapper) Ban(player, reason string) error {
	cmd := strings.Join([]string{"ban", player, reason}, " ")
	return w.console.WriteCmd(cmd)
}

func (w *Wrapper) BanList(t BanListType) ([]string, error) {
	cmd := fmt.Sprintf("banlist %s", t)
	evs, err := w.processCmdToEventArr(cmd, events.BanList, 10*time.Second)
	if err != nil {
		return nil, err
	}

	banList := []string{}
	for _, ev := range evs {
		banList = append(banList, ev.Data["entry_name"])
	}
	return banList, nil
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
	if err = snbt.Decode(rawData, resp); err != nil {
		return nil, err
	}
	return resp, err
}

// DefaultGameMode sets the default game mode for new players joining.
func (w *Wrapper) DefaultGameMode(mode GameMode) error {
	cmd := fmt.Sprintf("defaultgamemode %s", mode)
	return w.console.WriteCmd(cmd)
}

// DeOp removes a given player from the operator list.
func (w *Wrapper) DeOp(player string) error {
	return w.console.WriteCmd("deop " + player)
}

// Difficulty changes the game difficulty level of the world.
func (w *Wrapper) Difficulty(d GameDifficulty) error {
	cmd := fmt.Sprintf("difficulty %s", d)
	_, err := w.processCmdToEvent(cmd, events.Difficulty, 1*time.Second)
	return err
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

// SaveOn enables automatic saving. The server is allowed to write to the world files.
func (w *Wrapper) SaveOn() error {
	return w.console.WriteCmd("save-on")
}

// SaveOff disables automatic saving by preventing the server from writing to the world files.
func (w *Wrapper) SaveOff() error {
	return w.console.WriteCmd("save-off")
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
	err = snbt.Decode(rawData, &resp)
	return resp[0], err
}

// Start will initialize the minecraft java process and start
// orchestrating the wrapper machine.
func (w *Wrapper) Start() error {
	go w.processLogEvents()
	go w.processClock()
	return w.console.Start()
}

func (w *Wrapper) StartAndWait() (<-chan bool, error) {
	err := w.Start()
	return w.startedChan, err
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
