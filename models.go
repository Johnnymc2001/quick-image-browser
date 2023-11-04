package main

type ImageObj struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Base64 string `json:"base64"`
}

type Config struct {
	LastBrowseFolder string `json:"lastBrowseFolder"`
}
