// internal/domain/feed/atom.go
package feed

import (
	"encoding/xml"
	"fmt"
	"time"

	"gemini-bridge/internal/domain/model"
)

// Atom 1.0 フィード構造体（RFC 4287準拠）
type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	XMLNS   string      `xml:"xmlns,attr"`
	Title   string      `xml:"title"`
	Link    []atomLink  `xml:"link"`
	Updated string      `xml:"updated"`
	Author  atomAuthor  `xml:"author"`
	ID      string      `xml:"id"`
	Entries []atomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

type atomEntry struct {
	Title   string       `xml:"title"`
	Link    []atomLink   `xml:"link"`
	ID      string       `xml:"id"`
	Updated string       `xml:"updated"`
	Summary string       `xml:"summary,omitempty"`
	Content *atomContent `xml:"content,omitempty"`
}

type atomContent struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// GenerateAtom はサイト情報からAtom 1.0フィードXMLを生成する。
func GenerateAtom(site model.Site) (string, error) {
	var updated string
	if len(site.Posts) > 0 {
		updated = site.Posts[0].Date.Format(time.RFC3339)
	} else {
		updated = time.Now().Format(time.RFC3339)
	}

	f := atomFeed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		Title:   site.Title,
		Updated: updated,
		ID:      site.BaseURL + "/",
		Author:  atomAuthor{Name: site.Author},
		Link: []atomLink{
			{Href: site.BaseURL + "/", Rel: "alternate", Type: "text/html"},
			{Href: site.BaseURL + "/feed/atom.xml", Rel: "self", Type: "application/atom+xml"},
		},
	}

	for _, post := range site.Posts {
		postURL := fmt.Sprintf("%s/posts/%s/", site.BaseURL, post.Slug)
		entry := atomEntry{
			Title:   post.Title,
			ID:      postURL,
			Updated: post.Date.Format(time.RFC3339),
			Summary: post.Description,
			Link: []atomLink{
				{Href: postURL, Rel: "alternate", Type: "text/html"},
				{Href: postURL + "index.gmi", Rel: "alternate", Type: "text/gemini"},
			},
		}
		f.Entries = append(f.Entries, entry)
	}

	output, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Atom feed: %w", err)
	}

	return xml.Header + string(output), nil
}
