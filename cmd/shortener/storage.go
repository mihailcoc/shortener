package main

type Body struct {
	URL string `json:"str_value"`
}

var (
	urls = make(map[string]string)
)
