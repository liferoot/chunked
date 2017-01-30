package chunked

import (
	"testing"
)

func TestPool_Nearest(t *testing.T) {
	pool := NewPool(128, 4096)
	cases := [...]struct {
		in  int
		exp int
	}{
		{-1, pool.Smallest()},
		{0, pool.Smallest()},
		{64, pool.Smallest()},
		{128, pool.Smallest()},
		{150, 256},
		{312, 512},
		{1024, 1024},
		{1700, 2048},
		{2100, pool.Largest()},
		{4096, pool.Largest()},
		{5000, pool.Largest()},
	}
	for _, c := range cases {
		out := pool.Nearest(c.in)
		if c.exp != out {
			t.Errorf("\n\tfor: %d\n\texp: %d\n\tgot: %d\n", c.in, c.exp, out)
		}
	}
}

func TestPool_Put(t *testing.T) {
	pool := NewPool(128, 4096)
	cases := [...]struct {
		in  []byte
		exp error
	}{
		{nil, ErrWrongSize},
		{make([]byte, 32), ErrWrongSize},
		{make([]byte, 64), ErrWrongSize},
		{make([]byte, 128), nil},
		{make([]byte, 200), ErrWrongSize},
		{make([]byte, 256), nil},
		{make([]byte, 312), ErrWrongSize},
		{make([]byte, 512), nil},
		{make([]byte, 1023), ErrWrongSize},
		{make([]byte, 1024), nil},
		{make([]byte, 1500), ErrWrongSize},
		{make([]byte, 2048), nil},
		{make([]byte, 4096), nil},
		{make([]byte, 5000), ErrWrongSize},
	}
	for _, c := range cases {
		out := pool.Put(c.in)
		if c.exp != out {
			t.Errorf("\n\tfor: [%d:%d]\n\texp: %v\n\tgot: %v\n", len(c.in), cap(c.in), c.exp, out)
		}
	}
}

func TestPool_Get(t *testing.T) {
	pool := NewPool(128, 4096)
	cases := [...]struct {
		in  int
		exp int
	}{
		{-1, pool.Smallest()},
		{0, pool.Smallest()},
		{64, pool.Smallest()},
		{128, pool.Smallest()},
		{150, 256},
		{312, 512},
		{1024, 1024},
		{1700, 2048},
		{2050, pool.Largest()},
		{5000, pool.Largest()},
	}
	for _, c := range cases {
		out := pool.Get(c.in)
		if c.exp != cap(out) || len(out) != 0 {
			t.Errorf("\n\tfor: %d\n\texp: %d/%d\n\tgot: %d/%d\n", c.in, 0, c.exp, len(out), cap(out))
		}
	}
}
