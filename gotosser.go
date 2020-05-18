package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	configFileName = "gotosser.yaml"
)

func main() {
	//загружаем конфиг

	cfg, err := reloadConfig(configFileName)
	if err != nil {
		if err != errNotModified {
			log.Fatalf("Не удалось загрузить %s: %s", configFileName, err)
		}
	}
	//log.Println("%#v", cfg)

	//запискаем цикл сканирования каталогов
	go scanLoop(cfg)

	//ожидаем завершения программы пл Ctrl - C

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	for {
		select {
		case <-sigChan:
			log.Println("CTRL+C: Завершаю работу.")
			return
		}
	}
}

func scanLoop(cfg *Config) {
	//периодически просматриваем конфиг и помечаем каталоги из него для сканирования файлов внутри
	for {
		for _, scandir := range cfg.ScanGroups {
			if scandir.Enabled != true {
				continue
			}
			for _, srcDir := range scandir.SrcDirs {
				log.Println("Сканируем каталог", srcDir)
			}
		}
		time.Sleep(time.Duration(cfg.RescanInterval) * time.Second)
	}
}
