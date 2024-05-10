# Background
## Fixed window
Every period the X requests are available
- All requests can be used up instantaneously
- Rate limit doesn't reset until end of period

## Sliding window/log
- Keep track of request counts, either via small bucket or by storing all in a log
- To check rate limit count #requests in last period
- Requests outside of the period need to be deleted
- Allows for burst and refills smoothly

## Token bucket
- Essentially similar to fixed window, X requests in period
- However, replenishment happens continuously, either via separate process or clever use of timestamps to determine how much should get added

## Leaky bucket
- Essentially the inverse of Token bucket, requests can happen until a bucket is full
- Old requests are continuously removed from bucket to free up capacity
- Same as with Token bucket this can either be a background process or clever use of timestamps to determine how much should be removed

# GCRA
Generic Rate Cell Algorithm is essentially a generalisation of both Token Bucket and Leaky Bucket that has some nice properties:
- There exists an implementation (which we will look at) that only needs to store a single timestamp per rate limit key
- The same implementation can be used to simulate different rate limit behaviours just by changing parameters
- Rate limits can easily be changed on the fly
