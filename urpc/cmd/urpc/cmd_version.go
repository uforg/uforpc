package main

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/version"
)

type cmdVersionArgs struct{}

func cmdVersion(_ *cmdVersionArgs) {
	fmt.Printf("UFO RPC %s\n", version.VersionWithPrefix)
}
