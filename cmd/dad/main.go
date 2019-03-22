package main

import (
	"strings"

	"github.com/alee792/dad"
)

func main() {
	cfg := dad.Config{}
	chain := dad.NewChain(cfg)
	r := strings.NewReader("hello bye see you")
	chain.Read(r)
}
