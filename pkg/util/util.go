package util

import (
	"strconv"
)

func StringToInt64(str string) (int64, error) {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}
