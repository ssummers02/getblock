package main

import "math/big"

// subHex subtracts the two hex strings a and b and returns the result as a hex string
func subHex(a, b string) string {
	// convert a and b to big.Int
	ai, _ := new(big.Int).SetString(a[2:], 16)
	bi, _ := new(big.Int).SetString(b[2:], 16)

	// subtract b from a
	res := new(big.Int).Sub(ai, bi)

	// convert the result back to hex string
	return "0x" + res.Text(16)
}

// addHex adds the two hex strings a and b and returns the result as a hex string
func addHex(a, b string) string {
	// convert a and b to big.Int
	ai, _ := new(big.Int).SetString(a[2:], 16)
	bi, _ := new(big.Int).SetString(b[2:], 16)

	// add a and b
	res := new(big.Int).Add(ai, bi)

	// convert the result back to hex string
	return "0x" + res.Text(16)
}

// calculateBalanceChanges calculates the balance changes for each address in each block
func calculateBalanceChanges(blocks []Block) map[string]*big.Int {
	balances := make(map[string]*big.Int)

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			// Update the balance of the 'from' address
			if tx.From != "" {
				fromBalance := balances[tx.From]
				if fromBalance == nil {
					fromBalance = big.NewInt(0)
				}
				value, ok := new(big.Int).SetString(tx.Value[2:], 16)
				if !ok {
					// handle invalid hex string
				}
				fromBalance = fromBalance.Sub(fromBalance, value)
				balances[tx.From] = fromBalance
			}
			// Update the balance of the 'to' address
			if tx.To != "" {
				toBalance := balances[tx.To]
				if toBalance == nil {
					toBalance = big.NewInt(0)
				}
				value, ok := new(big.Int).SetString(tx.Value[2:], 16)
				if !ok {
					// handle invalid hex string
				}
				toBalance = toBalance.Add(toBalance, value)
				balances[tx.To] = toBalance
			}
		}
	}

	return balances
}

// findMaxBalanceChange finds the address with the largest absolute balance change
func findMaxBalanceChange(balanceChanges map[string]*big.Int) (string, *big.Int) {
	maxAddress := ""
	maxChange := big.NewInt(0)

	for address, balance := range balanceChanges {
		change := new(big.Int).Abs(balance)
		if change.Cmp(maxChange) == 1 {
			maxChange = change
			maxAddress = address
		}
	}

	return maxAddress, maxChange
}
