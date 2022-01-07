package download

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

// HttpRequester is an abstraction for a default HTTP GET requests.
type HttpRequester interface {
	// GetContent fetches the content of the remote url.
	GetContent(url string) (string, error)
}

// NewHttpRequester creates a new HttpRequester with the give client.
func NewHttpRequester(client *http.Client) HttpRequester {
	return &httpRequester{client: client}
}

// httpRequester implements the HttpRequester interface
type httpRequester struct {
	client *http.Client
}

func (hr *httpRequester) GetContent(url string) (content string, err error) {
	var rsp *http.Response
	rsp, err = hr.client.Get(url)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Download failed")
		return
	}

	statuscode := rsp.StatusCode

	if statuscode != http.StatusOK {
		err = fmt.Errorf("server response status: %d", rsp.StatusCode)
		log.Err(err).Str("url", url).Msg("Download failed")
		return
	}
	var rawBody []byte
	rawBody, err = io.ReadAll(rsp.Body)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Download failed to read body")
		return
	}
	content = string(rawBody)
	return
}
