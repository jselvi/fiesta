package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSortInfo(t *testing.T) {
	// Data for testing
	j := make(map[int]string)
	j[00] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":100,\"Version\":\"TLS 1.2\"}"
	j[01] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[02] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[03] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[04] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[10] = "{\"Client\":\"127.0.0.1:50343\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":100,\"Version\":\"TLS 1.2\"}"
	j[11] = "{\"Client\":\"127.0.0.1:50343\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[12] = "{\"Client\":\"127.0.0.1:50343\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[13] = "{\"Client\":\"127.0.0.1:50343\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[14] = "{\"Client\":\"127.0.0.1:50343\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[20] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":100,\"Version\":\"TLS 1.2\"}"
	j[21] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[22] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[23] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[24] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[30] = "{\"Client\":\"127.0.0.1:50344\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":100,\"Version\":\"TLS 1.2\"}"
	j[31] = "{\"Client\":\"127.0.0.1:50344\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[32] = "{\"Client\":\"127.0.0.1:50344\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[33] = "{\"Client\":\"127.0.0.1:50344\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[34] = "{\"Client\":\"127.0.0.1:50344\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"

	// Prepare tests
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			"Already sorted",
			[]int{00, 01, 02, 03, 04, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34},
			[]int{00, 01, 02, 03, 04, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34},
		},
		{
			"Full overlap",
			[]int{00, 01, 02, 10, 11, 12, 13, 14, 03, 04, 20, 21, 22, 30, 31, 32, 33, 34, 23, 24},
			[]int{00, 01, 02, 03, 04, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34},
		},
		{
			"Two paralel connections",
			[]int{00, 10, 01, 11, 02, 12, 03, 13, 04, 14, 20, 30, 21, 31, 22, 32, 23, 33, 24, 34},
			[]int{00, 01, 02, 03, 04, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34},
		},
		{
			"Many paralel connections",
			[]int{00, 01, 10, 02, 11, 03, 04, 20, 30, 31, 21, 32, 22, 23, 12, 33, 34, 13, 14, 24},
			[]int{00, 01, 02, 03, 04, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34},
		},
	}

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input
			var input []string
			for _, i := range tt.input {
				input = append(input, j[i])
			}

			// Create want
			var want []string
			for _, i := range tt.want {
				want = append(want, j[i])
			}

			// Sort input
			var c Core
			got, e := c.sortInfo(input)
			if e != nil {
				t.Errorf("sortInfo error %v", e)
				return
			}

			// Compare
			if !reflect.DeepEqual(got, want) {
				t.Errorf("sortInfo got = %v, want = %v", got, want)
			}
		})
	}
}

func TestCompactAppData(t *testing.T) {
	// Data for testing
	j := make(map[int]string)
	j[00] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":100,\"Version\":\"TLS 1.2\"}"
	j[01] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[02] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[03] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[04] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"AppData\",\"Length\":10,\"Version\":\"TLS 1.2\"}"
	j[10] = "{\"Client\":\"127.0.0.1:50342\",\"IsRequest\":false,\"Type\":\"Raw\",\"Length\":10}"

	// Prepare tests
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			"No Raw in the middle",
			[]int{00, 01, 02, 03, 04, 10},
			[]int{100, 40, 10},
		},
		{
			"Raw in the middle",
			[]int{00, 01, 10, 02, 03, 10, 10, 04, 10},
			[]int{100, 40, 10, 10, 10, 10},
		},
		{
			"Two requests",
			[]int{00, 01, 02, 03, 04, 00, 01, 02},
			[]int{100, 40, 100, 20},
		},
		{
			"Many Raws",
			[]int{10, 00, 10, 01, 02, 10, 03, 04, 00, 10, 01, 02},
			[]int{10, 100, 10, 40, 10, 100, 10, 20},
		},
	}

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input
			var input []string
			for _, i := range tt.input {
				input = append(input, j[i])
			}

			// Compact input
			var c Core
			gotStr, e := c.compactAppData(input)
			if e != nil {
				t.Errorf("compactAppData error %v", e)
				return
			}

			// Calculate length
			var got []int
			for _, x := range gotStr {
				b := []byte(x)
				var m proxyInfo
				e = json.Unmarshal(b, &m)

				got = append(got, m.Length)
			}

			// Compare
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compactAppData got = %v, want = %v", got, tt.want)
			}
		})
	}
}
