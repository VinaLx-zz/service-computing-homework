package scenario1

import "travel"

// Start ..
func Start() [10]travel.Summary {
	c := travel.GetCustomerDetails()
	ds := travel.GetDestination(c)
	var summaries [10]travel.Summary
	for i, d := range ds {
		q := travel.GetQuote(d)
		w := travel.GetWeatherForcast(d)
		summaries[i] = travel.Summary{Destination: d, Quote: q, Weather: w}
	}
	return summaries
}
