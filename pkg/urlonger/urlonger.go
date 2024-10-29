package urlonger

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Redirect struct {
	Url             string            `json:"url"`
	StatusCode      int               `json:"status_code"`
	Destination     string            `json:"destination"`
	ResponseHeaders map[string]string `json:"response_headers"`
}

func (redir Redirect) String() string {
	if redir.Destination == "" {
		return fmt.Sprintf("%s (%d)", redir.Url, redir.StatusCode)
	}
	return fmt.Sprintf("%s (%d) -> %s", redir.Url, redir.StatusCode, redir.Destination)
}

func NewRedirect(url string) (Redirect, error) {
	return Redirect{
		Url:             url,
		ResponseHeaders: map[string]string{},
	}, nil
}

func MustNewRedirect(url string) Redirect {
	redir, err := NewRedirect(url)
	if err != nil {
		panic(err)
	}
	return redir
}

func Resolve(url string, headersFilter []string) ([]Redirect, error) {

	tr := http.Transport{}
	client := http.Client{
		Transport: &tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	locations := []string{url}
	redirs := []Redirect{}

	for {
		redir := MustNewRedirect(locations[len(locations)-1])

		req, err := http.NewRequest("HEAD", redir.Url, nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// Set response status code
		redir.StatusCode = resp.StatusCode
		// Set location
		redir.Destination = resp.Header.Get("location")
		// Only pull certain headers
		for _, headerName := range headersFilter {
			headerVal := resp.Header.Get(headerName)
			if headerVal != "" {
				redir.ResponseHeaders[headerName] = headerVal
			}
		}

		log.Println(redir.String())

		redirs = append(redirs, redir)

		if redir.Destination == "" {
			// end of chain
			break
		}

		for _, prevLocation := range locations {
			if strings.EqualFold(redir.Destination, prevLocation) {
				// loop detected
				log.Println("loop detected")
				break
			}
		}

		locations = append(locations, redir.Destination)
	}

	return redirs, nil
}
