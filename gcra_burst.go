package main

import (
	"fmt"
	"time"
)

// tat = theoretical arrival time
type rlburst struct {
	tat float64
}

func (r *rlburst) use(rate float64, p time.Duration, burst float64) res {
	defer func() {
		fmt.Println("------------------------------------------------")
	}()

	period := float64(p.Milliseconds())
	interval := period / rate

	// Burst represents how many requests we can make in one go on top of the
	// You can only burst once per period, after you've exhausted your burst you need to
	// wait for rate limit to become available gradually
	if burst == 0 {
		burst = 1
	}
	burstOffset := interval * burst
	fmt.Printf("burstOffset: %f\n", burstOffset)

	tat, now := r.getTatAndNow()

	newTat := tat + 1*interval
	fmt.Printf("newTat: %f\n", newTat)

	// Pull back that point in the future by the burstOffset, essentially allowing requests to come in earlier'
	// This is where burst = rate comes into play, essentially on the first request if cost = rate this will set
	// allowAt = now, letting through the whole request
	// newTat however is now set to "rate from now", e.g. 1 second in the future. If the client makes another request right away
	// it will be rejected because even with the burst offset allowAt will be 1 second in the future
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
		r.tat = newTat
	}

	return res{}
}

func (r *rlburst) getTatAndNow() (float64, float64) {
	now := float64(time.Now().UnixMilli())
	fmt.Printf("now: %f\n", now)
	tat := now
	fmt.Printf("r.nextAllowedRequestTime: %f\n", r.tat)
	if r.tat != 0 && r.tat > now {
		tat = r.tat
	} else {
		r.tat = 0 // expire old nextAllowedRequestTime
	}
	fmt.Printf("nextAllowedRequestTime: %f\n", tat)

	return tat, now
}
