package ipipnet

import (
	"fmt"
	"github.com/yangchenxing/go-ip-index"
	"github.com/yangchenxing/go-ipipnet-downloader"
	"github.com/yangchenxing/go-regionid"
	"net"
)

// Index provide ip to regions index ability. It supports data in "DAT" format only.
type Index struct {
	// The downloader to fetch remote data.
	*downloader.Downloader

	// Minimal range for binary search. It is used by ip index.
	MinBinarySearchRange int

	// Keep unknown isp with ID as 0
	KeepUnknownISP bool

	index       *ipindex.IPIndex
	unknownISPs map[string]*regionid.ISP
}

// Result is the search result
type Result struct {
	Location regionid.Location
	ISPs     []*regionid.ISP
}

// Equal compare two Result instance. It returns true if they are equal.
func (result Result) Equal(o interface{}) bool {
	other, ok := o.(Result)
	if !ok || result.Location != other.Location || len(result.ISPs) != len(other.ISPs) {
		return false
	}
	for i, a := range result.ISPs {
		if a != other.ISPs[i] {
			return false
		}
	}
	return true
}

// Initialize setup the Index instance.
func (index *Index) Initialize() error {
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

// Search returns the region and isp assiciated with the ip.
func (index *Index) Search(ip net.IP) (result Result, err error) {
	var value ipindex.Value
	value, err = index.index.Search(ip)
	if err == nil && value != nil {
		result = value.(Result)
	}
	return
}

func (index *Index) load() error {
	return index.loadDat()
}

func (index *Index) getUnknownISP(name string) *regionid.ISP {
	isp := index.unknownISPs[name]
	if isp == nil {
		isp = &regionid.ISP{
			ID:   0,
			Name: name,
		}
		index.unknownISPs[name] = isp
	}
	return isp
}
