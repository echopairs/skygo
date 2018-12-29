package model

type Book struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"Title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}
