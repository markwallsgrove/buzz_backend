package models

import "strconv"

// StringToInt convert a string into a integer. If the
// conversation fails return the default value.
func StringToInt(str string, def int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return def
	}

	return i
}
