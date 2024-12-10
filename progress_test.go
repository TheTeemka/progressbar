package progressbar_test

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/TheTeemka/progressbar"
)

func TestBar(t *testing.T) {
	p := progressbar.ProgressBar{
		TotalBytes: 1 << 20,
		DownBytes:  rand.Int63n(1 << 20),
	}
	p.Print()
}

func TestMal(t *testing.T) {
	url := "http://212.183.159.230/200MB.zip"

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()

	bar := progressbar.New(resp.ContentLength)
	os, err := os.Create("file.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(io.MultiWriter(os, bar), resp.Body)
}
