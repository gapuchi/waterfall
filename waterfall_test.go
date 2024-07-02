package main

import (
	"github.com/gapuchi/waterfall/v2/balance"
	"github.com/gapuchi/waterfall/v2/util"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestCalcRocAndPrefOwed(t *testing.T) {

	runDate := time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)

	t.Run("succeeds", func(t *testing.T) {
		txns := []util.Transaction{
			{
				CommitmentId: "1",
				Amount:       balance.NewBalanceFromDollarAmount(100),
				Date:         time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
			},
			{
				CommitmentId: "1",
				Amount:       balance.NewBalanceFromDollarAmount(300),
				Date:         time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC),
			},
			{
				CommitmentId: "2",
				Amount:       balance.NewBalanceFromDollarAmount(1_000),
				Date:         time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
			},
		}

		roc, pref := CalcRocAndPrefOwed(txns, runDate)

		expectedRoc := map[string]balance.Balance{
			"1": balance.NewBalanceFromDollarAmount(400),
			"2": balance.NewBalanceFromDollarAmount(1_000),
		}

		expectedPref := map[string]balance.Balance{
			"1": balance.NewBalanceFromDollarAmount(21.85),
			"2": balance.NewBalanceFromDollarAmount(80),
		}

		assert.Equal(t, expectedRoc, roc)
		assert.Equal(t, expectedPref, pref)
	})
}

func TestAllocate(t *testing.T) {
	t.Run("Enough distribution for all tiers", func(t *testing.T) {

		distributions := []util.Transaction{
			{
				CommitmentId: "four tiers",
				Amount:       balance.NewBalanceFromDollarAmount(2_000),
			},
			{
				CommitmentId: "three tiers",
				Amount:       balance.NewBalanceFromDollarAmount(1_100),
			},
			{
				CommitmentId: "two tiers",
				Amount:       balance.NewBalanceFromDollarAmount(1_070.78),
			},
			{
				CommitmentId: "one tier",
				Amount:       balance.NewBalanceFromDollarAmount(900),
			},
		}

		rocOwed := map[string]balance.Balance{
			"four tiers":  balance.NewBalanceFromDollarAmount(1000),
			"three tiers": balance.NewBalanceFromDollarAmount(1000),
			"two tiers":   balance.NewBalanceFromDollarAmount(1000),
			"one tier":    balance.NewBalanceFromDollarAmount(1000),
		}

		prefOwed := map[string]balance.Balance{
			"four tiers":  balance.NewBalanceFromDollarAmount(80),
			"three tiers": balance.NewBalanceFromDollarAmount(80),
			"two tiers":   balance.NewBalanceFromDollarAmount(80),
			"one tier":    balance.NewBalanceFromDollarAmount(80),
		}

		actual := Allocate(distributions, rocOwed, prefOwed)
		expected := map[string]*util.Allocation{
			"four tiers": {
				ReturnOnCapital: balance.NewBalanceFromDollarAmount(1_000),
				PreferredReturn: balance.NewBalanceFromDollarAmount(80),
				CatchUp:         balance.NewBalanceFromDollarAmount(20),
				FinalSplitLP:    balance.NewBalanceFromDollarAmount(720),
				FinalSplitGP:    balance.NewBalanceFromDollarAmount(180),
			},
			"three tiers": {
				ReturnOnCapital: balance.NewBalanceFromDollarAmount(1_000),
				PreferredReturn: balance.NewBalanceFromDollarAmount(80),
				CatchUp:         balance.NewBalanceFromDollarAmount(20),
			},
			"two tiers": {
				ReturnOnCapital: balance.NewBalanceFromDollarAmount(1_000),
				PreferredReturn: balance.NewBalanceFromDollarAmount(70.78),
			},
			"one tier": {
				ReturnOnCapital: balance.NewBalanceFromDollarAmount(900),
			},
		}

		assert.Equal(t, expected, actual)
	})
}
