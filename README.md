# minecraft-wrapper

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/minecraft-gopher.png?raw=true" alt="Minecraft Gopher"/>
</p>

[![GoDoc](https://godoc.org/github.com/wlwanpan/minecraft-wrapper?status.svg)](https://godoc.org/github.com/wlwanpan/minecraft-wrapper)
[![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoicmdSUjltNjdIODN0dFNQQXgzUUZHajB4WnFxbVVzWDlMOW41VnYvZ2pTUEN5MnBKR1djVUtwNzdraFlNblQyV01HSldGY2w1OXhIZDljOGRqYzlyU3NRPSIsIml2UGFyYW1ldGVyU3BlYyI6IlJieFV3NjZycnM5MGo2QVYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)](https://us-west-2.console.aws.amazon.com/codesuite/codebuild/597927659010/projects/minecraft-wrapper)

## What is minecraft-wrapper?

Wrapper is a go package that wraps a Minecraft Server (JE) and interacts with it by pushing in commands and reading the server logs. This package is meant to be used as an interface between your minecraft server and your go application.

## Installation

```bash
go get github.com/wlwanpan/minecraft-wrapper
```

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

- Triggers the running game to save immediately.
```go
if err := wpr.SaveAll(true); err != nil {
  ...
}
```

Note: This package is developed and tested on Minecraft 1.16, though most functionalities (`Start`, `Stop`, `Seed`, ...) works across all versions. Commands like `/data get` was introduced in version 1.13 and might not work for earlier versions. :warning: 

## Overview

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/architecture.png?raw=true" alt="Minecraft Wrapper Overview"/>
</p>

If you are interested in learning the basic inner working of the wrapper, you can check out my [Medium article](https://levelup.gitconnected.com/lets-build-a-minecraft-server-wrapper-in-go-122c087e0023) for more details.

## Commands :construction:

All the following methods/commands are from the official [list of commands](https://minecraft.gamepedia.com/Commands#List_and_summary_of_commands) unless otherwise specified.

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
- [x] [DeOp](https://minecraft.gamepedia.com/Commands/deop)
- [x] [Difficulty](https://minecraft.gamepedia.com/Commands/difficulty)
- [ ] [Effect](https://minecraft.gamepedia.com/Commands/effect)
- [ ] [Enchant](https://minecraft.gamepedia.com/Commands/enchant)
- [ ] [Experience](https://minecraft.gamepedia.com/Commands/experience)
- [ ] [Fill](https://minecraft.gamepedia.com/Commands/fill)
- [ ] [ForceLoad](https://minecraft.gamepedia.com/Commands/forceload)
- [ ] [Function](https://minecraft.gamepedia.com/Commands/function)
- [x] [GameEvents](https://pkg.go.dev/github.com/wlwanpan/minecraft-wrapper#Wrapper.GameEvents) (Unofficial)
- [ ] [GameMode](https://minecraft.gamepedia.com/Commands/gamemode)
- [ ] [GameRule](https://minecraft.gamepedia.com/Commands/gamerule)
- [x] [SaveAll](https://minecraft.gamepedia.com/Commands/save#save-all)
- [x] [SaveOff](https://minecraft.gamepedia.com/Commands/save#save-off)
- [x] [SaveOn](https://minecraft.gamepedia.com/Commands/save#save-on)
- [x] [Say](https://minecraft.gamepedia.com/Commands/say)
- [x] [Seed](https://minecraft.gamepedia.com/Commands/seed)
- [x] [Start](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Start) (Unofficial)
- [x] [State](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.State) - (Unofficial)
- [x] [Stop](https://minecraft.gamepedia.com/Commands/stop)
- [x] [Kill](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Kill) - Terminates the Java Process (Unofficial)
- [x] [Tick](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Tick) (Unofficial)

This list might be incomplete...

## Minecraft resources

- [Gamepedia](https://minecraft.gamepedia.com)
- [DigMinecraft](https://www.digminecraft.com/game_commands)

## Help and contributions

Feel free to drop a PR, file an issue or proposal of changes you want to have.
