package main

import (
	"fmt"
	"time"
)

// tat = theoretical arrival time
type rlbasic struct {
	tat float64
}

func (r *rlbasic) use(rate float64, p time.Duration) res {
	defer func() {
		fmt.Println("------------------------------------------------")
	}()

	// Represent the expected arrival rate of requests
	// e.g. a rate of 1 second and 10 req/s limit would imply that we allow a request every 0.1 seconds
	period := float64(p.Milliseconds())
	interval := period / rate
	fmt.Printf("interval: %f\n", interval)

	// By default nextAllowedRequestTime is now for first-time requests
	now := float64(time.Now().UnixMilli())
	fmt.Printf("now: %f\n", now)
	tat := now
	fmt.Printf("r.nextAllowedRequestTime: %f\n", r.tat)
	// If nextAllowedRequestTime exists and is in the future the client has already used part of their rate limit
	if r.tat != 0 && r.tat > now {
		tat = r.tat
	} else {
		r.tat = 0 // expire old nextAllowedRequestTime
	}
	fmt.Printf("nextAllowedRequestTime: %f\n", tat)

	// newTat is the theoretical point in the future at which the next request is allowed
	newTat := tat + 1*interval
	fmt.Printf("newTat: %f\n", newTat)

	// Calculate whether the client has waited for at least 1 interval
	diff := now - tat
	fmt.Printf("diff: %f\n", diff)
	remaining := diff / interval
	fmt.Printf("remaining: %f\n", remaining)

	// it is very important to
	// - only check < 0 since == 0 is a valid state, 0 means we make request exactly at interval
	// - return if there are no remaining. If newTat gets stored even if remaining < 0 the rate limit would never reset
	if remaining < 0 {
		fmt.Printf("hit rate limit\n")
		return res{}
	} else {
		// If there is remaining rate limit store it
		// Represents how many milliseconds until the next request is allowed
		r.tat = newTat
		resetAfter := newTat - now
		fmt.Printf("resetAfter: %f\n", resetAfter)
		return res{}
	}
}
