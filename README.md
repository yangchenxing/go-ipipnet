# go-ipipnet

[![Go Report Card](https://goreportcard.com/badge/github.com/yangchenxing/go-ipipnet)](https://goreportcard.com/report/github.com/yangchenxing/go-ipipnet)
[![Build Status](https://travis-ci.org/yangchenxing/go-ipipnet.svg?branch=master)](https://travis-ci.org/yangchenxing/go-ipipnet)
[![GoDoc](http://godoc.org/github.com/yangchenxing/go-ipipnet?status.svg)](http://godoc.org/github.com/yangchenxing/go-ipipnet)
[![Coverage Status](https://coveralls.io/repos/github/yangchenxing/go-ipipnet/badge.svg?branch=master)](https://coveralls.io/github/yangchenxing/go-ipipnet?branch=master)

My golang library for IPIP.net.

## Example

    index := &Index{
      Downloader: &downloader.Downloader{
      LocalPath: "sample/mydata4vipweek2.dat",
      CheckETag: false,
    }
    
    result, _ := index.Search(net.ParseIP("58.32.100.100"))
    fmt.Println(result.Location.Name()) // 上海
    
## See

For more information about downloader in Index, see [go-ipipnet-downloader](http://github.com/yangchenxing/go-ipipnet-downloader).

For more information about ip index library used by Index, see [go-ip-index](http://github.com/yangchenxing/go-ip-index).
