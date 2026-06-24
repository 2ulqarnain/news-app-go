package models

type News struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title"`
	NewsURL  string `json:"newsURL"`
	ImageURL string `json:"imageURL"`
}
