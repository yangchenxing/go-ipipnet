package ipipnet

import (
	"errors"
	"fmt"
	"github.com/yangchenxing/go-ip-index"
	"github.com/yangchenxing/go-ipipnet-downloader"
	"github.com/yangchenxing/go-regionid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Index struct {
	*downloader.Downloader
	Format               string
	MinBinarySearchRange int
	index                *ipindex.IPIndex
}

type Result struct {
	Location regionid.Location
	ISPs     []*regionid.ISP
}

func (result Result) Equal(other interface{}) bool {
	if result.Location != other || len(result.ISPs) != len(other.ISPs) {
		return false
	}
	for i, a := range result.ISPs {
		if a != other.ISPs[i] {
			return false
		}
	}
	return true
}

func (index *Index) Initialize() error {
	if index.Format != "DAT" {
		return fmt.Errorf("unsupported format: %s", index.Format)
	}
	if !regionid.Initialized() {
		if err := regionid.LoadBuiltinWorld(); err != nil {
			return fmt.Errorf("load built-in regionid fail: %s", err.Error())
		}
	}
	if index.MinBinarySearchRange <= 0 {
		index.MinBinarySearchRange = ipindex.DefaultMinBinarySearchRange
	}
	if err := index.EnsureLocal(); err != nil {
		return fmt.Errorf("ensure local file fail: %s", err.Error())
	}
	if err := index.load(); err != nil {
		return fmt.Errorf("load local file fail: %s", err.Error())
	}
	index.UpdateCallback = func(string) { index.load() }
	go index.StartWatch()
	return nil
}

func (index *Index) Search(ip net.IP) (result Result, err error) {
	var value ipindex.Value
	value, err = index.index.Search(ip)
	if err == nil && value != nil {
		result = value.(Result)
	}
	return
}

func (index *Index) load() error {
	switch index.Format {
	case "DAT":
		return index.loadDat()
	}
	return fmt.Errorf("unsupported format: %q", format)
}
