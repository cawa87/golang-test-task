package models

import (
	"sync"
	"encoding/json"
)

type Result struct{
	sync.Mutex
	urlData []*UrlData
}

func (r* Result) Add(data *UrlData)  {
	r.Lock()
	r.urlData = append(r.urlData, data)
	r.Unlock()
}

func (r* Result) JSON() ([]byte, error) {
	return json.Marshal(r.urlData)
}
