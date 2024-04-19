package client

import (
	"fmt"
	"time"
)

var (
	RegisteredChains map[string]func()
)

func init() {
	RegisteredChains = make(map[string]func())
	RegisteredChains["eth"] = func() {
		fmt.Println("hello eth")
		time.Sleep(time.Second * 2)
	}
	RegisteredChains["avalanche"] = func() {
		fmt.Println("hello avalanche")
		time.Sleep(time.Second * 2)
	}
}
