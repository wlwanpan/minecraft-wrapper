package wrapper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
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
	// ErrWrapperNotOnline is returned when a commad is called but the wrapper
	// is not 'online'. The minecraft server is not loaded and ready to process
	// any commands.
	ErrWrapperNotOnline = errors.New("not online")
	// ErrPlayerNotFound is returned when a targetted command failed to process
	// due to the player not being connected to the server.
	ErrPlayerNotFound = errors.New("player not found")
	// ErrUnknownItem is returned when an item operation is called with an
	// invalid item type or structure.
	ErrUnknownItem = errors.New("unknown item")
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
	playerList     map[string]string
	ctxCancelFunc  context.CancelFunc
	gameEventsChan chan (events.GameEvent)
	loadedChan     chan bool
}

// NewDefaultWrapper returns a new instance of the Wrapper. This is
// the main method to use for your wrapper but if you wish to read
// and parse your own log lines to events, see 'NewWrapper'. This
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
		playerList:     map[string]string{},
		ctxCancelFunc:  func() {},
		gameEventsChan: make(chan events.GameEvent, 10),
		loadedChan:     make(chan bool, 1),
	}
	wpr.newFSM()
	return wpr
}

func (w *Wrapper) newFSM() {
	w.machine = fsm.NewFSM(
		WrapperOffline,
		wrapperFsmEvents,
		fsm.Callbacks{
			"enter_online": func(ev *fsm.Event) {
				if ev.Src == WrapperStarting {
					select {
					case w.loadedChan <- true:
					default:
					}
				}
			},
			"enter_offline": func(ev *fsm.Event) {
				w.ctxCancelFunc()
			},
		},
	)
}

func (w *Wrapper) processLogEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := w.console.ReadLine()
			if err == io.EOF {
				w.updateState(events.StoppedEvent)
				return
			}

			ev, t := w.parseLineToEvent(line)
			switch t {
			case events.TypeState:
				w.updateState(ev.(events.StateEvent))
			case events.TypeCmd:
				w.handleCmdEvent(ev.(events.GameEvent))
			case events.TypeGame:
				w.handleGameEvent(ev.(events.GameEvent))
			default:
			}
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

func (w *Wrapper) handleGameEvent(ev events.GameEvent) {
	if ev.Is(events.PlayerLeftEvent) {
		delete(w.playerList, ev.Data["player_name"])
	}
	if ev.Is(events.PlayerUUIDEvent) {
		w.playerList[ev.Data["player_name"]] = ev.Data["player_uuid"]
	}
	select {
	case w.gameEventsChan <- ev:
	default:
	}
}

func (w *Wrapper) writeToConsole(cmd string) error {
	if !w.machine.Is(WrapperOnline) {
		return ErrWrapperNotOnline
	}
	return w.console.WriteCmd(cmd)
}

func (w *Wrapper) processClock(ctx context.Context) {
	w.clock.start(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.clock.requestSync():
			w.clock.resetLastSync()
			w.writeToConsole("time query daytime")
		}
	}
}

func (w *Wrapper) processCmdToEvent(cmd string, timeout time.Duration, evs ...string) (events.GameEvent, error) {
	gchns := make([]<-chan events.GameEvent, len(evs))
	for i, ev := range evs {
		gchns[i] = w.eq.get(ev)
	}

	timeoutCaseIdx := len(evs)
	cases := make([]reflect.SelectCase, timeoutCaseIdx+1)
	for i, ch := range gchns {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}
	cases[timeoutCaseIdx] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(time.After(timeout)),
	}

	if err := w.writeToConsole(cmd); err != nil {
		return events.NilGameEvent, err
	}

	chosen, value, _ := reflect.Select(cases)
	if chosen == timeoutCaseIdx {
		return events.NilGameEvent, ErrWrapperResponseTimeout
	}

	ev := value.Interface().(events.GameEvent)
	errMessage, ok := ev.Data["error_message"]
	if ok {
		// If the game event carries an 'error_message' in its Data field,
		// wrap and propagate the error message as an error.
		return events.NilGameEvent, errors.New(errMessage)
	}
	return ev, nil
}

