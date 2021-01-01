# SNBT

SNBT (Stringified Name Binary Tag) is a file format introduced by Minecraft to save its data. While there are some already exist good packages to decode those files, the server logs prints its stringified counterpart. This lightweight package is meant to decode a given SNBT to a Go struct.

## Basic Usage
```go
bytesToDecode := []byte(`{Base: 1.0d, Name: "minecraft:generic.attack_damage"}`)

bytesStruct := struct {
    Base float64
    Name string
}{}

snbt.Decode(bytesToDecode, &bytesStruct)

fmt.Printf("%+v", bytesStruct) // {Base:1 Name:minecraft:generic.attack_damage}
```

## Resources
- https://minecraft.gamepedia.com/NBT_format
