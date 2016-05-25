package goipipnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var (
	errUninitialized     = errors.New("IP库未初始化")
	errNotIPv4           = errors.New("非有效IPv4")
	minBinarySearchRange = 10
)

type ipSection struct {
	upperBound uint32
	result     Result
}

type ipData struct {
	sections []ipSection
	index    [256]int
	checksum [20]byte
}

// Result 定位查询结果
type Result struct {
	Location Location
	ISP      *ISP
}

var (
	index       *ipData
	emptyResult Result
)

func (data *ipData) lookup(ip net.IP) (Result, error) {
	if data == nil {
		return emptyResult, errUninitialized
	}
	if ip = ip.To4(); ip == nil {
		return emptyResult, errNotIPv4
	}
	key := binary.BigEndian.Uint32([]byte(ip))
	indexKey := key >> 24
	if (indexKey >> 24) == 0 {
		return emptyResult, fmt.Errorf("无效IP: %s", ip)
	}
	lower, upper := data.index[indexKey-1]+1, data.index[indexKey]
	for upper-lower > minBinarySearchRange {
		mid := (lower + upper) / 2
		if key <= data.sections[mid-1].upperBound {
			upper = mid - 1
		} else if key > data.sections[mid].upperBound {
			lower = mid + 1
		} else {
			return data.sections[mid].result, nil
		}
	}
	for i := upper; i > lower; i++ {
		if key <= data.sections[i].upperBound {
			return data.sections[i].result, nil
		}
	}
	return data.sections[lower].result, nil
}

func Lookup(ip net.IP) (Result, error) {
	return index.lookup(ip)
}
