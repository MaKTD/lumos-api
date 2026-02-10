package httpapi

import (
	"strconv"
	"strings"
)

func parsePrice(p string) (float64, error) {
	s := strings.TrimSpace(p)
	if s == "" {
		return 0, strconv.ErrSyntax
	}

	neg := false
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		neg = true
		s = strings.TrimSpace(s[1 : len(s)-1])
	}
	if strings.Contains(s, "-") {
		neg = true
	}

	// Keep only digits and common separators.
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || c == '.' || c == ',' {
			buf = append(buf, c)
		}
	}
	s = string(buf)
	if s == "" {
		return 0, strconv.ErrSyntax
	}

	lastDot := strings.LastIndexByte(s, '.')
	lastComma := strings.LastIndexByte(s, ',')

	switch {
	case lastDot >= 0 && lastComma >= 0:
		// Assume the last separator is the decimal; the other is thousands.
		if lastDot > lastComma {
			s = strings.ReplaceAll(s, ",", "")
		} else {
			s = strings.ReplaceAll(s, ".", "")
			idx := strings.LastIndexByte(s, ',')
			if idx >= 0 {
				pos := idx - strings.Count(s[:idx], ",")
				noCommas := strings.ReplaceAll(s, ",", "")
				if pos < 0 || pos > len(noCommas) {
					return 0, strconv.ErrSyntax
				}
				s = noCommas[:pos] + "." + noCommas[pos:]
			}
		}

	case lastComma >= 0:
		// Decide whether comma is decimal or thousands.
		after := len(s) - lastComma - 1
		if after > 0 && after <= 2 {
			idx := strings.LastIndexByte(s, ',')
			pos := idx - strings.Count(s[:idx], ",")
			noCommas := strings.ReplaceAll(s, ",", "")
			if pos < 0 || pos > len(noCommas) {
				return 0, strconv.ErrSyntax
			}
			s = noCommas[:pos] + "." + noCommas[pos:]
		} else {
			s = strings.ReplaceAll(s, ",", "")
		}

	case lastDot >= 0:
		// Heuristic: treat a single ".xxx" (3 digits) as thousands separator if it looks like grouping.
		after := len(s) - lastDot - 1
		if strings.Count(s, ".") == 1 && after == 3 && lastDot > 0 {
			s = strings.ReplaceAll(s, ".", "")
		} else if strings.Count(s, ".") > 1 {
			// Multiple dots: assume last is decimal if it has 1-2 digits, otherwise thousands.
			if after > 0 && after <= 2 {
				idx := strings.LastIndexByte(s, '.')
				pos := idx - strings.Count(s[:idx], ".")
				noDots := strings.ReplaceAll(s, ".", "")
				if pos < 0 || pos > len(noDots) {
					return 0, strconv.ErrSyntax
				}
				s = noDots[:pos] + "." + noDots[pos:]
			} else {
				s = strings.ReplaceAll(s, ".", "")
			}
		}
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if neg {
		v = -v
	}
	return v, nil
}
