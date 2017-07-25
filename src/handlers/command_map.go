package handlers

import (
	"github.com/3stadt/GoTBot/src/structs"
)

var CommandMap = map[string]func(channel string, sender string, params string) (*structs.Message, error) {
	"goSay": Echo,
	"slap":  Slap,
	"stats": Stats,
}