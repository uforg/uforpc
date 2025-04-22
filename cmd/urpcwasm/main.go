//go:build js && wasm

package main

import (
	"log"
	"syscall/js"
)

var wrappers map[string]js.Func = map[string]js.Func{
	"cmdFmt":       cmdFmtWrapper(),
	"cmdTranspile": cmdTranspileWrapper(),
}

func main() {
	log.Println("UFO RPC WASM: Initializing...")

	for name, wrapper := range wrappers {
		js.Global().Set(name, wrapper)
	}

	log.Println("UFO RPC WASM: Initialized")
	<-make(chan any)
}
