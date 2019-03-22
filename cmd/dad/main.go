package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/alee792/dad/pkg/getzit"

	"github.com/alee792/dad/pkg/dad"
)

func main() {
	ctx := context.Background()
	cfg := dad.Config{}
	chain := dad.NewChain(cfg)
	r := strings.NewReader("hello bye see you")
	chain.Read(ctx, r)

	c := getzit.NewGraphQLClient(getzit.Config{
		Addr: "https://icanhazdadjoke.com/graphql",
	})
	j, err := c.GetJoke(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(j)
}
