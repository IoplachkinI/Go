package main

import (
	"fmt"
	"task1/book"
	"task1/generator"
	"task1/library"
	"task1/storage"
)

func GetInfoByName(name string, lib *library.Library) {
	book, ok := lib.Search(name)
	if (ok) {
		fmt.Println("Book Found:")
		fmt.Println("Name:", book.Name)
		fmt.Println("Author:", book.Author)
		fmt.Println("Genre:", book.Genre)
		fmt.Println("Year:", book.Year)
	} else {
		fmt.Println("Error: no book named", name)
	}
}

func TestEntry(name string, lib *library.Library) {
	_, ok := lib.Search(name)
	if (ok) {
		fmt.Println("Book found: ", name)
	} else {
		fmt.Println("Book NOT found:", name)
	}
}

func NewLibrary() library.Library {
	books := []book.Book {
		{"1984", "Джордж Оруэлл", "Роман-антиутопия", 1948},
		{"Обломов", "Гончаров И. А", "Роман", 1859},
		{"Бедная Лиза", "Карамзин Н. М.", "Повесть", 1791},
		{"Гранатовый браслет", "Куприн А. И.", "Повесть", 1910},
	}

	simpleHasher := generator.SimpleHasher{}

	lib := library.Library{}
	lib.Init(&storage.MapStorage{}, simpleHasher)

	for _, book := range books {
		lib.AddBook(book)
	}

	return lib
}

func main() {
	lib := NewLibrary();

	GetInfoByName("1984", &lib)

	TestEntry("Обломов", &lib)
	TestEntry("Капитанская дочка", &lib)

	polyHasher := generator.PolynomialHasher{Mod: 107_381_377, Mult:11}
	lib.ChangeGenerator(polyHasher)
	fmt.Println("Generator function changed!")

	TestEntry("1984", &lib)
	TestEntry("Бедная Лиза", &lib)

	lib.ChangeStorage(&storage.SliceStorage{})
	fmt.Println("Storage changed!")

	TestEntry("Обломов", &lib)
	TestEntry("Обломовъ", &lib)
	TestEntry("Гранатовый браслет", &lib)
}