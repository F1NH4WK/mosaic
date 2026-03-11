package producer

import (
	"fmt"
	"strings"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if !encountered[elements[v]] && elements[v] != "" {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func GenerateCombinations(words []string, year int) []string {
	var variants []string

	for _, w := range words {
		lower := strings.ToLower(strings.TrimSpace(w))
		variants = append(variants, lower)
		variants = append(variants, capitalize(lower))
		variants = append(variants, strings.ToUpper(lower))
	}
	variants = removeDuplicates(variants)

	var combinations []string
	combinations = append(combinations, variants...)

	for i, w1 := range variants {
		for j, w2 := range variants {
			if i != j {
				if !strings.EqualFold(w1, w2) {
					combinations = append(combinations, w1+w2)
					combinations = append(combinations, w1+"_"+w2)
					combinations = append(combinations, w1+"."+w2)
				}
			}
		}
	}

	var withYears []string
	if year != 0 {
		yearStr := fmt.Sprintf("%d", year)
		shortYear := yearStr[len(yearStr)-2:]

		yearsToAppend := []string{yearStr, shortYear}

		for _, combo := range combinations {
			for _, y := range yearsToAppend {
				withYears = append(withYears, combo+y)
				withYears = append(withYears, combo+"_"+y)
			}
		}
	}
	combinations = append(combinations, withYears...)

	var finalCombinations []string
	finalCombinations = append(finalCombinations, combinations...)
	
	specialChars := []string{"!", "@", "123"}
	for _, combo := range combinations {
		for _, char := range specialChars {
			finalCombinations = append(finalCombinations, combo+char)
		}
	}

	return removeDuplicates(finalCombinations)
}