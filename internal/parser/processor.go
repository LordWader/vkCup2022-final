package parser

import (
	"log"
	"regexp"
)

type Parser struct {
	ImageCh    chan string
	DownloadCh chan string
}

func NewParser() *Parser {
	return &Parser{
		ImageCh:    make(chan string, 10),
		DownloadCh: make(chan string, 10),
	}
}

func (p *Parser) GetImagePaths(out chan<- []string) {
	for content := range p.ImageCh {
		imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
		subMatchSlice := imageRegExp.FindAllStringSubmatch(content, -1)
		outBuf := make([]string, len(subMatchSlice))
		for ind, item := range subMatchSlice {
			outBuf[ind] = item[1]
			log.Println("Image found : ", item[1])
		}
		out <- outBuf
	}
}

func (p *Parser) GetHrefPaths(out chan<- []string) {
	for content := range p.DownloadCh {
		hrefRegExp := regexp.MustCompile(`<a[^>]+\bhref=["']([^"']+)["']`)
		subMatchSlice := hrefRegExp.FindAllStringSubmatch(content, -1)
		outBuf := make([]string, len(subMatchSlice))
		for ind, item := range subMatchSlice {
			outBuf[ind] = item[1]
			log.Println("Href found : ", item[1])
		}
		out <- outBuf
	}
}
