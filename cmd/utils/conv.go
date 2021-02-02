package utils

import "strconv"

// StringsToUints converts a slice of strings to a slice of uints.
func StringsToUints(strings []string) (uints []uint, err error) {
	for i := range strings {
		u, err := strconv.ParseUint(strings[i], 10, 0)
		if err != nil {
			return nil, err
		}

		uints = append(uints, uint(u))
	}

	return
}
