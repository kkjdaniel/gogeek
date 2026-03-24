//go:build contract

package contract

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/beevik/etree"
	gogeek "github.com/kkjdaniel/gogeek/v2"
)

// fetchRawXML makes an authenticated, rate-limited request and returns the raw XML bytes.
func fetchRawXML(client *gogeek.Client, url string) ([]byte, error) {
	client.Limiter().Take()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	switch client.AuthMode() {
	case gogeek.AuthAPIKey:
		req.Header.Set("Authorization", "Bearer "+client.APIKey())
	case gogeek.AuthCookie:
		req.Header.Set("Cookie", client.CookieString())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle 202 (processing) with a single retry
	if resp.StatusCode == http.StatusAccepted {
		time.Sleep(3 * time.Second)
		client.Limiter().Take()
		req2, _ := http.NewRequest("GET", url, nil)
		switch client.AuthMode() {
		case gogeek.AuthAPIKey:
			req2.Header.Set("Authorization", "Bearer "+client.APIKey())
		case gogeek.AuthCookie:
			req2.Header.Set("Cookie", client.CookieString())
		}
		resp, err = http.DefaultClient.Do(req2)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// globalIgnored contains XML paths that appear in every BGG API response
// but are metadata we intentionally don't model.
var globalIgnored = map[string]bool{
	"[@termsofuse]": true, // present on every root element
}

// assertFieldCoverage parses raw XML and compares every element and attribute
// against the xml struct tags of the given model. It reports any XML fields
// present in the API response that the Go struct doesn't capture.
// Extra ignored path suffixes can be passed for endpoint-specific exclusions.
func assertFieldCoverage(t *testing.T, rawXML []byte, model interface{}, endpointName string, ignored ...string) {
	t.Helper()

	doc := etree.NewDocument()
	doc.ReadSettings.Permissive = true
	if err := doc.ReadFromBytes(rawXML); err != nil {
		t.Errorf("[%s] coverage: failed to parse XML: %v", endpointName, err)
		return
	}

	root := doc.Root()
	if root == nil {
		t.Errorf("[%s] coverage: no root element", endpointName)
		return
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	ignoreSet := map[string]bool{}
	for k, v := range globalIgnored {
		ignoreSet[k] = v
	}
	for _, p := range ignored {
		ignoreSet[p] = true
	}

	var unmatched []string
	checkElement(root, modelType, root.Tag, &unmatched)

	// Filter out ignored paths (match by suffix so "[@termsofuse]" matches
	// "items[@termsofuse]", "user[@termsofuse]", etc.)
	var filtered []string
	for _, u := range unmatched {
		skip := false
		for pattern := range ignoreSet {
			if strings.HasSuffix(u, pattern) || u == pattern {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, u)
		}
	}

	if len(filtered) > 0 {
		t.Errorf("[%s] API response contains fields not in model:\n  %s",
			endpointName, strings.Join(filtered, "\n  "))
	}
}

// checkElement recursively compares an XML element's children and attributes
// against the fields of a Go struct type.
func checkElement(elem *etree.Element, structType reflect.Type, path string, unmatched *[]string) {
	if structType.Kind() != reflect.Struct {
		return
	}

	// Build lookup maps from the struct's xml tags.
	elemFields := map[string]reflect.Type{}   // element name -> resolved field type
	attrFields := map[string]bool{}           // attribute name -> true
	nestedPaths := map[string]nestedField{}   // wrapper element -> child info
	hasChardata := false

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("xml")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		opts := parts[1:]

		isAttr := false
		for _, o := range opts {
			if o == "attr" {
				isAttr = true
			}
			if o == "chardata" {
				hasChardata = true
			}
		}

		if isAttr {
			attrFields[name] = true
			continue
		}

		if hasChardata {
			continue
		}

		// Handle nested paths like "statistics>ratings" or "threads>thread"
		if strings.Contains(name, ">") {
			segments := strings.SplitN(name, ">", 2)
			ft := resolveType(field.Type)
			nestedPaths[segments[0]] = nestedField{childName: segments[1], fieldType: ft}
			continue
		}

		elemFields[name] = resolveType(field.Type)
	}

	// Check attributes on this element.
	for _, attr := range elem.Attr {
		if attr.Space == "xmlns" || attr.Key == "xmlns" || attr.Space == "xml" {
			continue
		}
		if !attrFields[attr.Key] {
			*unmatched = append(*unmatched, fmt.Sprintf("%s[@%s]", path, attr.Key))
		}
	}

	// Check child elements.
	seen := map[string]bool{}
	for _, child := range elem.ChildElements() {
		tag := child.Tag
		if seen[tag] {
			continue // only check each unique element name once
		}
		seen[tag] = true

		childPath := path + "." + tag

		if ft, ok := elemFields[tag]; ok {
			checkElement(child, ft, childPath, unmatched)
		} else if nested, ok := nestedPaths[tag]; ok {
			// For "threads>thread" style: the wrapper element "threads" contains
			// child elements "thread" that map to the struct field.
			for _, grandchild := range child.ChildElements() {
				if grandchild.Tag == nested.childName {
					checkElement(grandchild, nested.fieldType, childPath+"."+nested.childName, unmatched)
					break
				}
			}
		} else {
			*unmatched = append(*unmatched, childPath)
		}
	}
}

type nestedField struct {
	childName string
	fieldType reflect.Type
}

// resolveType unwraps pointers and slices to get the underlying (usually struct) type.
func resolveType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t
}
