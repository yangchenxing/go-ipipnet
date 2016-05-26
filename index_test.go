package ipipnet

import (
	"fmt"
	"github.com/yangchenxing/go-ipipnet-downloader"
	"net"
	"os"
	"testing"
)

var (
	idx *Index
)

func TestSearch(t *testing.T) {
	if result, err := idx.Search(net.ParseIP("58.32.100.100")); err != nil {
		t.Error("search 58.32.100.100 fail:", err.Error())
		return
	} else if result.Location.Subdivision().Name() != "上海" {
		t.Error("unexpected location of 58.32.100.100:", result.Location)
		return
	}
	if result, err := idx.Search(net.ParseIP("58.30.100.100")); err != nil {
		t.Error("search 58.30.100.100 fail:", err.Error())
		return
	} else if result.Location.Subdivision().Name() != "北京" {
		t.Error("unexpected location of 58.30.100.100:", result.Location)
		return
	}
}

func BenchmarkSearch(b *testing.B) {
	testIP := []net.IP{
		net.ParseIP("58.32.100.100"),
		net.ParseIP("58.30.100.100"),
	}
	for i := 0; i < b.N; i++ {
		idx.Search(testIP[i%len(testIP)])
	}
}

func TestMain(m *testing.M) {
	idx = &Index{
		Downloader: &downloader.Downloader{
			LocalPath: "sample/mydata4vipweek2.dat",
			CheckETag: false,
		},
	}
	if err := idx.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "initialize index fail: %s\n", err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}
