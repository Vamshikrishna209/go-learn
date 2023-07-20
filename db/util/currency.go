package util

import "sort"

var CurrencyList = []string{
	"USD",
	"INR",
	"EUR",
	"CAD",
}

func IsSupportedCurrency(currency string) bool {
	sort.Strings(CurrencyList)
	res := sort.SearchStrings(CurrencyList, currency)
	return res < len(CurrencyList) && currency == CurrencyList[res]
}
