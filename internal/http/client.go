package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var htmlPool = sync.Pool{New: func() interface{} {
	b := make([]byte, 0, 4096)
	return b
},
}

type Client struct {
	queue      []string
	used       sync.Map
	DownloadCh chan []string
	NewHrefCh  chan []string
}

func NewHttpClient(rootPath string) *Client {
	return &Client{queue: []string{rootPath},
		used:       sync.Map{},
		DownloadCh: make(chan []string, 10),
		NewHrefCh:  make(chan []string, 10),
	}
}

func (c *Client) MakeBFSWalk(imageCh, downloadCh chan<- string) {
	defer func() {
		close(imageCh)
		close(downloadCh)
	}()
	for len(c.queue) > 0 {
		next := c.queue[0]
		if _, ok := c.used.Load(next); ok {
			c.queue = c.queue[1:]
			continue
		}
		c.used.Store(next, true)
		data, _ := c.GetHtmlContent(next)
		imageCh <- data
		downloadCh <- data
		c.queue = c.queue[1:]
		newHref := <-c.NewHrefCh
		for _, nh := range newHref {
			c.queue = append(c.queue, nh)
		}
	}
}

func (c *Client) GetHtmlContent(prefix string) (string, error) {
	url := fmt.Sprintf("http://localhost:8080%s", prefix)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("can't get html content of page")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return string(data), nil
}

func (c *Client) DownloadImage() {
	for batch := range c.DownloadCh {
		for _, prefix := range batch {
			url := fmt.Sprintf("http://localhost:8080%s", prefix)
			// don't worry about errors
			response, e := http.Get(url)
			if e != nil {
				log.Fatal(e)
			}
			defer response.Body.Close()

			//open a file for writing
			prefSplit := strings.Split(prefix, "/")
			file, err := os.Create(fmt.Sprintf("tmp/%s", prefSplit[len(prefSplit)-1]))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			// Use io.Copy to just dump the response body to the file. This supports huge files
			_, err = io.Copy(file, response.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Success!")
		}
	}
}
