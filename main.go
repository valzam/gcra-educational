package main

import (
	"fmt"
	"time"
)

type res struct {
	newTat    time.Time
	remaining float64
}

// Illustration https://brandur.org/rate-limiting

func main() {
	runBurstCost()
}

func runBasic() {
	rate := float64(10)

	r := rlbasic{}
	for range 10 {
		r.use(rate, time.Second)
		time.Sleep(50 * time.Millisecond)
	}
}

func runBurst() {
	rate := float64(10)
	// the burst acts a bit like an overdraft facility on your debit card.
	// with burst=3 you can go down to -3 on your account and every 100ms you get +1 credit to your account.
	// so you can make 3 requests in quick succession before you get a +1.
	// if you then wait 200ms you will have gotten +2, so are at -1 on your account and can dip into the overdraft again until it hits -3.
	burst := float64(3)

	r := rlburst{}
	// You can burst 3 requests before getting rate limited
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst) // fails

	// sleeping for 100ms gives you a regular request but no burst capacity
	time.Sleep(100 * time.Millisecond)
	r.use(rate, time.Second, burst) // succeeds
	r.use(rate, time.Second, burst) // fails

	// sleeping 200ms gives you back one regular and one burst request
	time.Sleep(200 * time.Millisecond)
	r.use(rate, time.Second, burst) // succeeds
	r.use(rate, time.Second, burst) // succeeds
	r.use(rate, time.Second, burst) // fails

	// sleeping for 300ms gives you back your full burst capacity of 3 requests
	time.Sleep(300 * time.Millisecond)
	r.use(rate, time.Second, burst) // succeeds
	r.use(rate, time.Second, burst) // succeeds
	r.use(rate, time.Second, burst) // succeeds
	r.use(rate, time.Second, burst) // fails
}

func runBurstCost() {
	rate := float64(10)

	// cost just means you can dip into your overdraft with one request where cost > 1 instead of making multiple requests
	// but nothing fundamentally changes with the behaviour of the algorithm
	// to be able to use your full balance in one request set burst = rate
	// this is the same as setting burst = rate but constant burst is easier to understand
	r := rlburstcost{}

	burst := float64(4)
	// You can use 4 capacity right away, in any number of requests
	r.use(rate, time.Second, burst, 1)
	r.use(rate, time.Second, burst, 3)
	r.use(rate, time.Second, burst, 1) // fails

	// Sleeping for 200ms gives you back 2 capacity
	time.Sleep(200 * time.Millisecond)
	r.use(rate, time.Second, burst, 2)
	r.use(rate, time.Second, burst, 1) // fails

	// Even if you sleep for 500ms you can never use more than 4 capacity in on request
	time.Sleep(500 * time.Millisecond)
	r.use(rate, time.Second, burst, 5) // fails
	r.use(rate, time.Second, burst, 4) // succeeds

	// The use case here is rate limiting across longer periods
	// Within one second it probably doesn't make sense to set burst < rate
	// But if you rate limit hourly then preventing the user from burning through their whole rate limit instantly
	// (maybe due to a bug on their side) can be beneficial. It seems like it would be finicky to get right though
	// Set burst too low and we might rate limit too aggressively, set it too high and it defeats the purpose
}

func runExtended() {
	rate := float64(10)
	cost := float64(10)
	burst := float64(10)

	r := rlextended{}
	r.use(rate, time.Second, burst, cost) // use up all
	remain := r.hasCapacity(rate, time.Second, burst, cost)
	fmt.Printf("has capacity: %t\n", remain)
	r.use(rate, time.Second, burst, cost) // will rate limit

	r.refund(rate, time.Second, cost) // refund all
	remain = r.hasCapacity(rate, time.Second, burst, cost)
	fmt.Printf("has capacity: %t\n", remain)

	r.use(rate, time.Second, burst, cost) // will work
}
