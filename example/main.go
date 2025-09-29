package main

import (
	"fmt"
	"github.com/Blocktunium/gonyx/contrib/gormkit"
)

func main() {
	fmt.Println("Using gormkit:", gormkit.GetManager())
}
