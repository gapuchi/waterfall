### Waterfall

This repo contains an implementation of a 4-tier waterfall. It accepts a commitments csv, transactions csv and date and will output a csv with the allocations for each tier, for each investor.

The transactions CSV should have the following columns in this order:

1. `transaction_date`
1. `transaction_amount`
1. `contribution_or_distribution`
1. `commitment_id`

The output CSV will have the following columns

1. Commitment ID
1. Return On Capital
1. Preferred Return
1. Catch Up
1. Final Split (LP)
1. Final Split (GP)

### Usage

This application is a CLI written in Go. To run, please install [Go 1.22](https://go.dev/dl/) on your machine.

To see its usage run:

```
go run .
```

Example command:

```
go run . -c commitments.csv -t transactions.csv -d 01/01/2025 -o output.csv
```

Example output CSV:

```
Commitment ID,Return On Capital,Preferred Return,Catch Up,Final Split (LP),Final Split (GP)
5,"250,000.00","125,514.99","31,378.75","170,485.01","42,621.25"
7,"350,000.00","175,721.01","43,930.25","232,278.99","58,069.75"
9,"476,666.67","241,594.83","60,398.71","249,071.83","62,267.96"
```

### Assumptions

- 8% hurdle, 100% catch-up, and 80/20 split.
- Fractional cents are round to the nearest cent after every operation.
- CSV is properly formatted and uses the `,` as the separator.
- Preferred Returns is calculated over 365 (leap years may need to be addressed.)

### Potential Improvements

- Avoid loading entire CSV into memory by operating on a row-by-row basis.
- Enable CLI to take in a single commitment ID and load only that specific ID.
- Implement Banker's Rounding
- Additional unit testing for thorough-ness
- Make the 4-tier values configurable