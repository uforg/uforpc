package main

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/version"
)

type cmdVersionArgs struct{}

func cmdVersion(_ *cmdVersionArgs) {
	fmt.Printf("UFO RPC %s\n", version.VersionWithPrefix)
	fmt.Println("Please star the repo at https://github.com/uforg/uforpc")
}
