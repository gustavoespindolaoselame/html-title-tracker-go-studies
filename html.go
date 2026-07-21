package titletracker

import (
	"io"
	"net/http"
	"regexp"
)

// Titulo fetches the <title> tag asynchronously for given URLs.
// Capitalized 'T' makes it exported (public).
func Titulo(urls ...string) <-chan string {
	c := make(chan string)
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				c <- "Error fetching URL"
				return
			}
			defer resp.Body.Close()

			html, err := io.ReadAll(resp.Body)
			if err != nil {
				c <- "Error reading body"
				return
			}

			r, _ := regexp.Compile(`(?i)<title>(.*?)</title>`)
			matches := r.FindStringSubmatch(string(html))
			if len(matches) > 1 {
				c <- matches[1]
			} else {
				c <- "No title found"
			}
		}(url)
	}
	return c
}