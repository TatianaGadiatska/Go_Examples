package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var (
	configModtime  int64
	errNotModified = errors.New("Not modified")
)

//CopyRule - хранит правила копирования для ScanDir
type CopyRule struct {
	Masks        []string `yaml:"masks"`
	DstDir       string   `yaml:"dst_dir"`
	IfExists     string   `yaml:"ifexists"`
	Mode         string   `yaml:"mode"`
	ExcludeMasks []string `yaml:"excludemasks"`
}

//ScanGroup - хранит настройки для папки, которая будет сканироваться
type ScanGroup struct {
	SrcDirs      []string         `yaml:"src_dirs"`
	Enabled      bool             `yaml:"enabled"`
	Rules        map[int]CopyRule `yaml:"rules"`
	CreateSrc    bool             `yaml:"create_src"`
	ExcludeMasks []string         `yaml:"excludemasks"`
}

// Config - структура для считывания конфигурационного файла
type Config struct {
	ScanGroups         []ScanGroup `yaml:"scangroups"`
	EnableHTTP         bool        `yaml:"enable_http"`
	Listen             string      `yaml:"listen"`
	MaxScanThreads     int         `yaml:"max_scan_threads"`
	MaxCopyThreads     int         `yaml:"max_copy_threads"`
	RescanInterval     int         `yaml:"rescaninterval"`
	LogLevel           string      `yaml:"loglevel"`
	GlobalExcludeMasks []string    `yaml:"excludemasks"`
}

func readConfig(ConfigName string) (x *Config, err error) {
	var file []byte
	if file, err = ioutil.ReadFile(ConfigName); err != nil {
		return nil, err
	}
	x = new(Config)
	if err = yaml.Unmarshal(file, x); err != nil {
		return nil, err
	}
	if x.LogLevel == "" {
		x.LogLevel = "Debug"
	}
	return x, nil
}

//Проверяет время изменения конфигурационного файла
//и перезагружает его если он изменился
//Возвращает errNotModified если изменений нет
func reloadConfig(configName string) (cfg *Config, err error) {
	info, err := os.Stat(configName)
	if err != nil {
		return nil, err
	}
	if configModtime != info.ModTime().UnixNano() {
		configModtime = info.ModTime().UnixNano()
		cfg, err = readConfig(configName)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return nil, errNotModified
}
