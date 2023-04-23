package entity

type Product struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
