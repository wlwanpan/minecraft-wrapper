# minecraft-wrapper

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/minecraft-gopher.png?raw=true" alt="Minecraft Gopher"/>
</p>

[![GoDoc](https://godoc.org/github.com/wlwanpan/minecraft-wrapper?status.svg)](https://godoc.org/github.com/wlwanpan/minecraft-wrapper)

## What is minecraft-wrapper?

Wrapper is a go package (under construction :construction:) that wraps a Minecraft Server (JE) and interacts with it by pushing in commands and reading the server logs. This package is meant to be used as an interface between your minecraft server and your go application. 

## Usage

- Usage with default configs:
```go
wpr := wrapper.NewDefaultWrapper("server.jar", 1024, 1024)
wpr.Start()
defer wpr.Stop()

// Listening to game events...
for {
  select {
  case e := <-wpr.GameEvents():
    log.Println(e.String())
  }
}
```

- Listening to the `Wrapper` state changes:
```go
wpr.RegisterStateChangeCBs(func (w *wrapper.Wrapper, from events.Event, to events.Event) {
	log.Printf("%s -> %s", from.String(), to.String())
})
```

- Retrieving a player position from the [`/data get`](https://minecraft.gamepedia.com/Commands/data#get) command:
```go
out, err := wpr.DataGet("entity", PLAYER_NAME|PLAYER_UUID)
if err != nil {
	...
}
fmt.Println(out.Pos) // [PLAYER_X, PLAYER_Y, PLAYER_Z]
```

:warning: This package was developed/tested for Minecraft 1.16, though the basic functionality should work across all version. APIs like `wrapper.DataGet` was introduce as from 1.13 but not tested. :warning: 

## Overview

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/architecture.png?raw=true" alt="Minecraft Wrapper Overview"/>
</p>

If you are interested in learning the basic inner working of the wrapper, you can check out my [Medium article](https://levelup.gitconnected.com/lets-build-a-minecraft-server-wrapper-in-go-122c087e0023) for more details.

## Commands

All the folllowing commands/APIs are from the official [list of commands](https://minecraft.gamepedia.com/Commands#List_and_summary_of_commands) unless otherwise specified.

- [ ] [Attributes](https://minecraft.gamepedia.com/Commands/attribute)
- [ ] [Advancement](https://minecraft.gamepedia.com/Commands/advancement)
- [x] [Ban](https://minecraft.gamepedia.com/Commands/ban)
- [ ] [BanIp](https://minecraft.gamepedia.com/Commands/ban#ban-ip)
- [ ] [BanList](https://minecraft.gamepedia.com/Commands/ban#banlist)
- [ ] [Bossbar](https://minecraft.gamepedia.com/Commands/bossbar)
- [x] [DataGet](https://minecraft.gamepedia.com/Commands/data#get)
- [ ] [DataMerge](https://minecraft.gamepedia.com/Commands/data#merge)
- [ ] [DataModify](https://minecraft.gamepedia.com/Commands/data#modify)
- [ ] [DataRemove](https://minecraft.gamepedia.com/Commands/data#remove)
- [ ] [DefaultGameMode](https://minecraft.gamepedia.com/Commands/defaultgamemode)
- [x] [SaveAll](https://minecraft.gamepedia.com/Commands/save#save-all)
- [x] [Say](https://minecraft.gamepedia.com/Commands/say)
- [x] [Seed](https://minecraft.gamepedia.com/Commands/seed)
- [x] [Start](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Start) (Unofficial)
- [x] [State](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.State) - (Unofficial)
- [x] [Stop](https://minecraft.gamepedia.com/Commands/stop)
- [x] [Kill](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Kill) - Terminates the Java Process (Unofficial)
- [x] [Tick](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Tick) (Unofficial)

More incoming...

## Minecraft resources

- [Daylight cycle](https://minecraft.gamepedia.com/Daylight_cycle)

## Help and contributions

Feel free to drop a PR, file an issue or proposal of changes you want to have.
