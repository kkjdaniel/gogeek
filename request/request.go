package request

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/clbanning/mxj"
	"go.uber.org/ratelimit"
)

// https://boardgamegeek.com/thread/2388502/updated-api-rate-limit-recommendation
var limiter = ratelimit.New(2, ratelimit.WithoutSlack)

func FetchAndUnmarshal(url string, v interface{}) error {
	limiter.Take()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from BGG API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	mv, err := mxj.NewMapXml(body)
	if err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}

	cleanXML, err := mv.Xml()
	if err != nil {
		return fmt.Errorf("failed to regenerate XML: %w", err)
	}

	if err := xml.Unmarshal(cleanXML, v); err != nil {
		return fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	return nil
}
