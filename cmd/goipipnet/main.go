package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/yangchenxing/goipipnet"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "goipipnet"
	app.Usage = "goipipnet实用工具"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "local",
			Value: "ipipnet.dat",
			Usage: "本地数据路径",
		},
		&cli.StringFlag{
			Name:  "remote",
			Value: "",
			Usage: "远程下载路径",
		},
		&cli.StringFlag{
			Name:  "format",
			Value: "",
			Usage: "文件格式",
		},
		&cli.DurationFlag{
			Name:  "interval",
			Value: time.Hour,
			Usage: "检查时间间隔",
		},
		&cli.StringFlag{
			Name:  "world",
			Value: "",
			Usage: "世界文件路径",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "locate",
			Usage:  "IP定位",
			Action: locate,
		},
	}
	app.Run(os.Args)
}

type ether struct{}

var (
	gopath string
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	gopath = file[:len(dir)-len("/github.com/yangchenxing/goipipnet/conf/goipipnet")]
}

func log(skip int, level string, format string, args ...interface{}) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = ""
		line = 0
	}
}

func locate(ctx *cli.Context) {
	unknowns := map[string]map[string]ether{
		"country":     make(map[string]ether),
		"subdivision": make(map[string]ether),
		"city":        make(map[string]ether),
		"isp":         make(map[string]ether),
	}
	var e ether
	config := goipipnet.Config{
		Local:         ctx.GlobalString("local"),
		Remote:        ctx.GlobalString("remote"),
		Format:        ctx.GlobalString("format"),
		World:         ctx.GlobalString("world"),
		CheckInterval: ctx.GlobalDuration("interval"),
		UnknownCallback: func(level string, data []string) {
			unknowns[level][strings.Join(data, "/")] = e
		},
	}
	goipipnet.Initialize(config)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimRight(scanner.Text(), "\n")
		ip := net.ParseIP(text)
		if ip == nil {
			fmt.Println("无效IP:", text)
			continue
		}
		result, err := goipipnet.Lookup(ip)
		if err != nil {
			fmt.Println("定位出错:", err.Error())
			continue
		}
		if result.Location == nil {
			fmt.Println("定位失败")
			continue
		}
		fmt.Println("国家:", result.Location.Country(), ", 省份:", result.Location.Subdivision(),
			", 城市:", result.Location.City(), ", 运营商:", result.ISP)
	}
}
