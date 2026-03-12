// internal/producer/generator.go
package producer

import (
	"strings"
	"unicode"

	"github.com/F1NH4WK/mosaic/internal/models"
)

type Rules struct {
	MinLength    int
	RequireUpper bool
	RequireLower bool
	RequireNum   bool
	RequireSpec  bool
	UseLeetspeak bool
}

func GenerateCombinations(profile models.Profile) []string {
	var results []string
	wordSet := make(map[string]struct{})

	add := func(s string) {
		if s != "" {
			wordSet[s] = struct{}{}
		}
	}

	var bases []string
	allWords := append(profile.Names, profile.Keywords...)
	for _, w := range allWords {
		if w == "" {
			continue
		}
		bases = append(bases, strings.ToLower(w))
		bases = append(bases, strings.Title(strings.ToLower(w)))
		bases = append(bases, strings.ToUpper(w))
	}

	var years, months, days []string
	if len(profile.DOB) == 8 {
		d := profile.DOB[0:2]
		m := profile.DOB[2:4]
		y := profile.DOB[4:8]
		shortY := profile.DOB[6:8]

		days = []string{d}
		months = []string{m}
		years = []string{y, shortY}

	
		bds := []string{
			d + m + y,      
			m + d + y,      
			d + m + shortY, 
			m + d + shortY, 
			d + m,          
			m + d,          
			d + y,          
			d + shortY,     
			m + y,          
			m + shortY,     
			y + m + d,      
		}
		bases = append(bases, bds...) 
	}


	for _, w := range bases {
		add(w)
	}


	commonSuffixes := []string{"123", "1234", "12345", "123456", "1", "12", "!", "@", "!!", "123!", "2023", "2024"}
	commonSuffixes = append(commonSuffixes, years...)
	commonSuffixes = append(commonSuffixes, months...)
	commonSuffixes = append(commonSuffixes, days...)

	for _, w1 := range bases {
		
	
		for _, suf := range commonSuffixes {
			add(w1 + suf)
			add(w1 + "_" + suf)
		}

		for _, w2 := range bases {
			if strings.EqualFold(w1, w2) {
				continue 
			}
			
		
			combos := []string{w1 + w2, w1 + "_" + w2, w1 + "." + w2}
			for _, c := range combos {
				add(c)
				
			
				for _, suf := range commonSuffixes {
					add(c + suf)
				}
			}
		}
	}

	for k := range wordSet {
		results = append(results, k)
	}

	return results
}

func GeneratePasswords(baseWord string, outChan chan<- string, rules Rules) {
	buffer := []byte(baseWord)
	
	if len(buffer) < rules.MinLength {
		return
	}
	
	if !rules.UseLeetspeak {
		if isValid(buffer, rules) {
			outChan <- string(buffer)
		}
		return
	}

	backtrackLeetspeak(buffer, 0, outChan, rules)
}

func backtrackLeetspeak(buffer []byte, index int, outChan chan<- string, rules Rules) {
	if index == len(buffer) {
		if isValid(buffer, rules) {
			outChan <- string(buffer)
		}
		return
	}

	originalChar := buffer[index]
	leetOptions := getLeetChars(originalChar)

	for _, mut := range leetOptions {
		buffer[index] = mut
		backtrackLeetspeak(buffer, index+1, outChan, rules)
	}
	
	buffer[index] = originalChar
}

func getLeetChars(c byte) []byte {
	switch c {
	case 'a', 'A': return []byte{c, '@', '4'}
	case 'e', 'E': return []byte{c, '3'}
	case 'i', 'I': return []byte{c, '1', '!'}
	case 'o', 'O': return []byte{c, '0'}
	case 's', 'S': return []byte{c, '$', '5'}
	case 't', 'T': return []byte{c, '7'}
	default:       return []byte{c}
	}
}

func isValid(pass []byte, rules Rules) bool {
	if len(pass) < rules.MinLength { return false }
	
	hasUpper, hasLower, hasNum, hasSpec := false, false, false, false
	for _, b := range pass {
		r := rune(b)
		if unicode.IsUpper(r) { hasUpper = true }
		if unicode.IsLower(r) { hasLower = true }
		if unicode.IsDigit(r) { hasNum = true }
		if unicode.IsPunct(r) || unicode.IsSymbol(r) { hasSpec = true }
	}

	if rules.RequireUpper && !hasUpper { return false }
	if rules.RequireLower && !hasLower { return false }
	if rules.RequireNum && !hasNum { return false }
	if rules.RequireSpec && !hasSpec { return false }

	return true
}