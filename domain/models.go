package domain

type Attraction struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Url         string  `json:"url"`
	ImageUrl    string  `json:"image_url"`
	Location    string  `json:"location"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
