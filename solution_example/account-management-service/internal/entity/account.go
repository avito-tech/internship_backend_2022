package entity

import "time"

type Account struct {
	Id        int       `db:"id"`
	Balance   int       `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}
