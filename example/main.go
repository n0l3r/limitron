package main

import "github.com/n0l3r/limitron"

func main() {
	cfg := limitron.Config{
		Rate:      1, // 1 token leaks per second
		Capacity:  3,
		StoreType: limitron.StoreTypeMemory,
	}

	lb, _ := limitron.NewLeakyBucketLimiter(cfg)

	// Use the limiter (example)
	for i := 0; i < 10; i++ {
		allowed, _ := lb.Allow("example-client")
		if allowed {
			println("Request allowed")
		} else {
			println("Request denied")
		}
	}
}
