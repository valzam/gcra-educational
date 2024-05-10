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
	runBurst()
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
	burst := float64(3) // in addition to one every 100ms you can make 3 requests at any time

	r := rlburst{}
	// Use 1 burst request
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
	time.Sleep(100 * time.Millisecond)

	// regular request
	r.use(rate, time.Second, burst)
	time.Sleep(100 * time.Millisecond)

	// use 2 burst requests + one regular request
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)

	// rate limited
	r.use(rate, time.Second, burst)
	time.Sleep(100 * time.Millisecond)

	// use one regular request, second request is rate limited, no burst available
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)

	// wait for 1 second theoretical period to be over
	time.Sleep(700 * time.Millisecond)

	// can burst 3 requests plus 1 regular request again
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
	r.use(rate, time.Second, burst)
}

func runBurstCost() {
	rate := float64(10)
	burst := float64(2)

	cost := float64(1)

	r := rlburstcost{}
	r.use(rate, time.Second, burst, cost)
	r.use(rate, time.Second, burst, cost)
	r.use(rate, time.Second, burst, cost)
	time.Sleep(100 * time.Millisecond)
	r.use(rate, time.Second, burst, 2)
	//for range 10 {
	//}
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
