package crud

import "github.com/dohr-michael/go-libs/storage"

type Created struct {
	Id   storage.ID  `json:"id"`
	Item interface{} `json:"item"`
}

type Deleted struct {
	Id storage.ID `json:"id"`
}
