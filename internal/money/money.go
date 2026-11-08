package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Amount int64

var ErrFormat = errors.New("amount must be a decimal with exactly two places, e.g. \"-12.50\"")

var ErrOverflow = errors.New("amount is too large")

func Parse(s string) (Amount, error) {
	neg := strings.HasPrefix(s, "-")
	body := strings.TrimPrefix(s, "-")

	dot := strings.IndexByte(body, '.')
	if dot < 1 || len(body)-dot-1 != 2 {
		return 0, ErrFormat
	}
	whole, frac := body[:dot], body[dot+1:]
	if !allDigits(whole) || !allDigits(frac) {
		return 0, ErrFormat
	}

	units, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, ErrOverflow
	}
	cents, _ := strconv.ParseInt(frac, 10, 64)
	if units > (math.MaxInt64-cents)/100 {
		return 0, ErrOverflow
	}
	total := units*100 + cents
	if neg {
		total = -total
	}
	return Amount(total), nil
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func (a Amount) Add(b Amount) (Amount, error) {
	if (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b) {
		return 0, ErrOverflow
	}
	return a + b, nil
}

func (a Amount) Cents() int64 { return int64(a) }

func (a Amount) String() string {

	u := uint64(a)
	sign := ""
	if a < 0 {
		sign = "-"
		u = -u
	}
	return fmt.Sprintf("%s%d.%02d", sign, u/100, u%100)
}

func (a Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || data[0] != '"' {
		return errors.New("amount must be a JSON string, not a number")
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := Parse(s)
	if err != nil {
		return err
	}
	*a = v
	return nil
}
