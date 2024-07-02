package main

import (
	"flag"
	"log"
	"os"

	"github.com/gapuchi/waterfall/v2/util"
)

const (
	hurdleRate       = 0.08
	catchUpRate      = 0.25
	lpFinalTierSplit = 0.8
)

func main() {

	var commitmentsFileLoc, transactionsFileLoc, outputFileLoc, dateString string
	flag.StringVar(&commitmentsFileLoc, "c", "", "location of the commitment file")
	flag.StringVar(&transactionsFileLoc, "t", "", "location of the transaction file")
	flag.StringVar(&dateString, "d", "", "date to run the waterfall")
	flag.StringVar(&outputFileLoc, "o", "", "name of output file")
	flag.Parse()

	if commitmentsFileLoc == "" || transactionsFileLoc == "" || dateString == "" || outputFileLoc == "" {
		flag.Usage()
		return
	}

	runDate, err := util.ParseDate(dateString)
	if err != nil {
		log.Fatalf("Date [%s] is incorrectly formatted. It must be in the form of mm/dd/yyyy", dateString)
	}

	_, err = util.ParseCommitments(commitmentsFileLoc)
	if os.IsNotExist(err) {
		log.Fatalf("File %s not found", commitmentsFileLoc)
	} else if err != nil {
		log.Fatal(err)
	}

	contributions, distributions, err := util.ParseTransactions(transactionsFileLoc)
	if os.IsNotExist(err) {
		log.Fatalf("File %s not found", transactionsFileLoc)
	} else if err != nil {
		log.Fatal(err)
	}

	rocOwed, prefOwed := CalcRocAndPrefOwed(contributions, runDate)
	allocations := Allocate(distributions, rocOwed, prefOwed)

	if err = util.Print(outputFileLoc, allocations); err != nil {
		log.Fatalf("failed to write to csv: %s", err.Error())
	}
}
