package database

type Student struct {
	Surname   string `json:"surname"`
	Group     string `json:"group"`
	ID        int64  `json:"id"`
	Value     int64  `json:"value"`
	IsPresent int    `json:"isPresent"`
}
