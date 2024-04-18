package main

import (
	"context"
	"fmt"

	"github.com/matthiasbruns/gnutella-go/cache"
)

func main() {
	ctx := context.Background()
	ips := cache.QueryBazookaVendors(ctx)

	fmt.Printf("Found %d ips\n", len(ips))
}
