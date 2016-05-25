package goipipnet

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strconv"
)

// Location 三级定位查询结果，无法确定的层级返回nil
type Location interface {
	Country() *Country
	Subdivision() *Subdivision
	City() *City
}

// Iso3166_1 ISO-3166-1国家代码
type Iso3166_1 struct {
	Alpha2  string
	Alpha3  string
	Numeric int
}

// Country 国家
type Country struct {
	ID           int       // 用户系统定义的城市ID，可为0
	Iso          Iso3166_1 // ISO-3166-1定义的国家ID
	Name         string    // 国家名称
	subdivisions map[string]*Subdivision
}

// Subdivision 子区域(省)
type Subdivision struct {
	ID      int    // 用户系统定义
	Iso     string // ISO-3166-2定义的子区域ID
	Name    string // 子区域名称
	country *Country
	cities  map[string]*City
}

// City 城市
type City struct {
	ID          int    // 用户系统定义的城市ID
	Name        string // 城市名称
	subdivision *Subdivision
}

// ISP 网络服务提供商
type ISP struct {
	ID   int    // ISP ID
	Name string // ISP名称
}

var (
	worldCountries map[string]*Country
	worldIsps      map[string]*ISP
)

// Country 返回国家本省
func (country *Country) Country() *Country {
	return country
}

// Subdivision 返回空
func (country *Country) Subdivision() *Subdivision {
	return nil
}

// City 返回空
func (country *Country) City() *City {
	return nil
}

func (country *Country) String() string {
	if country == nil {
		return "未知国家(0/-/-/0)"
	}
	return fmt.Sprintf("%s(%d/%s/%s/%d)", country.ID, country.Iso.Alpha2, country.Iso.Alpha3, country.Iso.Numeric)
}

// Country 返回子区域所在国家
func (subdivision *Subdivision) Country() *Country {
	return subdivision.country
}

// Subdivision 返回子区域本身
func (subdivision *Subdivision) Subdivision() *Subdivision {
	return subdivision
}

func (subdivision *Subdivision) String() string {
	if subdivision == nil {
		return "未知子区域(0/-)"
	}
	return fmt.Sprintf("%s(%d/%s)", subdivision.Name, subdivision.ID, subdivision.Iso)
}

// City 返回空
func (subdivision *Subdivision) City() *City {
	return nil
}

// Country 返回城市所在国家
func (city *City) Country() *Country {
	return city.subdivision.Country()
}

// Subdivision 返回城市所在子区域
func (city *City) Subdivision() *Subdivision {
	return city.subdivision
}

// City 返回城市本身
func (city *City) City() *City {
	return nil
}

func (city *City) String() string {
	if city == nil {
		return "未知城市(0)"
	}
	return fmt.Sprintf("%s(%d)", city.Name, city.ID)
}

func (isp *ISP) String() string {
	if isp == nil {
		return "未知运营商(0)"
	}
	return fmt.Sprintf("%s(%d)", isp.Name, isp.ID)
}

func loadWorld(data []byte) error {
	countries := make(map[string]*Country)
	isps := make(map[string]*ISP)
	reader := csv.NewReader(bytes.NewReader(data))
	var country *Country
	var subdivision *Subdivision
	for record, err := reader.Read(); err == nil; record, err = reader.Read() {
		switch record[0] {
		case "country":
			country = parseCountry(record[1:])
			if country == nil {
				return fmt.Errorf("错误的国家ID数据: %v", record)
			}
			countries[country.Name] = country
		case "subdivision":
			subdivision = parseSubdivision(record[1:], country)
			if subdivision == nil {
				return fmt.Errorf("错误的子区域ID数据: %v", record)
			}
			country.subdivisions[subdivision.Name] = subdivision
		case "city":
			city := parseCity(record[1:], subdivision)
			if city == nil {
				return fmt.Errorf("错误的城市ID数据: %v", record)
			}
			subdivision.cities[city.Name] = city
		case "isp":
			isp, aliases := parseISP(record[1:])
			if isp == nil {
				return fmt.Errorf("错误的ISP ID数据: %v", record)
			}
			isps[isp.Name] = isp
			for _, alias := range aliases {
				isps[alias] = isp
			}
		case "":
		default:
			return fmt.Errorf("未知数据行格式: %v", record)
		}
	}
	worldCountries = countries
	worldIsps = isps
	return nil
}

func loadWorldFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return loadWorld(content)
}

func parseCountry(fields []string) *Country {
	if len(fields) != 5 {
		return nil
	}
	id, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return nil
	}
	isoNumeric, err := strconv.ParseInt(fields[3], 10, 32)
	if err != nil {
		return nil
	}
	return &Country{
		ID: int(id),
		Iso: Iso3166_1{
			Alpha2:  fields[1],
			Alpha3:  fields[2],
			Numeric: int(isoNumeric),
		},
		Name:         fields[4],
		subdivisions: make(map[string]*Subdivision),
	}
}

func parseSubdivision(fields []string, country *Country) *Subdivision {
	if len(fields) != 3 || country == nil {
		return nil
	}
	id, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return nil
	}
	return &Subdivision{
		ID:      int(id),
		Iso:     fields[1],
		Name:    fields[2],
		country: country,
	}
}

func parseCity(fields []string, subdivision *Subdivision) *City {
	if len(fields) != 2 {
		return nil
	}
	id, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return nil
	}
	return &City{
		ID:          int(id),
		Name:        fields[1],
		subdivision: subdivision,
	}
}

func parseISP(fields []string) (*ISP, []string) {
	if len(fields) < 2 {
		return nil, nil
	}
	id, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return nil, nil
	}
	if len(fields) == 2 {
		return &ISP{
			ID:   int(id),
			Name: fields[1],
		}, nil
	}
	return &ISP{
		ID:   int(id),
		Name: fields[1],
	}, fields[2:]
}
