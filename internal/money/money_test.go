package money

import (
	"encoding/json"
	"errors"
	"testing"
)

func mustParse(t *testing.T, s string) Amount {
	t.Helper()
	a, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q): %v", s, err)
	}
	return a
}

func TestTenPlusTwentyIsExactlyThirty(t *testing.T) {
	sum, err := mustParse(t, "0.10").Add(mustParse(t, "0.20"))
	if err != nil {
		t.Fatal(err)
	}
	if got := sum.String(); got != "0.30" {
		t.Errorf("0.10 + 0.20 = %q, want \"0.30\"", got)
	}
	out, _ := mustParse(t, "-0.10").Add(mustParse(t, "-0.20"))
	if got := out.String(); got != "-0.30" {
		t.Errorf("-0.10 + -0.20 = %q, want \"-0.30\"", got)
	}
}

func TestManySmallAmountsDoNotDrift(t *testing.T) {
	var total Amount
	dime := mustParse(t, "0.10")
	for i := 0; i < 1000; i++ {
		var err error
		if total, err = total.Add(dime); err != nil {
			t.Fatal(err)
		}
	}
	if got := total.String(); got != "100.00" {
		t.Errorf("1000 x 0.10 = %q, want \"100.00\"", got)
	}
}

func TestParseAndString(t *testing.T) {
	tests := []struct {
		in    string
		cents int64
		out   string
	}{
		{"0.00", 0, "0.00"},
		{"-0.00", 0, "0.00"},
		{"0.05", 5, "0.05"},
		{"-0.05", -5, "-0.05"},
		{"12.50", 1250, "12.50"},
		{"-12.50", -1250, "-12.50"},
		{"007.10", 710, "7.10"},
		{"1234567.89", 123456789, "1234567.89"},
	}
	for _, tc := range tests {
		a, err := Parse(tc.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.in, err)
			continue
		}
		if a.Cents() != tc.cents {
			t.Errorf("Parse(%q) = %d cents, want %d", tc.in, a.Cents(), tc.cents)
		}
		if a.String() != tc.out {
			t.Errorf("Parse(%q).String() = %q, want %q", tc.in, a.String(), tc.out)
		}
	}
}

func TestParseRejects(t *testing.T) {
	for _, in := range []string{
		"", "-", ".", "abc", "1", "1.5", "1.500", "-.50", "1.",
		"+1.00", " 1.00", "1.00 ", "1,50", "1e3", "1.2e1", "--1.00", "1.0a", "0x10.00",
	} {
		if _, err := Parse(in); !errors.Is(err, ErrFormat) {
			t.Errorf("Parse(%q) error = %v, want ErrFormat", in, err)
		}
	}
}

func TestParseOverflow(t *testing.T) {
	if _, err := Parse("92233720368547758.07"); err != nil {
		t.Errorf("largest amount should parse: %v", err)
	}
	for _, in := range []string{"92233720368547758.08", "99999999999999999999.00"} {
		if _, err := Parse(in); !errors.Is(err, ErrOverflow) {
			t.Errorf("Parse(%q) error = %v, want ErrOverflow", in, err)
		}
	}
}

func TestAddOverflow(t *testing.T) {
	max, min := Amount(1<<63-1), Amount(-1<<63)
	if _, err := max.Add(1); !errors.Is(err, ErrOverflow) {
		t.Errorf("max+1 error = %v, want ErrOverflow", err)
	}
	if _, err := min.Add(-1); !errors.Is(err, ErrOverflow) {
		t.Errorf("min-1 error = %v, want ErrOverflow", err)
	}
	if s := min.String(); s != "-92233720368547758.08" {
		t.Errorf("min.String() = %q", s)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type payload struct {
		Amount Amount `json:"amount"`
	}
	b, err := json.Marshal(payload{Amount: -1250})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"amount":"-12.50"}` {
		t.Errorf("marshal = %s", b)
	}
	var p payload
	if err := json.Unmarshal(b, &p); err != nil || p.Amount != -1250 {
		t.Errorf("unmarshal = %v, %v", p.Amount, err)
	}
}

func TestJSONRejectsNumbersAndBadStrings(t *testing.T) {
	type payload struct {
		Amount Amount `json:"amount"`
	}
	for _, in := range []string{
		`{"amount":-12.5}`, `{"amount":12}`, `{"amount":null}`, `{"amount":true}`,
		`{"amount":"12.5"}`, `{"amount":"abc"}`, `{"amount":""}`,
	} {
		var p payload
		if err := json.Unmarshal([]byte(in), &p); err == nil {
			t.Errorf("Unmarshal(%s) succeeded with %v, want an error", in, p.Amount)
		}
	}
}
