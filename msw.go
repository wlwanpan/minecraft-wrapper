package main

import (
	"context"
	"log"
	"regexp"

	"github.com/looplab/fsm"
)

const (
	Offline  = "offline"
	Starting = "starting"
	Online   = "online"
	Stopping = "stopping"
)

var (
	logRegex = regexp.MustCompile(`(\[[0-9:]*\]) \[([A-z #0-9]*)\/([A-z #]*)\](.*)`)
)

type Update struct {
}

func logLineToUpdate(line string) *Update {
	// logData := logRegex.FindAllStringSubmatch(line, 5)

	// log_time = r.group(1)
	// server_thread = r.group(2)
	// log_level = r.group(3)
	// output = r.group(4)
	// for _, i := range logData {
	// 	log.Println(i)
	// }
	return nil
}

type Callback func(*MSW) error

type MSW struct {
	console *Console
	machine *fsm.FSM

	onlineCallbacks  []Callback
	offlineCallbacks []Callback
}

func NewMSW(console *Console) *MSW {
	m := &MSW{
		console: console,
	}
	m.initMachine()
	return m
}

func (m *MSW) initMachine() {
	m.machine = fsm.NewFSM(
		Offline,
		fsm.Events{
			fsm.EventDesc{
				Name: Offline,
				Src: []string{
					Starting,
					Online,
					Stopping,
				},
				Dst: Starting,
			},
			fsm.EventDesc{
				Name: Starting,
				Src: []string{
					Offline,
				},
				Dst: Online,
			},
			fsm.EventDesc{
				Name: Online,
				Src: []string{
					Starting,
				},
				Dst: Stopping,
			},
		},
		fsm.Callbacks{
			"enter_offline": func(e *fsm.Event) { m.triggerOfflineCallbacks() },
			"enter_online":  func(e *fsm.Event) { m.triggerOnlineCallbacks() },
		},
	)
}

func (m *MSW) triggerOfflineCallbacks() {
	for _, cb := range m.offlineCallbacks {
		cb(m)
	}
}

func (m *MSW) triggerOnlineCallbacks() {
	for _, cb := range m.onlineCallbacks {
		cb(m)
	}
}

func (m *MSW) State() string {
	return m.machine.Current()
}

func (m *MSW) Start(ctx context.Context) error {
	if err := m.console.Start(); err != nil {
		return err
	}
	go m.processConsoleStdout(ctx)
	return nil
}

func (m *MSW) processConsoleStdout(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := m.console.Read()
			if err != nil {
				// log.Println(err)
				continue
			}
			log.Println(line)
			update := logLineToUpdate(line)
			m.processUpdate(update)
		}
	}
}

func (m *MSW) processUpdate(update *Update) error {
	return nil
}

func (m *MSW) Stop(ctx context.Context) error {
	return nil
}
