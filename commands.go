package wrapper

// Reach objective: verify minecraft version (Bedrock vs. Java and 1.XX) to not build / prevent using certain commands

type Command interface {
	Command() string
	Events() []Event
}

// Could be a great spot to use the github.com/densestvoid/postoffice package.
// It was designed to be able to send and receive on channels identified by interface addresses.
// Each event type could be registered as an address
// Event is for any console resposne, error is for command processing only
func (w *Wrapper) ExecuteCommand(cmd Command) (Event, error) {
	// TODO: create/get channels for each event type on the wrapper

	// TODO: write the command to the console

	// TODO: wait to receive on one of the event channels, and return that event

	return nil, nil
}

/*
attribute
advancement
ban x
ban-ip
banlist x
bossbar
clear
clone
data (get)
datapack
debug
defaultgamemode x
deop x
difficulty x
effect
enchant
execute
experience (add,query)
fill
forceload
function
gamemode
gamerule
give x
help
kick x
kill
list x
locate
locatebiome
loot
me
msg
op
pardon
particle
playsound
publish
recipe
reload
save-all x
save-off x
save-on x
say x
schedule
scoreboard
seed
setblock
setidletimeout
setworldspawn
spawnpoint
spectate
spreadplayers
stop x
stopsound
summon
tag
team
teammsg
teleport
tell x
tellraw
time
title
tp
trigger
w
weather
whitelist
worldborder
xp
*/
