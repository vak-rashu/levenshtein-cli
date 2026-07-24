package levenshtein

import (
	"strings"
)

var db = map[string]int{
	"init":            1,
	"init-systemd(Ub": 2,
	"snapfuse":        162,
}

func calculateThreshold(dist int, s1, s2 string) bool {
	minLen := min(len(s1), len(s2))
	val := (dist / minLen) * 100
	if val < 45 {
		return true
	} else {
		return false
	}
}

func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	// initialize matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s2); i++ {
		matrix[i][0] = i
	}

	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min(
					matrix[i][j-1]+1,
					matrix[i-1][j]+1,
					matrix[i-1][j-1]+1,
				)
			}

		}
	}

	return matrix[len(s1)][len(s2)]
}

func hasPrefix(processName, arg string) string {
	has := strings.HasPrefix(processName, arg)
	if has {
		return processName
	} else {
		return ""
	}
}

var suggestions = make([]string, 0)

func Loop(arg string) ([]string, int) {

	for pname, pid := range db {
		if arg != pname {
			val := hasPrefix(pname, arg)
			if val != "" {
				return suggestions, pid
			} else {
				dist := levenshteinDistance(pname, arg)
				val := calculateThreshold(dist, pname, arg)
				if val {
					suggestions = append(suggestions, pname)
				} else {
					return suggestions, 0
				}
			}
		} else {
			return suggestions, pid
		}

	}

	return suggestions, 0
}
