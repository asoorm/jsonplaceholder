package model

type Geo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     Geo    `json:"geo"`
}

type User struct {
	Id       int     `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Username string  `json:"username,omitempty"`
	Email    string  `json:"email,omitempty"`
	Address  Address `json:"address"`
}
