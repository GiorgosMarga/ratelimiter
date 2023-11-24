# GoRateLimiter

GoRateLimiter is a simple and efficient rate limiter library written in Go. Rate limiting is a critical aspect of distributed systems to control the rate of incoming requests or events. This library provides a straightforward implementation of a rate limiter using the token bucket algorithm.

## What is a Rate Limiter?

A rate limiter is a mechanism used to control the rate at which operations or events are allowed to occur. It is commonly employed to prevent abuse, manage resources, and ensure a more predictable and stable system behavior.

## Why Use a Rate Limiter?

Rate limiting is essential in various scenarios, such as:

1. **Preventing Abuse:** Protect your services from abuse, such as brute force attacks or excessive API requests.
2. **Resource Management:** Ensure fair and efficient use of resources by limiting the rate of resource-intensive operations.
3. **Stability:** Avoid overloading your system during traffic spikes by controlling the rate of incoming requests.
4. **Compliance:** Adhere to API rate limits imposed by external services or comply with regulations.

## Bucket Algorithm Pros and Cons

### Pros:

1. **Simplicity:** The bucket algorithm is straightforward to implement and understand.
2. **Predictable:** Provides a predictable and stable rate limiting mechanism.
3. **Dynamic Adjustments:** Allows dynamic adjustments to the rate by modifying the bucket size or refill rate.

### Cons:

1. **Burstiness:** The bucket algorithm can exhibit burstiness when requests arrive in clusters, as it allows bursts of requests when the bucket is not empty.
2. **Precision vs. Efficiency Trade-off:** Achieving high precision in limiting rates might require smaller bucket sizes, leading to increased memory and computation overhead.
3. **Clock Dependency:** The bucket algorithm is dependent on the clock for timekeeping, which can introduce complexities in distributed systems.

## Usage

```go
package main

import (
	"fmt"
	"time"
    "http"
	"github.com/GiorgosMarga/ratelimiter"
)

type Server struct {
	Port        string
	Ratelimiter *ratelimiter.RateLimiter
}

func main() {
	bucketCapacity := 10 // Initial and Max capacity of each bucket -> 10 requests per refillTime
	refillTime := 60     // Every 10 seconds add refillRate tokens to each bucket
	refillRate := 5      // Number of tokens to add to each bucket every refillTime
	infoLog := true      // For debug and development environments
	s := &Server{
		Port:        ":3000",
		Ratelimiter: ratelimiter.NewRateLimiter(bucketCapacity, refillTime, refillRate, infoLog),
	}
	http.HandleFunc("/", s.Home)

	log.Fatal(http.ListenAndServe(s.Port, nil))
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	fmt.Println(ip)
	if ok := s.Ratelimiter.CheckRequest(ip); !ok {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	// Process request
	w.WriteHeader(http.StatusOK)

}
```