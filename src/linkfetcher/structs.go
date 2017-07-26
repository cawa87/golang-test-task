package main

// Meta contains basic fetching info
type Meta struct {
	Status        int    `json:"status"`
	ContentType   string `json:"content-type,omitempty"`
	ContentLength int    `json:"content-length,omitempty"`
	Error         string `json:"error,omitempty"`
}

// Element contains the name and the count of the HTML tag
type Element struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

// ResponseItem is a definition of one link
type ResponseItem struct {
	URL      string    `json:"url"`
	Meta     Meta      `json:"meta"`
	Elements []Element `json:"elements,omitempty"`
}

type Response []*ResponseItem

type Request []string
