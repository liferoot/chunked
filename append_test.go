package chunked

import "testing"
import "bytes"

func TestAppend(t *testing.T) {
	cases := [...]struct {
		in  [][]byte
		exp []byte
	}{
		{[][]byte{nil}, nil},
		{[][]byte{nil, nil}, nil},
		{[][]byte{nil, []byte{'1', '2', '3'}}, []byte{'1', '2', '3'}},
		{[][]byte{nil, []byte{'1'}, []byte{'2'}, []byte{'3'}}, []byte{'1', '2', '3'}},
		{[][]byte{[]byte{}, []byte{'1', '2', '3'}}, []byte{'1', '2', '3'}},
		{[][]byte{[]byte{}, []byte{'1'}, []byte{'2'}, []byte{'3'}}, []byte{'1', '2', '3'}},
		{[][]byte{[]byte{'1', '2', '3'}}, []byte{'1', '2', '3'}},
		{[][]byte{[]byte{'1', '2', '3'}, nil}, []byte{'1', '2', '3'}},
		{[][]byte{[]byte{'1', '2', '3'}, []byte{'4'}, []byte{'5'}}, []byte{'1', '2', '3', '4', '5'}},
		{[][]byte{[]byte{'1', '2', '3'}, []byte{'4', '5'}}, []byte{'1', '2', '3', '4', '5'}},
		{[][]byte{makeSlice(16, []byte{'1', '2', '3'}), []byte{'4'}, []byte{'5'}}, []byte{'1', '2', '3', '4', '5'}},
		{[][]byte{makeSlice(16, []byte{'1', '2', '3'}), []byte{'4', '5'}}, []byte{'1', '2', '3', '4', '5'}},
	}
	var out []byte

	for _, c := range cases {
		if len(c.in) > 1 {
			out = Append(c.in[0], c.in[1:]...)
		} else {
			out = Append(c.in[0])
		}
		if !bytes.Equal(c.exp, out) {
			t.Errorf("\n\tfor: %v, len/cap: %d/%d\n\texp: %v\n\tgot: %v, len/cap: %d/%d\n",
				c.in, len(c.in[0]), cap(c.in[0]), c.exp, out, len(out), cap(out))
		}
	}
}

func TestAppendByte(t *testing.T) {
	cases := [...]struct {
		in  [2][]byte
		exp []byte
	}{
		{[2][]byte{nil, nil}, nil},
		{[2][]byte{nil, []byte{'1', '2', '3'}}, []byte{'1', '2', '3'}},
		{[2][]byte{[]byte{}, []byte{'1', '2', '3'}}, []byte{'1', '2', '3'}},
		{[2][]byte{[]byte{'1', '2', '3'}, nil}, []byte{'1', '2', '3'}},
		{[2][]byte{[]byte{'1', '2', '3'}, []byte{'4', '5'}}, []byte{'1', '2', '3', '4', '5'}},
		{[2][]byte{makeSlice(16, []byte{'1', '2', '3'}), []byte{'4', '5'}}, []byte{'1', '2', '3', '4', '5'}},
	}
	for _, c := range cases {
		out := AppendByte(c.in[0], c.in[1]...)
		if !bytes.Equal(c.exp, out) {
			t.Errorf("\n\tfor: %v, len/cap: %d/%d\n\texp: %v\n\tgot: %v, len/cap: %d/%d\n",
				c.in, len(c.in[0]), cap(c.in[0]), c.exp, out, len(out), cap(out))
		}
	}
}

func makeSlice(capacity int, payload []byte) []byte {
	l := len(payload)
	if capacity < l {
		capacity = l
	}
	b := make([]byte, l, capacity)
	copy(b, payload)
	return b
}
