package lexer

import (
	"testing"
)

func testStringEqual(t *testing.T, want string, got string, msg string) {
	t.Helper()
	if want != got {
		t.Fatalf("want=%q, but got=%q, message: %s", want, got, msg)
	}
}

func TestPy3ReadInteger(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"0", "0"},
		{"0  ", "0"},
		{"0_  ", ""},
		{"0_0  ", "0_0"},
		{"00000", "00000"},
		{"7", "7"},
		{"2147483647", "2147483647"},
		{"100_000_000_000", "100_000_000_000"},
		{"0o177", "0o177"},
		{"0b100110111", "0b100110111"},
		{"0b_1110_0101", "0b_1110_0101"},
		{"3", "3"},
		{"79228162514264337593543950336", "79228162514264337593543950336"},
		{"0o377", "0o377"},
		{"0o_377_", ""},
		{"0xdeadbeef", "0xdeadbeef"},
		{"0x_deadbeef", "0x_deadbeef"},
		{"0xdeadbeefG", "0xdeadbeef"},
		{"0Xdeadbeef", "0Xdeadbeef"},
		{"0Gdeadbeef", "0"},
		{"01", ""},
	}
	for _, test := range tests {
		buffer := NewCodeBuffer(test.code)
		got := Py3ReadInteger(buffer)
		testStringEqual(t, test.want, got, test.code)
	}
}

func TestPy3ReadFloat(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{".   ", ""},
		{".0", ".0"},
		{".0   ", ".0"},
		{".001", ".001"},
		{".001e+0", ".001e+0"},
		{".001e+0   ", ".001e+0"},
		{".001e+123   ", ".001e+123"},
		{".001E+123   ", ".001E+123"},
		{".001E-123   ", ".001E-123"},
		{".001E123   ", ".001E123"},
		{".001ED123   ", ""},

		{"314", ""},
		{"3.14", "3.14"},
		{"10.", "10."},
		{"1e100", "1e100"},
		{"3.14e-10", "3.14e-10"},
		{"3.14e+10", "3.14e+10"},
		{"0e0 ", "0e0"},
	}
	for _, test := range tests {
		buffer := NewCodeBuffer(test.code)
		got := Py3ReadFloat(buffer)
		testStringEqual(t, test.want, got, test.code)
	}
}

func TestPy3ReadNumber(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"0", "0"},
		{"0  ", "0"},
		{"0_  ", ""},
		{"0_0  ", "0_0"},
		{"00000", "00000"},
		{"7", "7"},
		{"2147483647", "2147483647"},
		{"100_000_000_000", "100_000_000_000"},
		{"0o177", "0o177"},
		{"0b100110111", "0b100110111"},
		{"0b_1110_0101", "0b_1110_0101"},
		{"3", "3"},
		{"79228162514264337593543950336", "79228162514264337593543950336"},
		{"0o377", "0o377"},
		{"0o_377_", ""},
		{"0xdeadbeef", "0xdeadbeef"},
		{"0x_deadbeef", "0x_deadbeef"},
		{"0xdeadbeefG", "0xdeadbeef"},
		{"0Xdeadbeef", "0Xdeadbeef"},
		{"0Gdeadbeef", "0"},
		{"01", ""},

		{".   ", ""},
		{".0", ".0"},
		{".0   ", ".0"},
		{".001", ".001"},
		{".001e+0", ".001e+0"},
		{".001e+0   ", ".001e+0"},
		{".001e+123   ", ".001e+123"},
		{".001E+123   ", ".001E+123"},
		{".001E-123   ", ".001E-123"},
		{".001E123   ", ".001E123"},
		{".001ED123   ", ""},

		{"3.14", "3.14"},
		{"10.", "10."},
		{"1e100", "1e100"},
		{"3.14e-10", "3.14e-10"},
		{"3.14e+10", "3.14e+10"},
		{"0e0 ", "0e0"},
	}
	for _, test := range tests {
		buffer := NewCodeBuffer(test.code)
		got := Py3ReadNumber(buffer)
		testStringEqual(t, test.want, got, test.code)
	}
}
