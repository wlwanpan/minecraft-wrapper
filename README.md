# minecraft-wrapper

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/minecraft-gopher.png?raw=true" alt="Minecraft Gopher"/>
</p>

[![GoDoc](https://godoc.org/github.com/wlwanpan/minecraft-wrapper?status.svg)](https://godoc.org/github.com/wlwanpan/minecraft-wrapper)

## What is minecraft-wrapper?

Wrapper is a go package that wraps a Minecraft Server and interacts with it by pushing in commands and reading the server logs. This package is meant to be used as an interface between your minecraft server and your go application. 

## Usage

- Usage with default configs.
```go
wpr := wrapper.NewDefaultWrapper("server.jar", 1024, 1024)
wpr.Start()
defer wrp.Stop()

// Listening to game events...
for {
  select {
  case e := <-wpr.GameEvents():
    log.Println(e.String())
  }
}
```

- Listening to the `Wrapper` state changes.
```go
wpr.RegisterStateChangeCBs(func (w *wrapper.Wrapper, from events.Event, to events.Event) {
	log.Printf("%s -> %s", from.String(), to.String())
})
```

## Overview

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/architecture.png?raw=true" alt="Minecraft Wrapper Overview"/>
</p>

If you are interested in learning about inner working of the wrapper, you can check out my [Medium article](https://levelup.gitconnected.com/lets-build-a-minecraft-server-wrapper-in-go-122c087e0023) for more details.

## Minecraft resources

- [Daylight cycle](https://minecraft.gamepedia.com/Daylight_cycle)

## Help and contributions

Feel free to drop a PR, file an issue or proposal of changes you want to have.
