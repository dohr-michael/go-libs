package crud

type Created struct {
	Id   string      `json:"id"`
	Item interface{} `json:"item"`
}

type Deleted struct {
	Id string `json:"id"`
}
