package main

import (
	"context"
	"log"

	"github.com/looplab/fsm"
)

const (
	Offline  = "offline"
	Online   = "online"
	Starting = "starting"
	Stopping = "stopping"
	Saving   = "saving"
)

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
			"enter_offline": func(e *fsm.Event) {
				m.triggerOfflineCallbacks()
			},
			"enter_online": func(e *fsm.Event) {
				m.triggerOnlineCallbacks()
			},
		},
	)
}

func (m *MSW) RegisterOnlineCallbacks(cbs ...Callback) {
	m.onlineCallbacks = append(m.onlineCallbacks, cbs...)
}

func (m *MSW) RegisterOfflineCallback(cbs ...Callback) {
	m.offlineCallbacks = append(m.offlineCallbacks, cbs...)
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
	if err := m.machine.Transition(); err != nil {
		return err
	}
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
				log.Println(err)
				continue
			}
			if err := m.processLog(line); err != nil {
				log.Println(err)
			}
		}
	}
}

func (m *MSW) processLog(line string) error {
	logLine, err := strToLogLine(line)
	if err != nil {
		return err
	}
	log.Println(logLine.output)
	return nil
}

func (m *MSW) Stop(ctx context.Context) error {
	return nil
}
