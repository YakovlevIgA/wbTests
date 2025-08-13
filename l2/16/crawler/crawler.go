package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"wget/fetcher"
	"wget/parser"
	"wget/saver"
)

type Crawler struct {
	BaseURL     *url.URL
	OutputDir   string
	MaxDepth    int
	Client      *http.Client
	Concurrency int

	visited   map[string]bool
	visitedMu sync.Mutex
	queue     chan Task
	wg        sync.WaitGroup
}

type Task struct {
	URL   *url.URL
	Depth int
}

func NewCrawler(u *url.URL, out string, depth int, client *http.Client, concurrency int) *Crawler {
	return &Crawler{
		BaseURL:     u,
		OutputDir:   out,
		MaxDepth:    depth,
		Client:      client,
		Concurrency: concurrency,
		visited:     make(map[string]bool),
		queue:       make(chan Task, 100000),
	}
}

func (c *Crawler) Run() {
	for i := 0; i < c.Concurrency; i++ {
		go c.worker()
	}

	c.enqueue(c.BaseURL, 0)

	c.wg.Wait()
	close(c.queue)
}

func (c *Crawler) worker() {
	for task := range c.queue {
		c.process(task)
		c.wg.Done()
	}
}

func (c *Crawler) process(task Task) {
	if task.Depth > c.MaxDepth {
		return
	}

	if c.isVisited(task.URL) {
		return
	}
	c.markVisited(task.URL)

	body, contentType, err := fetcher.Fetch(task.URL, c.Client)
	if err != nil {
		log.Printf("[ERROR] %s: %v", task.URL, err)
		return
	}

	localPath := saver.Save(task.URL, body, contentType, c.OutputDir)
	log.Printf("[SAVED] %s → %s", task.URL, localPath)

	if strings.HasPrefix(contentType, "text/html") {
		links, resources := parser.ExtractLinksAndResources(task.URL, body)

		for _, res := range resources {
			if res.Host == c.BaseURL.Host && !c.isVisited(res) {
				c.enqueue(res, task.Depth)
			}
		}

		for _, link := range links {
			if link.Host == c.BaseURL.Host && !c.isVisited(link) {
				c.enqueue(link, task.Depth+1)
			}
		}
	}
}

func (c *Crawler) enqueue(u *url.URL, depth int) {
	if depth > c.MaxDepth {
		return
	}

	c.wg.Add(1)
	c.queue <- Task{URL: u, Depth: depth}
}

func (c *Crawler) isVisited(u *url.URL) bool {
	c.visitedMu.Lock()
	defer c.visitedMu.Unlock()
	return c.visited[u.String()]
}

func (c *Crawler) markVisited(u *url.URL) {
	c.visitedMu.Lock()
	defer c.visitedMu.Unlock()
	c.visited[u.String()] = true
}
