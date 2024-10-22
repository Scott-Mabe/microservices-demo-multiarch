package money

import (
    "fmt"
    "context"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
    "errors"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
)


const (
	nanosMin = -999999999
	nanosMax = +999999999
	nanosMod = 1000000000
)

var (
	ErrInvalidValue        = errors.New("one of the specified money values is invalid")
	ErrMismatchingCurrency = errors.New("mismatching currency codes")
)

// IsValid checks if specified value has a valid units/nanos signs and ranges.
func IsValid(m pb.Money) bool {
	return signMatches(m) && validNanos(m.GetNanos())
}

func signMatches(m pb.Money) bool {
	return m.GetNanos() == 0 || m.GetUnits() == 0 || (m.GetNanos() < 0) == (m.GetUnits() < 0)
}

func validNanos(nanos int32) bool { return nanosMin <= nanos && nanos <= nanosMax }

// IsZero returns true if the specified money value is equal to zero.
func IsZero(m pb.Money) bool { return m.GetUnits() == 0 && m.GetNanos() == 0 }

// IsPositive returns true if the specified money value is valid and is
// positive.
func IsPositive(m pb.Money) bool {
	return IsValid(m) && m.GetUnits() > 0 || (m.GetUnits() == 0 && m.GetNanos() > 0)
}

// IsNegative returns true if the specified money value is valid and is
// negative.
func IsNegative(m pb.Money) bool {
	return IsValid(m) && m.GetUnits() < 0 || (m.GetUnits() == 0 && m.GetNanos() < 0)
}

// AreSameCurrency returns true if values l and r have a currency code and
// they are the same values.
func AreSameCurrency(l, r pb.Money) bool {
	return l.GetCurrencyCode() == r.GetCurrencyCode() && l.GetCurrencyCode() != ""
}

// AreEquals returns true if values l and r are the equal, including the
// currency. This does not check validity of the provided values.
func AreEquals(l, r pb.Money) bool {
	return l.GetCurrencyCode() == r.GetCurrencyCode() &&
		l.GetUnits() == r.GetUnits() && l.GetNanos() == r.GetNanos()
}

// Negate returns the same amount with the sign negated.
func Negate(m pb.Money) pb.Money {
	return pb.Money{
		Units:        -m.GetUnits(),
		Nanos:        -m.GetNanos(),
		CurrencyCode: m.GetCurrencyCode()}
}

// Must panics if the given error is not nil. This can be used with other
// functions like: "m := Must(Sum(a,b))".
func Must(v pb.Money, err error) pb.Money {
	if err != nil {
		panic(err)
	}
	return v
}

// Sum adds two values. Returns an error if one of the values are invalid or
// currency codes are not matching (unless currency code is unspecified for
// both).
func Sum(l, r pb.Money) (pb.Money, error) {
	if !IsValid(l) || !IsValid(r) {
		return pb.Money{}, ErrInvalidValue
	} else if l.GetCurrencyCode() != r.GetCurrencyCode() {
		return pb.Money{}, ErrMismatchingCurrency
	}
	units := l.GetUnits() + r.GetUnits()
	nanos := l.GetNanos() + r.GetNanos()

	if (units == 0 && nanos == 0) || (units > 0 && nanos >= 0) || (units < 0 && nanos <= 0) {
		// same sign <units, nanos>
		units += int64(nanos / nanosMod)
		nanos = nanos % nanosMod
	} else {
		// different sign. nanos guaranteed to not to go over the limit
		if units > 0 {
			units--
			nanos += nanosMod
		} else {
			units++
			nanos -= nanosMod
		}
	}

	return pb.Money{
		Units:        units,
		Nanos:        nanos,
		CurrencyCode: l.GetCurrencyCode()}, nil
}

// MultiplySlow is a slow multiplication operation done through adding the value
// to itself n-1 times.
func MultiplySlow(m pb.Money, n uint32) pb.Money {
	out := m
	for n > 1 {
		out = Must(Sum(out, m))
		n--
	}
	return out
}

// Money represents an amount of money with a specific currency.
type Money struct {
    currency string
    amount   int64
}

// Add adds two amounts of money with the same currency.
func (m Money) Add(ctx context.Context, other Money) (Money, error) {
    // Start a Datadog span for the Add operation
    span, _ := tracer.StartSpanFromContext(ctx, "money.add")
    defer span.Finish()

    // Check if the currencies match
    if m.currency != other.currency {
        return Money{}, fmt.Errorf("cannot add two different currencies")
    }

    // Set span tags for price and currency
    span.SetTag("price.amount", m.amount)
    span.SetTag("price.currency", m.currency)

    result := Money{
        currency: m.currency,
        amount:   m.amount + other.amount,
    }

    // Set span tag for the resulting amount
    span.SetTag("resulting_price.amount", result.amount)
    span.SetTag("resulting_price.currency", result.currency)

    return result, nil
}

// Multiply multiplies the amount of money by a factor.
func (m Money) Multiply(ctx context.Context, factor int64) Money {
    // Start a Datadog span for the Multiply operation
    span, _ := tracer.StartSpanFromContext(ctx, "money.multiply")
    defer span.Finish()

    // Set span tags for price and currency
    span.SetTag("price.amount", m.amount)
    span.SetTag("price.currency", m.currency)
    span.SetTag("factor", factor)

    result := Money{
        currency: m.currency,
        amount:   m.amount * factor,
    }

    // Set span tag for the resulting amount
    span.SetTag("resulting_price.amount", result.amount)
    span.SetTag("resulting_price.currency", result.currency)

    return result
}
