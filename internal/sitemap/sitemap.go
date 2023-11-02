// Package sitemap provides a simple way to parse a sitemap and get a list of
// URLs.
package sitemap

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

// ErrSitemap is returned when a sitemap cannot be parsed.
const ErrSitemap xerrors.Error = "failed to parse sitemap"

// AverageSitemapSize is the average size of a sitemap.
const AverageSitemapSize = 1000

// URL represents a single URL element in the XML sitemap.
type URL struct {
	Loc string `xml:"loc"`
}

// Parse reads a sitemap from an io.Reader an returns a slice of URLs.
func Parse(r io.Reader) ([]string, error) {
	var (
		urls           = make([]string, 0, AverageSitemapSize)
		bufferedReader = bufio.NewReader(r)
		decoder        = xml.NewDecoder(bufferedReader)
		inURL          = false
		inImage        = false
	)

	for {
		t, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrSitemap, err)
		}

		switch elem := t.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "url":
				inURL = true
			case "image":
				inImage = true
			case "loc":
				if inURL && !inImage {
					var loc string

					if err := decoder.DecodeElement(&loc, &elem); err != nil {
						return nil, fmt.Errorf("%w: %w", ErrSitemap, err)
					}

					urls = append(urls, loc)
				}
			}
		case xml.EndElement:
			switch elem.Name.Local {
			case "url":
				inURL = false
			case "image":
				inImage = false
			}
		}
	}

	return urls, nil
}
