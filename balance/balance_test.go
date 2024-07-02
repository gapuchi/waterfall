package balance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	t.Run("creation", func(t *testing.T) {
		actual := NewBalanceFromDollarAmount(100.47)
		expected := Balance{AmountInCents: 100_47}
		assert.Equal(t, expected, actual)
	})

	t.Run("add", func(t *testing.T) {
		actual := NewBalanceFromDollarAmount(100.47).Add(NewBalanceFromDollarAmount(100.53))
		expected := Balance{AmountInCents: 201_00}
		assert.Equal(t, expected, actual)
	})

	t.Run("minus", func(t *testing.T) {
		actual := NewBalanceFromDollarAmount(100.47).Minus(NewBalanceFromDollarAmount(50.53))
		expected := Balance{AmountInCents: 49_94}
		assert.Equal(t, expected, actual)
	})

	t.Run("greater than", func(t *testing.T) {
		testCases := []struct {
			first    float64
			second   int64
			expected bool
		}{
			{100.01, 100_00, true},
			{100.00, 100_00, false},
			{99.99, 100_00, false},
		}

		for _, testCase := range testCases {
			actual := NewBalanceFromDollarAmount(testCase.first).GreaterThan(testCase.second)
			assert.Equal(t, testCase.expected, actual)
		}
	})

	t.Run("times", func(t *testing.T) {
		actual := NewBalanceFromDollarAmount(100.47).Times(0.7)
		expected := Balance{AmountInCents: 70_33}
		assert.Equal(t, expected, actual)
	})

	t.Run("to string", func(t *testing.T) {
		actual := NewBalanceFromDollarAmount(100.4734534).ToString()
		expected := "100.47"
		assert.Equal(t, expected, actual)
	})

	t.Run("min", func(t *testing.T) {
		testCases := []struct {
			first    float64
			second   float64
			expected float64
		}{
			{100.01, 100, 100},
			{100.00, 100, 100},
			{99.99, 100, 99.99},
		}

		for _, testCase := range testCases {
			actual := Min(NewBalanceFromDollarAmount(testCase.first), NewBalanceFromDollarAmount(testCase.second))
			assert.Equal(t, NewBalanceFromDollarAmount(testCase.expected), actual)
		}
	})
}
