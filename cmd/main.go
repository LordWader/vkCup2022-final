package main

import (
	"os"
	"vkCup2022-final/internal/http"
	"vkCup2022-final/internal/parser"
	"vkCup2022-final/internal/taskio"
)

// TODO - Добавить конкатенацию и разобраться как резать изображение в разных плоскостях

func main() {
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
	taskio.ReadFile()
	taskio.CheckFile()
}
