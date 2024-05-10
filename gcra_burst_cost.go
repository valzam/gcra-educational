package main

import (
	"fmt"
	"time"
)

// nextAllowedRequestTime = theoretical arrival time
type rlburstcost struct {
	nextAllowedRequestTime float64
}

func (r *rlburstcost) use(rate float64, p time.Duration, burst float64, cost float64) res {
	defer func() {
		fmt.Println("------------------------------------------------")
	}()

	period := float64(p.Milliseconds())
	interval := period / rate

	// Represents the "virtual" number of requests this invocation has used up in terms of time
	// e.g. 0.1 second intervals, if cost is 5 we have used up 0.5 seconds worth of requests in one go
	increment := interval * cost
	fmt.Printf("increment: %f\n", increment)

	// Burst offset is essentially "negative cost".
	// Burst should either be:
	// - equal to cost. This allows us to gradually use up the rate and it also gradually refills
	// - equal to rate. This allows us to use all of the rate instantly but we retain the gradual refills
	burstOffset := interval * burst
	fmt.Printf("burstOffset: %f\n", burstOffset)

	tat, now := r.getTatAndNow()

	// We now accept a variable cost so we also have to move nextAllowedRequestTime by more than 1 interval
	newTat := tat + increment
	fmt.Printf("newTat: %f\n", newTat)

	allowAt := newTat - burstOffset
	fmt.Printf("allowAt: %f\n", allowAt)

	diff := now - allowAt
	fmt.Printf("diff: %f\n", diff)

	remaining := diff / interval
	fmt.Printf("remaining: %f\n", remaining)

	if remaining < 0 {
		fmt.Printf("hit rate limit\n")
		return res{time.UnixMilli(int64(newTat)), 0}
	}

	resetAfter := newTat - now
	fmt.Printf("resetAfter: %f\n", resetAfter)
	if resetAfter > 0 {
		r.nextAllowedRequestTime = newTat
	}

	return res{}
}

func (r *rlburstcost) getTatAndNow() (float64, float64) {
	now := float64(time.Now().UnixMilli())
	fmt.Printf("now: %f\n", now)
	tat := now
	fmt.Printf("r.nextAllowedRequestTime: %f\n", r.nextAllowedRequestTime)
	if r.nextAllowedRequestTime != 0 && r.nextAllowedRequestTime > now {
		tat = r.nextAllowedRequestTime
	} else {
		r.nextAllowedRequestTime = 0 // expire old nextAllowedRequestTime
	}
	fmt.Printf("nextAllowedRequestTime: %f\n", tat)

	return tat, now
}
