package main

import (
	"math"
	"time"

	"github.com/gapuchi/waterfall/v2/balance"
	"github.com/gapuchi/waterfall/v2/util"
)

func Allocate(distributions []util.Transaction, rocOwed, prefOwed map[string]balance.Balance) map[string]*util.Allocation {
	allocations := make(map[string]*util.Allocation)
	for _, distribution := range distributions {
		totalAllocation, ok := allocations[distribution.CommitmentId]
		if !ok {
			allocations[distribution.CommitmentId] = &util.Allocation{}
			totalAllocation = allocations[distribution.CommitmentId]
		}

		amount := distribution.Amount

		rocNeeded := rocOwed[distribution.CommitmentId].Minus(totalAllocation.ReturnOnCapital)
		if rocNeeded.GreaterThan(0) {
			allocation := balance.Min(rocNeeded, amount)
			amount = amount.Minus(allocation)
			totalAllocation.ReturnOnCapital = totalAllocation.ReturnOnCapital.Add(allocation)
		}

		prefNeeded := prefOwed[distribution.CommitmentId].Minus(totalAllocation.PreferredReturn)
		if prefNeeded.GreaterThan(0) {
			allocation := balance.Min(prefNeeded, amount)
			amount = amount.Minus(allocation)
			totalAllocation.PreferredReturn = totalAllocation.PreferredReturn.Add(allocation)
		}

		catchUpNeeded := prefOwed[distribution.CommitmentId].Times(catchUpRate).Minus(totalAllocation.CatchUp)
		if catchUpNeeded.GreaterThan(0) {
			allocation := balance.Min(catchUpNeeded, amount)
			amount = amount.Minus(allocation)
			totalAllocation.CatchUp = totalAllocation.CatchUp.Add(allocation)
		}

		finalSplitLPAllocation := amount.Times(lpFinalTierSplit)
		amount = amount.Minus(finalSplitLPAllocation)
		totalAllocation.FinalSplitLP = totalAllocation.FinalSplitLP.Add(finalSplitLPAllocation)

		totalAllocation.FinalSplitGP = totalAllocation.FinalSplitGP.Add(amount)
	}

	return allocations
}

func CalcRocAndPrefOwed(contributions []util.Transaction, runDate time.Time) (map[string]balance.Balance, map[string]balance.Balance) {
	rocOwed := make(map[string]balance.Balance)
	prefOwed := make(map[string]balance.Balance)

	for _, contribution := range contributions {
		rocOwed[contribution.CommitmentId] = rocOwed[contribution.CommitmentId].Add(contribution.Amount)

		days := int(runDate.Sub(contribution.Date).Hours() / 24)

		prefOwed[contribution.CommitmentId] = prefOwed[contribution.CommitmentId].Add(calculatePreferredReturn(contribution.Amount, days))
	}

	return rocOwed, prefOwed
}

func calculatePreferredReturn(amount balance.Balance, days int) balance.Balance {
	return amount.Times(math.Pow(1+hurdleRate, float64(days)/365)).Minus(amount)
}
