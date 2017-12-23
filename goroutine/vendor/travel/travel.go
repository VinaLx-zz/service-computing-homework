package travel

import (
	"time"
)

// Customer ..
type Customer struct{}

// Destination ..
type Destination struct{}

// Quote ..
type Quote struct{}

// Weather ..
type Weather struct{}

// Summary ..
type Summary struct {
	Destination Destination
	Quote       Quote
	Weather     Weather
}

// GetCustomerDetails ..
func GetCustomerDetails() Customer {
	time.Sleep(150 * time.Millisecond)
	return Customer{}
}

// GetDestination ..
func GetDestination(c Customer) [10]Destination {
	time.Sleep(250 * time.Millisecond)
	return [10]Destination{}
}

// GetWeatherForcast ..
func GetWeatherForcast(d Destination) Weather {
	time.Sleep(330 * time.Millisecond)
	return Weather{}
}

// GetQuote ..
func GetQuote(d Destination) Quote {
	time.Sleep(170 * time.Millisecond)
	return Quote{}
}
