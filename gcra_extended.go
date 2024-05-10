package main

import (
	"fmt"
	"math"
	"time"
)

// tat = theoretical arrival time
type rlextended struct {
	tat float64
}

func (r *rlextended) use(rate float64, p time.Duration, burst float64, cost float64) res {
	fmt.Println("using rate limit")
	defer func() {
		fmt.Println("------------------------------------------------")
	}()

	interval, increment, burstOffset := r.parseInputs(rate, p, burst, cost)
	tat, now := r.getTatAndNow()

	newTat := tat + increment
	allowAt := newTat - burstOffset

	diff := now - allowAt
	remaining := diff / interval

	if remaining < 0 {
		fmt.Printf("hit rate limit\n")
		return res{time.UnixMilli(int64(newTat)), 0}
	}

	resetAfter := newTat - now
	if resetAfter > 0 {
		r.tat = newTat
	}

	return res{}
}

func (r *rlextended) hasCapacity(rate float64, p time.Duration, burst float64, cost float64) bool {
	interval, increment, burstOffset := r.parseInputs(rate, p, burst, cost)
	tat, now := r.getTatAndNow()

	newTat := tat + increment
	allowAt := newTat - burstOffset

	diff := now - allowAt
	remaining := diff / interval
	remaining = math.Min(0, remaining)

	return remaining >= 0
}

func (r *rlextended) refund(rate float64, p time.Duration, refund float64) res {
	fmt.Println("refunding rate limit")
	defer func() {
		fmt.Println("------------------------------------------------")
	}()

	interval, increment, _ := r.parseInputs(rate, p, 0, refund)
	tat, now := r.getTatAndNow()

	// Refunding essentially means allowing the next request to come earlier
	newTat := tat - increment
	r.tat = newTat

	diff := now - newTat
	remaining := diff / interval
	remaining = math.Min(0, remaining)

	return res{time.UnixMilli(int64(newTat)), remaining}
}

func (r *rlextended) getTatAndNow() (float64, float64) {
	now := float64(time.Now().UnixMilli())
	tat := now
	if r.tat != 0 && r.tat > now {
		tat = r.tat
	} else {
		r.tat = 0 // expire old nextAllowedRequestTime
	}

	return tat, now
}

func (r *rlextended) parseInputs(rate float64, p time.Duration, burst float64, cost float64) (interval float64, increment float64, burstOffset float64) {
	period := float64(p.Milliseconds())
	interval = period / rate
	increment = interval * cost
	burstOffset = interval * burst

	return
}
