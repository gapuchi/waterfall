package balance

import (
	"math"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Balance defines the behavior for dollar amount operations. This was done so over using native
// int or floats to make it consistent across the application.
type Balance struct {
	AmountInCents int64
}

func NewBalanceFromDollarAmount(amount float64) Balance {
	return Balance{AmountInCents: int64(amount * 100)}
}

func (b Balance) Add(other Balance) Balance {
	return Balance{
		AmountInCents: b.AmountInCents + other.AmountInCents,
	}
}

func (b Balance) Minus(other Balance) Balance {
	return Balance{
		AmountInCents: b.AmountInCents - other.AmountInCents,
	}
}

func (b Balance) GreaterThan(amount int64) bool {
	return b.AmountInCents > amount
}

func (b Balance) Times(multiplier float64) Balance {
	product := float64(b.AmountInCents) * multiplier

	return Balance{
		AmountInCents: int64(math.Round(product)),
	}
}

func (b Balance) ToString() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", float64(b.AmountInCents)/100)
}

func Min(x, y Balance) Balance {
	if x.GreaterThan(y.AmountInCents) {
		return y
	} else {
		return x
	}
}
