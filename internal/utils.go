package internal

import "slices"

func OnError(err error, callback func(err error)) {
	if err != nil {
		callback(err)
	}
}

func SubtractSlice(left []string, right []string) []string {
	remains := []string{}

	for _, key := range left {
		if !slices.Contains(right, key) {
			remains = append(remains, key)
		}
	}

	return remains
}
