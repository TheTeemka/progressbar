package progressbar_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/TheTeemka/progressbar"
	"github.com/TheTeemka/progressbar/bar"
)

type Slower struct{}

func (s Slower) Write(p []byte) (int, error) {
	time.Sleep(time.Millisecond * 1)
	return len(p), nil
}

func TestWithTotal(t *testing.T) {
	input, err := os.Open("ex.txt")
	if err != nil {
		t.Fatal(err)
	}
	st, err := input.Stat()
	if err != nil {
		t.Fatal(err)
	}

	bar, err := progressbar.New(st.Size(), bar.NewCarBar())
	if err != nil {
		t.Fatal(err)
	}
	os, err := os.Create("file.txt")
	if err != nil {
		t.Fatal(err)
	}

	var slower Slower
	io.Copy(io.MultiWriter(os, bar, slower), input)
}

// func TestWithoutTotal(t *testing.T) {
// 	input, err := os.Open("ex.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	bar, err := progressbar.New(-1, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	os, err := os.Create("file.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	io.Copy(io.MultiWriter(os, bar), input)
// }
