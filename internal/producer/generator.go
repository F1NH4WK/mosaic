package producer

import "unicode"

type Rules struct {
	MinLength int
	RequireUpper bool
	RequireLower bool
	RequireNum   bool
	RequireSpec  bool
	UseLeetspeak bool
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
		// O(N)
		if isValid(buffer, rules) {
			outChan <- string(buffer)
		}
		return
	}

	originalChar := buffer[index]
	substitutions := getLeetChars(originalChar)

	for _, sub := range substitutions {
		buffer[index] = sub
		backtrackLeetspeak(buffer, index+1, outChan, rules)
	}

	buffer[index] = originalChar
}


func isValid(pwd []byte, rules Rules) bool {
	if !rules.RequireUpper && !rules.RequireLower && !rules.RequireNum && !rules.RequireSpec {
		return true 
	}

	var hasU, hasL, hasN, hasS bool
	for _, b := range pwd {
		if b >= 'A' && b <= 'Z' { hasU = true; continue }
		if b >= 'a' && b <= 'z' { hasL = true; continue }
		if b >= '0' && b <= '9' { hasN = true; continue }
		hasS = true
	}

	return (!rules.RequireUpper || hasU) &&
		(!rules.RequireLower || hasL) &&
		(!rules.RequireNum || hasN) &&
		(!rules.RequireSpec || hasS)
}

func getLeetChars(b byte) []byte {
	lower := byte(unicode.ToLower(rune(b)))

	switch lower {
	case 'a':
		return []byte{'a', 'A', '@', '4'}
	case 'e':
		return []byte{'e', 'E', '3'}
	case 'i':
		return []byte{'i', 'I', '1', '!'}
	case 'o':
		return []byte{'o', 'O', '0'}
	case 's':
		return []byte{'s', 'S', '$', '5'}
	case 't':
		return []byte{'t', 'T', '7'}
	default:
		if unicode.IsLetter(rune(b)) {
			return []byte{lower, byte(unicode.ToUpper(rune(b)))}
		}
		return []byte{b}
	}
}