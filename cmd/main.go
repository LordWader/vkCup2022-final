package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"vkCup2022-final/internal/http"
	"vkCup2022-final/internal/parser"
	"vkCup2022-final/internal/taskio"
)

// TODO - Добавить конкатенацию и разобраться как резать изображение в разных плоскостях
func TrackMemoryUsage() {
	for {
		select {
		case <-time.Tick(time.Microsecond * 500):
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
			fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
			fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
			fmt.Printf("\tNumGC = %v\n", m.NumGC)
		}
	}
}

func main() {
	// start goroutime with memory tracker
	go TrackMemoryUsage()
	// store all files here
	_ = os.Mkdir("tmp", os.ModePerm)
	c := http.NewHttpClient("/")
	p := parser.NewParser()
	for i := 0; i < 10; i++ {
		go p.GetHrefPaths(c.NewHrefCh)
		go p.GetImagePaths(c.DownloadCh)
		go c.DownloadImage()
	}
	c.MakeBFSWalk(p.ImageCh, p.DownloadCh)
	taskio.CreateImage()
	taskio.CheckFile()
}
