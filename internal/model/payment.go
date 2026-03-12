package model

type Payment struct {
	Id     int     `json:"id"`
	Person string  `json:"person"`
	Amount float32 `json:"amount"`
	Date   string  `json:"date"`
}

type Person struct {
	Name string `json:"name"`
}
