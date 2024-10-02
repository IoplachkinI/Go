package book

type Book struct {
	Name string
	Author string
	Genre string
	Year int
}

type IdBook struct {
	Id int
	Book Book
}