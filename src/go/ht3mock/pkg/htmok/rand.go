package htmok

import (
	"math/rand"
	"time"
)

func rnd(min, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Int()%(max-min)
}

func rndDelay(min, max int) time.Duration {
	return time.Duration(rnd(min, max)) * time.Millisecond
}
