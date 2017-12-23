package scenario2

import (
	"travel"
)

// Start ..
func Start() [10]travel.Summary {
	c := travel.GetCustomerDetails()
	ds := travel.GetDestination(c)

	quoteDone := [10]chan struct{}{}
	weatherDone := [10]chan struct{}{}

	for i := range quoteDone {
		quoteDone[i] = make(chan struct{})
	}
	for i := range weatherDone {
		weatherDone[i] = make(chan struct{})
	}

	var summaries [10]travel.Summary

	for i, d := range ds {
		idx := i
		dest := d
		go func() {
			summaries[idx].Quote = travel.GetQuote(dest)
			quoteDone[idx] <- struct{}{}
		}()
		go func() {
			summaries[idx].Weather = travel.GetWeatherForcast(dest)
			weatherDone[idx] <- struct{}{}
		}()
	}

	for _, done := range quoteDone {
		<-done
	}
	for _, done := range weatherDone {
		<-done
	}
	return summaries
}
