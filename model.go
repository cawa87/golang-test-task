package main

type Document struct {
	Url      string            `json:"url"`
	Meta     DocumentMeta      `json:"meta"`
	Elements []DocumentElement `json:"elements,omitempty"`
}

type RawDocument struct {
	Url     string
	Meta    DocumentMeta
	Content []byte
}

type DocumentMeta struct {
	Status        int    `json:"status"`
	ContentType   string `json:"conten-type,omitempty"`
	ContentLength int    `json:"content-length,omitempty"`
}

type DocumentElement struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}