func (w *Wrapper) processCmdToEventArr(cmd string, timeout time.Duration, ev string) ([]events.GameEvent, error) {
	evChan := w.eq.get(ev)
	if err := w.writeToConsole(cmd); err != nil {
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

func (w *Wrapper) Ban(player, reason string) error {
	cmd := strings.Join([]string{"ban", player, reason}, " ")
	return w.writeToConsole(cmd)
}

func (w *Wrapper) BanList(t BanListType) ([]string, error) {
	cmd := fmt.Sprintf("banlist %s", t)
	evs, err := w.processCmdToEventArr(cmd, 3*time.Second, events.BanList)
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
	ev, err := w.processCmdToEvent(cmd, 3*time.Second, events.DataGet)
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
	return w.writeToConsole(cmd)
}

// DeOp removes a given player from the operator list.
func (w *Wrapper) DeOp(player string) error {
	return w.writeToConsole("deop " + player)
}

// Difficulty changes the game difficulty level of the world.
func (w *Wrapper) Difficulty(d GameDifficulty) error {
	cmd := fmt.Sprintf("difficulty %s", d)
	_, err := w.processCmdToEvent(cmd, 1*time.Second, events.Difficulty)
	return err
}

// ExperienceAdd adds a given amount of experience either:
// - levels or
// - points
// to the provided player.
func (w *Wrapper) ExperienceAdd(target string, xp int32, xpType ExperienceType) error {
	cmd := fmt.Sprintf("experience add %s %d %s", target, xp, xpType)
	ev, err := w.processCmdToEvent(cmd, 1*time.Second, events.ExperienceAdd, events.NoPlayerFound)
	if err != nil {
		return err
	}
	if ev.Is(events.NoPlayerFoundEvent) {
		return ErrPlayerNotFound
	}
	return nil
}

// ExperienceQuery returns the amount of experience of the provided player.
// The 'target' arg should be a single target, multi-targets query might fail.
func (w *Wrapper) ExperienceQuery(target string, xpType ExperienceType) (int, error) {
	cmd := fmt.Sprintf("experience query %s %s", target, xpType)
	ev, err := w.processCmdToEvent(cmd, 1*time.Second, events.ExperienceQuery, events.NoPlayerFound)
	if err != nil {
		return 0, err
	}
	if ev.Is(events.NoPlayerFoundEvent) {
		return 0, ErrPlayerNotFound
	}
	return strconv.Atoi(ev.Data["amount"])
}

// GameEvents returns a receive-only channel of game related event. For example:
// - Player joined, left, died, was banned.
// - Game updates like game mode changes.
// - Player sends messages...
func (w *Wrapper) GameEvents() <-chan events.GameEvent {
	return w.gameEventsChan
}

// Give give a target player entity some given items.
func (w *Wrapper) Give(target, item string, count int) error {
	cmd := fmt.Sprintf("give %s %s %d", target, item, count)
	ev, err := w.processCmdToEvent(cmd, 1*time.Second, events.Give, events.NoPlayerFound, events.UnknownItem)
	if err != nil {
		return err
	}
	if ev.Is(events.NoPlayerFoundEvent) {
		return ErrPlayerNotFound
	}
	if ev.Is(events.UnknownItemEvent) {
		return ErrUnknownItem
	}
	return nil
}

// Kill the java process, use with caution since it will not trigger a save game.
// Kill manually perform some cleanup task and hard reset the state to 'offline'.
func (w *Wrapper) Kill() error {
	if err := w.console.Kill(); err != nil {
		return err
	}

	// Hard reset the wrapper machine state the 'offline'.
	w.machine.SetState(WrapperOffline)
	// Manually trigger the context cancellation since 'SetState'
	// does not trigger any callbacks on the fsm.
	w.ctxCancelFunc()
	return nil
}

// Kick kicks the provided player from the server. If a reason is provided,
// the message will display on the players screen when disconnected.
func (w *Wrapper) Kick(target, reason string) error {
	cmd := strings.Join([]string{"kick", target, reason}, " ")
	ev, err := w.processCmdToEvent(cmd, 1*time.Second, events.Kicked, events.NoPlayerFound)
	if err != nil {
		return err
	}
	if ev.Is(events.NoPlayerFoundEvent) {
		return ErrPlayerNotFound
	}
	return nil
}

// List returns a list of connected players on the server.
func (w *Wrapper) List() []Player {
	players := []Player{}
	for name, uuid := range w.playerList {
		players = append(players, Player{
			Name: name,
			UUID: uuid,
		})
	}
	return players
}

func (w *Wrapper) Loaded() <-chan bool {
	return w.loadedChan
}

// Reload reloads the server datapack.
func (w *Wrapper) Reload() error {
	return w.writeToConsole("reload")
}

// SaveAll marks all chunks and player data to be saved to the data storage device.
// When flush is true, the marked data are saved immediately.
func (w *Wrapper) SaveAll(flush bool) error {
	cmd := "save-all"
	if flush {
		cmd += " flush"
	}
	return w.writeToConsole(cmd)
}

// SaveOn enables automatic saving. The server is allowed to write to the world files.
func (w *Wrapper) SaveOn() error {
	return w.writeToConsole("save-on")
}

// SaveOff disables automatic saving by preventing the server from writing to the world files.
func (w *Wrapper) SaveOff() error {
	return w.writeToConsole("save-off")
}

// Say sends the given message in the minecraft in-game chat.
func (w *Wrapper) Say(msg string) error {
	return w.writeToConsole("say " + msg)
}

// Seed returns the world seed.
func (w *Wrapper) Seed() (int, error) {
	ev, err := w.processCmdToEvent("seed", 1*time.Second, events.Seed)
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
	if !w.machine.Is(WrapperOffline) {
		return fmt.Errorf("cannot Start when wrapper is in %s state", w.State())
	}
	ctx, cancel := context.WithCancel(context.Background())
	w.ctxCancelFunc = cancel
	go w.processLogEvents(ctx)
	go w.processClock(ctx)
	return w.console.Start()
}

// State returns the current state of the server, it can be one of:
// 'offline', 'online', 'starting' or 'stopping'.
func (w *Wrapper) State() string {
	return w.machine.Current()
}

// Stop pipes a 'stop' command to the minecraft java process.
func (w *Wrapper) Stop() error {
	if !w.machine.Is(WrapperOnline) {
		return ErrWrapperNotOnline
	}
	return w.console.WriteCmd("stop")
}

// Tell sends a message to a specific target in the server.
func (w *Wrapper) Tell(target, msg string) error {
	cmd := fmt.Sprintf("tell %s %s", target, msg)
	ev, err := w.processCmdToEvent(cmd, 3*time.Second, events.WhisperTo, events.NoPlayerFound)
	if err != nil {
		return err
	}
	if ev.Is(events.NoPlayerFoundEvent) {
		return ErrPlayerNotFound
	}
	return nil
}

// Tick returns the current minecraft game tick, which runs at a fixed rate
// of 20 ticks per second, src: https://minecraft.gamepedia.com/Tick.
func (w *Wrapper) Tick() int {
	return w.clock.Tick
}
