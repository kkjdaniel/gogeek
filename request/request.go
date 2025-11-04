package request

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/clbanning/mxj"
	"github.com/kkjdaniel/gogeek/v2"
)

var (
	maxRetries = 5
	retryDelay = 3 * time.Second
)

var (
	// ErrEmptyResponse is returned when the response body is empty
	ErrEmptyResponse = fmt.Errorf("empty response body")
	// ErrHTTPError is returned when the HTTP request fails
	ErrHTTPError = fmt.Errorf("HTTP request failed")
	// ErrUnexpectedStatusCode is returned when the HTTP status code is not 200
	ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")
	// ErrMaxRetriesExceeded is returned when the maximum number of retries is exceeded
	ErrMaxRetriesExceeded = fmt.Errorf("exceeded maximum retries while waiting for BGG to process request")
	// ErrUnmarshalError is returned when the XML response cannot be unmarshalled
	ErrUnmarshalError = fmt.Errorf("failed to unmarshal XML response")
	// ErrRegenerateError is returned when the XML response cannot be regenerated
	ErrRegenerateError = fmt.Errorf("failed to regenerate XML response")
	// ErrXMLParseError is returned when the XML response cannot be parsed
	ErrXMLParseError = fmt.Errorf("failed to parse XML response")
)

func FetchAndUnmarshal(client *gogeek.Client, url string, v interface{}) error {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		client.Limiter().Take()

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ErrHTTPError
		}

		// Add authentication headers based on client configuration
		switch client.AuthMode() {
		case gogeek.AuthAPIKey:
			req.Header.Set("Authorization", "Bearer "+client.APIKey())
		case gogeek.AuthCookie:
			req.Header.Set("Cookie", client.CookieString())
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrHTTPError
		}

		// Handle 202 status - request accepted but still processing
		// https://boardgamegeek.com/wiki/page/BGG_XML_API2#toc12
		if resp.StatusCode == http.StatusAccepted {
			resp.Body.Close()
			if attempt == maxRetries {
				return ErrMaxRetriesExceeded
			}
			time.Sleep(retryDelay)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrEmptyResponse
		}

		body = fixMalformedXML(body)

		if err := xml.Unmarshal(body, v); err != nil {
			mv, err := mxj.NewMapXml(body)
			if err != nil {
				return ErrXMLParseError
			}

			cleanXML, err := mv.Xml()
			if err != nil {
				return ErrRegenerateError
			}

			if err := xml.Unmarshal(cleanXML, v); err != nil {
				typeName := fmt.Sprintf("%T", v)
				return fmt.Errorf("%w: failed to unmarshal into %s: %v", ErrUnmarshalError, typeName, err)
			}
		}

		return nil // Success
	}

	return fmt.Errorf("failed to get response from BGG API after retries")
}

func fixMalformedXML(data []byte) []byte {
	xmlStr := string(data)

	cdataRegex := regexp.MustCompile(`<!\[CDATA\[(.*?)\]\]>`)
	cdataSections := make(map[string]string)
	xmlStr = cdataRegex.ReplaceAllStringFunc(xmlStr, func(match string) string {
		placeholder := fmt.Sprintf("CDATA_PLACEHOLDER_%d", len(cdataSections))
		cdataSections[placeholder] = match
		return placeholder
	})

	entityRegex := regexp.MustCompile(`&(amp|lt|gt|apos|quot|#[0-9]+|#x[0-9a-fA-F]+);`)
	xmlStr = entityRegex.ReplaceAllStringFunc(xmlStr, func(s string) string {
		return "ENTITY_PLACEHOLDER" + s[1:]
	})

	htmlEntityRegex := regexp.MustCompile(`&([a-zA-Z]+);`)
	xmlStr = htmlEntityRegex.ReplaceAllStringFunc(xmlStr, func(s string) string {

		unescaped := html.UnescapeString(s)
		if unescaped != s {

			return unescaped
		}
		return s
	})

	xmlStr = strings.Replace(xmlStr, "&", "&amp;", -1)

	xmlStr = strings.Replace(xmlStr, "ENTITY_PLACEHOLDER", "&", -1)

	for placeholder, cdata := range cdataSections {
		xmlStr = strings.Replace(xmlStr, placeholder, cdata, -1)
	}

	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F]`)
	xmlStr = re.ReplaceAllString(xmlStr, "")

	return []byte(xmlStr)
}
