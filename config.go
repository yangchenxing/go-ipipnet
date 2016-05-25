package goipipnet

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Local           string
	Remote          string
	Format          string
	World           string
	CheckInterval   time.Duration
	UnknownCallback func(string, []string)
}

var (
	errNoETag     = errors.New("HTTP应答缺少ETag头")
	errUnmodified = errors.New("远程文件未发生修改")

	defaultCheckInterval = time.Hour
)

func Initialize(config Config) error {
	if config.World == "" {
		if err := loadWorld(builtinWorld); err != nil {
			return err
		}
		logDebug("加载IP库內建世界完成")
	} else if err := loadWorldFile(config.World); err != nil {
		return err
	} else {
		logDebug("加载IP库外部世界完成")
	}
	if err := guardLocal(config); err != nil {
		return err
	}
	if err := load(config); err != nil {
		return err
	}
	logDebug("加载IP库完成")
	go watch(config)
	return nil
}

func load(config Config) error {
	format := strings.ToUpper(config.Format)
	if format == "" {
		format = strings.ToUpper(filepath.Ext(config.Local))
	}
	switch format {
	case "DAT":
		if err := LoadDatFile(config.Local, config.UnknownCallback); err != nil {
			return err
		}
	default:
		return fmt.Errorf("不支持的文件格式: %q", format)
	}
	return nil
}

func guardLocal(config Config) error {
	if _, err := os.Stat(config.Local); err == nil {
		return nil
	}
	return download(config.Local, config.Remote)
}

func download(local, remote string) error {
	response, err := http.Get(remote)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	etag := response.Header.Get("ETag")
	if etag == "" {
		return errNoETag
	} else if !strings.HasPrefix(etag, "sha1-") {
		return fmt.Errorf("不支持的ETag: %q", etag)
	}
	if index != nil && fmt.Sprintf("sha1-%x", index.checksum) == etag {
		return nil
	}
	localTemp := local + ".tmp"
	if file, err := os.OpenFile(localTemp, os.O_CREATE|os.O_WRONLY, 0755); err != nil {
		return fmt.Errorf("新建临时文件出错: path=%q, error=%q", localTemp, err.Error())
	} else if _, err := io.Copy(file, response.Body); err != nil {
		file.Close()
		return fmt.Errorf("保存临时文件出错: path=%q, error=%q", localTemp, err.Error())
	} else {
		file.Close()
	}
	if err := os.Rename(localTemp, local); err != nil {
		return fmt.Errorf("重命名临时文件出错: old=%q, new=%q", localTemp, local)
	}
	logDebug("下载新文件完成: etag=%q", etag)
	return nil
}

func watch(config Config) {
	if config.Remote == "" {
		watchLocal(config)
	} else {
		watchRemote(config)
	}
}

func watchLocal(config Config) {
	var timestamp time.Time
	if info, err := os.Stat(config.Local); err == nil {
		timestamp = info.ModTime()
	} else {
		logError("本地IP数据库不存在: %s", err.Error())
	}
	interval := config.CheckInterval
	if interval == 0 {
		interval = defaultCheckInterval
	}
	for {
		time.Sleep(config.CheckInterval)
		if info, err := os.Stat(config.Local); err != nil {
			logError("本地IP数据库不存在: %s", err.Error())
			continue
		} else if info.ModTime().After(timestamp) {
			if err := load(config); err != nil {
				logError("加载IP库更新出错: %s", err.Error())
			} else {
				logDebug("加载IP库本地更新完成")
			}
		}
	}
}

func watchRemote(config Config) {
	for {
		if err := download(config.Local, config.Remote); err == nil {
			if err := load(config); err != nil {
				logError("加载IP库更新出错: %s", err.Error())
			} else {
				logDebug("加载IP库远程更新完成")
			}
		} else if err == errUnmodified {
			logDebug("IP库无更新")
		} else {
			logError("下载IP库出错: %s", err.Error())
		}
		time.Sleep(config.CheckInterval)
	}
}
