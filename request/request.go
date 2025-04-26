package request

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/clbanning/mxj"
	"go.uber.org/ratelimit"
)

// https://boardgamegeek.com/thread/2388502/updated-api-rate-limit-recommendation
var limiter = ratelimit.New(2, ratelimit.WithoutSlack)

var (
	maxRetries = 5
	retryDelay = 3 * time.Second
)

func FetchAndUnmarshal(url string, v interface{}) error {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		limiter.Take()

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch data from BGG API: %w", err)
		}

		// Handle 202 status - request accepted but still processing
		// https://boardgamegeek.com/wiki/page/BGG_XML_API2#toc12
		if resp.StatusCode == http.StatusAccepted {
			resp.Body.Close()
			if attempt == maxRetries {
				return fmt.Errorf("exceeded maximum retries while waiting for BGG to process request")
			}
			time.Sleep(retryDelay)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		body = fixMalformedXML(body)

		if err := xml.Unmarshal(body, v); err != nil {
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
		}

		return nil // Success
	}

	return fmt.Errorf("failed to get response from BGG API after retries")
}

func fixMalformedXML(data []byte) []byte {
	xmlStr := string(data)

	re := regexp.MustCompile(`&(amp|lt|gt|apos|quot|#[0-9]+|#x[0-9a-fA-F]+);`)
	xmlStr = re.ReplaceAllStringFunc(xmlStr, func(s string) string {
		return "ENTITY_PLACEHOLDER" + s[1:]
	})
	xmlStr = strings.Replace(xmlStr, "&", "&amp;", -1)
	xmlStr = strings.Replace(xmlStr, "ENTITY_PLACEHOLDER", "&", -1)

	re = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F]`)
	xmlStr = re.ReplaceAllString(xmlStr, "")

	return []byte(xmlStr)
}
