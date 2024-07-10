package main


type Account struct {
	Id        int    `json:"id" db:"id" gorm:"primaryKey"`
	Firstname string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
}

