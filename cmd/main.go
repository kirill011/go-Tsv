package main

import (
	"go-Tsv/internal/app"
	"go-Tsv/internal/config"
	"sync"
)

func main() {
	cfg := config.Init()

	app := app.New(cfg)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		app.Run()
	}()
	go func() {
		defer wg.Done()
		app.MonitFiles()
	}()

	wg.Wait()
}
