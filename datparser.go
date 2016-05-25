package goipipnet

import (
	"crypto/sha1"
	"encoding/binary"
	"io/ioutil"
	"strings"
)

func LoadDatFile(path string, unknownCallback func(string, []string)) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return LoadDat(content, unknownCallback)
}

func LoadDat(content []byte, unknownCallback func(string, []string)) error {
	textOffset := binary.BigEndian.Uint32(content[:4]) - 1024
	data := &ipData{
		sections: make([]ipSection, (textOffset-4-1024)/8),
		checksum: sha1.Sum(content),
	}
	for i, offset := 0, uint32(1028); offset < textOffset; i, offset = i+1, offset+8 {
		ip := binary.BigEndian.Uint32(content[offset : offset+4])
		data.sections[i].upperBound = ip
		dataRange := binary.BigEndian.Uint32(content[offset+4 : offset+8])
		dataOffset := dataRange & uint32(0x00FFFFFF)
		dataLength := dataRange >> 24
		data.sections[i].result = parseResult(string(content[dataOffset:dataOffset+dataLength]), unknownCallback)
		data.index[ip>>24] = i
	}
	index = data
	return nil
}

func parseResult(text string, unknownCallback func(string, []string)) Result {
	fields := strings.Split(text, "\t")
	var location Location
	if country := worldCountries[fields[0]]; country != nil {
		location = country
		if subdivision := country.subdivisions[fields[1]]; subdivision != nil {
			location = subdivision
			if city := subdivision.cities[fields[2]]; city != nil {
				location = city
			} else if fields[2] != "" && unknownCallback != nil {
				unknownCallback("city", fields)
			}
		} else if fields[1] != "" && fields[0] != fields[1] && unknownCallback != nil {
			unknownCallback("subdivision", fields[:2])
		}
	} else if fields[0] != "" && unknownCallback != nil {
		unknownCallback("country", fields[:1])
	}
	ispNames := strings.Split(fields[len(fields)-1], "/")
	isp := worldIsps[ispNames[0]]
	if isp == nil && ispNames[0] != "" && unknownCallback != nil {
		unknownCallback("isp", ispNames)
	}
	return Result{
		Location: location,
		ISP:      isp,
	}
}
