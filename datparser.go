package ipipnet

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yangchenxing/go-ip-index"
	"github.com/yangchenxing/go-regionid"
)

func (index *Index) loadDat() error {
	builder := ipindex.NewIndexBuilder(index.MinBinarySearchRange)
	content, err := ioutil.ReadFile(index.LocalPath)
	if err != nil {
		return fmt.Errorf("read local file fail: %s", err.Error())
	}
	if len(content) < 1028 {
		return fmt.Errorf("bad content length: %d", len(content))
	}
	textOffset := binary.BigEndian.Uint32(content[:4]) - 1024
	lower := uint32(1)
	for i, offset := 0, uint32(1028); offset < textOffset; i, offset = i+1, offset+8 {
		upper := binary.BigEndian.Uint32(content[offset : offset+4])
		dataRange := binary.LittleEndian.Uint32(content[offset+4 : offset+8])
		dataOffset := textOffset + dataRange&uint32(0x00FFFFFF)
		dataLength := dataRange >> 24
		result := index.parseDatResult(string(content[dataOffset : dataOffset+dataLength]))
		builder.AddUint32(lower, upper, result)
		// 此处不检查错误添加错误情况，由ipip.net确保文件格式正确
		lower = upper + 1
	}
	index.index = builder.Build()
	return nil
}

func (index *Index) parseDatResult(text string) Result {
	fields := strings.Split(text, "\t")
	location := regionid.GetLocation(fields[0], fields[1], fields[2])
	ispNames := strings.Split(fields[len(fields)-1], "/")
	isps := make([]*regionid.ISP, 0, len(ispNames))
	for _, name := range ispNames {
		isp := regionid.GetISP(name)
		if isp == nil && index.KeepUnknownISP {
			isp = index.getUnknownISP(name)
		}
		if isp != nil {
			isps = append(isps, isp)
		}
	}
	if len(isps) == 0 {
		isps = nil
	}
	return Result{
		Location: location,
		ISPs:     isps,
	}
}
