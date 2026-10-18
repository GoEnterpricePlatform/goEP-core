package helpers

import (
	"strconv"
)

func ParsePaddleAmount(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	amount, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
