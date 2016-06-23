package ipipnet

import (
	"testing"

	"github.com/yangchenxing/go-ipipnet-downloader"
)

func TestLoadDatFail(t *testing.T) {
	index := &Index{
		Downloader: &downloader.Downloader{
			LocalPath: "no such path",
		},
	}
	if err := index.loadDat(); err == nil {
		t.Error("unexpected success")
		return
	}
}

func TestParseDatResult(t *testing.T) {
	index := &Index{}
	res := index.parseDatResult("中国\t中国\t\t\t")
	if res.ISPs != nil {
		t.Error("unexpected result:", res)
		return
	}
}
