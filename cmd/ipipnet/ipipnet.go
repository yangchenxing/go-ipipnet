package main

import (
	"flag"
	"github.com/yangchenxing/go-ipipnet"
)

var (
	localPath = flag.String("local", "ipipnet.dat", "本地数据文件路径")
	remoteURL = flag.String("remote", "", "远程下载地址")
)
