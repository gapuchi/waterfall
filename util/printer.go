package util

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gapuchi/waterfall/v2/balance"
)

type Allocation struct {
	ReturnOnCapital balance.Balance
	PreferredReturn balance.Balance
	CatchUp         balance.Balance
	FinalSplitLP    balance.Balance
	FinalSplitGP    balance.Balance
}

func Print(fileName string, allocations map[string]*Allocation) error {

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		"Commitment ID",
		"Return On Capital",
		"Preferred Return",
		"Catch Up",
		"Final Split (LP)",
		"Final Split (GP)",
	})
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for commitmentId, allocation := range allocations {
		err = writer.Write([]string{
			commitmentId,
			allocation.ReturnOnCapital.ToString(),
			allocation.PreferredReturn.ToString(),
			allocation.CatchUp.ToString(),
			allocation.FinalSplitLP.ToString(),
			allocation.FinalSplitGP.ToString(),
		})
		if err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}
