package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gapuchi/waterfall/v2/balance"
)

type Transaction struct {
	Date         time.Time
	Amount       balance.Balance
	CommitmentId string
}

type TransactionType string

const (
	Contribution TransactionType = "contribution"
	Distribution TransactionType = "distribution"
)

func ParseCommitments(fileLoc string) (map[string]balance.Balance, error) {
	commitmentToAmount := make(map[string]balance.Balance)

	csvFile, err := os.Open(fileLoc)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Read header
	if _, err = reader.Read(); err == io.EOF {
		return commitmentToAmount, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read from CSV: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read from CSV: %w", err)
		}

		commitmentId := record[1]

		amount, err := parseStringToAmount(record[3])
		if err != nil {
			return nil, err
		}

		totalAmount := commitmentToAmount[commitmentId]
		commitmentToAmount[commitmentId] = totalAmount.Add(amount)
	}

	return commitmentToAmount, nil
}

func ParseTransactions(fileLoc string) ([]Transaction, []Transaction, error) {
	contributions := make([]Transaction, 0)
	distributions := make([]Transaction, 0)

	csvFile, err := os.Open(fileLoc)
	if err != nil {
		return nil, nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Read header
	if _, err = reader.Read(); err == io.EOF {
		return contributions, nil, nil
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to read from CSV: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read from CSV: %w", err)
		}

		commitmentId := record[3]

		date, err := ParseDate(record[0])
		if err != nil {
			return nil, nil, err
		}

		amount, err := parseStringToAmount(record[1])
		if err != nil {
			return nil, nil, err
		}

		txn := Transaction{
			Date:         date,
			Amount:       amount,
			CommitmentId: commitmentId,
		}

		switch TransactionType(record[2]) {
		case Contribution:
			contributions = append(contributions, txn)
		case Distribution:
			distributions = append(distributions, txn)
		}
	}

	return contributions, distributions, nil
}

func parseStringToAmount(amountString string) (balance.Balance, error) {
	amountStringCleaned := strings.ReplaceAll(amountString, ",", "")
	amountStringCleaned = strings.ReplaceAll(amountStringCleaned, ".", "")
	amountStringCleaned = strings.ReplaceAll(amountStringCleaned, "$", "")
	amountStringCleaned = strings.ReplaceAll(amountStringCleaned, "(", "")
	amountStringCleaned = strings.ReplaceAll(amountStringCleaned, ")", "")
	amountStringCleaned = strings.TrimSpace(amountStringCleaned)

	amount, err := strconv.Atoi(amountStringCleaned)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("failed to parse amount: %w", err)
	}

	return balance.Balance{AmountInCents: int64(amount)}, nil
}

func ParseDate(dateString string) (time.Time, error) {
	date, err := time.Parse("01/02/2006", dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return date, nil
}
