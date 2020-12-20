# minecraft-wrapper

<p align="center">
  <img src="https://github.com/wlwanpan/minecraft-wrapper/blob/master/assets/minecraft-gopher.png?raw=true" alt="Minecraft Gopher"/>
</p>

[![GoDoc](https://godoc.org/github.com/wlwanpan/minecraft-wrapper?status.svg)](https://godoc.org/github.com/wlwanpan/minecraft-wrapper)

## What is minecraft-wrapper?

Wrapper is a go package that wraps a Minecraft Server and interacts with it by pushing in commands and reading the server logs.

## Usage

```go
wpr := wrapper.NewDefaultWrapper("server.jar", 1024, 1024)
wpr.Start()

// Listening to game events...
for {
  select {
  case e := <-wpr.GameEvents():
    log.Println(e.String())
  }
}
```
