package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set your own keypair
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	h2, err := libp2p.New(ctx,
		// Use your own created keypair
		libp2p.Identity(priv),

		// Set your own listen address
		// The config takes an array of addresses, specify as many as you want.
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello World, my second hosts ID is %s\n", h2.ID())
}
