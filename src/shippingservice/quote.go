package main

import (
	"fmt"
	"math"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Quote represents a price with dollars and cents.
type Quote struct {
	Dollars uint32
	Cents   uint32
}

// CreateQuoteFromCount takes a number of items and returns a Quote struct.
func CreateQuoteFromCount(count int) Quote {
	return CreateQuoteFromFloat(8.99)
}

// CreateQuoteFromFloat takes a price represented as a float and creates a Quote struct.
func CreateQuoteFromFloat(value float64) Quote {
	span := tracer.StartSpan("quote.create")
	defer span.Finish()

	units, fraction := math.Modf(value)
	quote := Quote{
		uint32(units),
		uint32(math.Trunc(fraction * 100)),
	}

	span.SetTag("quote.dollars", quote.Dollars)
	span.SetTag("quote.cents", quote.Cents)
	span.SetTag("quote.total", fmt.Sprintf("%d.%02d", quote.Dollars, quote.Cents))

	return quote
}

func main() {
	tracer.Start()
	defer tracer.Stop()

	quote := CreateQuoteFromCount(5)
	fmt.Printf("Quote: %d dollars and %d cents\n", quote.Dollars, quote.Cents)
}
