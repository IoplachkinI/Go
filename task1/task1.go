package main

import (
	"fmt"
)

type Book struct {
	Name string
	Author string
	Genre string
	Year int
}

type Storage interface {
	Init()
	Search(id int) (Book, bool)
	Add(id int, book Book)
	Rehash(generator IdGenerator)
	GetRawData() []struct {int; Book}
}

type MapStorage struct {
	storage map[int]Book
}

func (ms *MapStorage) Init() {
	ms.storage = make(map[int]Book)
}

func (ms *MapStorage) Add(id int, book Book) {
	ms.storage[id] = book
}

func (ms *MapStorage) Search(id int) (Book, bool) {
	book, ok := ms.storage[id]
	return book, ok
}

func (ms *MapStorage) Rehash (generator IdGenerator) {
	new_storage := make(map[int]Book)
	for _, book := range ms.storage {
		new_storage[generator.GenerateId(book.Name)] = book
	}
	ms.storage = new_storage
}

func (ms* MapStorage) GetRawData() []struct {int; Book} {
	tmp_slice := make([]struct {int; Book}, 0, 0)
	for id, book := range ms.storage {
		tmp_slice = append(tmp_slice, struct {int; Book}{id, book})
	}
	return tmp_slice
}

type SliceStorage struct {
	storage [] struct {int; Book}
}

func (ss *SliceStorage) Init() {
	ss.storage = make([]struct {int; Book}, 0, 0)
}

func (ss *SliceStorage) Add(id int, book Book) {
	ss.storage = append(ss.storage, struct {int; Book}{id, book})
}

func (ss *SliceStorage) Search(id int) (Book, bool) {
	for _, pair := range ss.storage {
		var book_id, book = pair.int, pair.Book
		if (book_id == id) {
			return book, true;
		}
	}
	return Book{}, false;
}

func (ss *SliceStorage) Rehash(generator IdGenerator) {
	new_storage := make([]struct {int; Book}, 0, 0)
	for _, pair := range ss.storage {
		new_storage = append(new_storage, 
			struct {int; Book}{generator.GenerateId(pair.Name), pair.Book})
	}
	ss.storage = new_storage
}

func (ss *SliceStorage) GetRawData() []struct {int; Book} {
	return ss.storage
}

type IdGenerator interface {
	GenerateId(name string) int
}

type PolynomialHasher struct {
	Mod int
	Mult int
}

func (hasher PolynomialHasher) GenerateId(str string) int {
	var hash, deg = 0, 1
	for _, char := range str {
		char_code := (int(char) - '0')
		hash = (hash + (char_code * deg) % hasher.Mod) % hasher.Mod
		deg = (deg * hasher.Mult) % hasher.Mod
	}
	return hash
}

type SimpleHasher struct {}

func (hasher SimpleHasher) GenerateId(str string) int {
	sum := 0
	for _, char := range str {
		char_code := (int(char) - '0')
		sum += char_code * char_code
	}
	return sum
}

type Library struct {
	storage Storage
	generator IdGenerator;
}

func (lib* Library) Init() {
	lib.storage.Init()
}

func (lib* Library) Search(name string) (Book, bool) {
	book, ok := lib.storage.Search(lib.generator.GenerateId(name))
	return book, ok
}

func (lib* Library) AddBook(book Book) {
	lib.storage.Add(lib.generator.GenerateId(book.Name), book)
}

func (lib* Library) ChangeGenerator(generator IdGenerator) {
	lib.generator = generator
	lib.storage.Rehash(generator)
}


// Needs empty Storage object
func (lib* Library) ChangeStorage(storage Storage) {
	for _, pair := range lib.storage.GetRawData() {
		storage.Add(pair.int, pair.Book);
	}
	lib.storage = storage
}

func GetInfoByName(name string, lib *Library) {
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

func TestEntry(name string, lib *Library) {
	_, ok := lib.Search(name)
	if (ok) {
		fmt.Println("Book found: ", name)
	} else {
		fmt.Println("Book NOT found:", name)
	}
}

func main() {
	books := []Book {
		{"1984", "Джордж Оруэлл", "Роман-антиутопия", 1948},
		{"Обломов", "Гончаров И. А", "Роман", 1859},
		{"Бедная Лиза", "Карамзин Н. М.", "Повесть", 1791},
		{"Гранатовый браслет", "Куприн А. И.", "Повесть", 1910},
	}

	poly_hasher := PolynomialHasher{107_381_377, 11}
	simple_hasher := SimpleHasher{}

	lib := Library{storage: &MapStorage{}, generator: simple_hasher}
	lib.Init()

	for _, book := range books {
		lib.AddBook(book)
	}

	GetInfoByName("1984", &lib)

	TestEntry("Обломов", &lib)
	TestEntry("Капитанская дочка", &lib)

	fmt.Println("Generator function changed!")
	lib.ChangeGenerator(poly_hasher)

	TestEntry("1984", &lib)
	TestEntry("Бедная Лиза", &lib)

	fmt.Println("Storage changed!")
	lib.ChangeStorage(&SliceStorage{})

	TestEntry("Обломов", &lib)
	TestEntry("Обломовъ", &lib)
	TestEntry("Гранатовый браслет", &lib)
}