// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package money

import (
    "fmt"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Money represents an amount of money with a specific currency.
type Money struct {
    currency string
    amount   int64
}

// Add adds two amounts of money with the same currency.
func (m Money) Add(other Money) (Money, error) {
    // Start a Datadog span
    span, _ := tracer.StartSpanFromContext(context.Background(), "money.add")
    defer span.Finish()

    // Check if the currencies match
    if m.currency != other.currency {
        return Money{}, fmt.Errorf("cannot add two different currencies")
    }

    // Set price as a span tag
    span.SetTag("price.amount", m.amount)
    span.SetTag("price.currency", m.currency)

    // Perform addition
    return Money{
        currency: m.currency,
        amount:   m.amount + other.amount,
    }, nil
}

// Multiply multiplies the amount of money by a factor.
func (m Money) Multiply(factor int64) Money {
    // Start a Datadog span
    span, _ := tracer.StartSpanFromContext(context.Background(), "money.multiply")
    defer span.Finish()

    // Set price as a span tag
    span.SetTag("price.amount", m.amount)
    span.SetTag("price.currency", m.currency)
    span.SetTag("factor", factor)

    // Perform multiplication
    return Money{
        currency: m.currency,
        amount:   m.amount * factor,
    }
}
