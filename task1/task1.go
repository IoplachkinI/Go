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

type IdBook struct {
	Id int
	Book Book
}

type Storage interface {
	Init()
	Search(id int) (Book, bool)
	Add(id int, book Book)
	Rehash(generator IdGenerator)
	GetRawData() []IdBook
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
	newStorage := make(map[int]Book)
	for _, book := range ms.storage {
		newStorage[generator.GenerateId(book.Name)] = book
	}
	ms.storage = newStorage
}

func (ms* MapStorage) GetRawData() []IdBook {
	tmpSlice := make([]IdBook, 0, 0)
	for id, book := range ms.storage {
		tmpSlice = append(tmpSlice, IdBook{id, book})
	}
	return tmpSlice
}

type SliceStorage struct {
	storage [] IdBook
}

func (ss *SliceStorage) Init() {
	ss.storage = make([]IdBook, 0, 0)
}

func (ss *SliceStorage) Add(id int, book Book) {
	ss.storage = append(ss.storage, IdBook{id, book})
}

func (ss *SliceStorage) Search(id int) (Book, bool) {
	for _, idBook := range ss.storage {
		if (idBook.Id == id) {
			return idBook.Book, true;
		}
	}
	return Book{}, false;
}

func (ss *SliceStorage) Rehash(generator IdGenerator) {
	newStorage := make([]IdBook, 0, 0)
	for _, idBook := range ss.storage {
		newStorage = append(newStorage, 
			IdBook{generator.GenerateId(idBook.Book.Name), idBook.Book})
	}
	ss.storage = newStorage
}

func (ss *SliceStorage) GetRawData() []IdBook {
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
		charCode := (int(char) - '0')
		hash = (hash + (charCode * deg) % hasher.Mod) % hasher.Mod
		deg = (deg * hasher.Mult) % hasher.Mod
	}
	return hash
}

type SimpleHasher struct {}

func (hasher SimpleHasher) GenerateId(str string) int {
	sum := 0
	for _, char := range str {
		charCode := (int(char) - '0')
		sum += charCode * charCode
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
	for _, idBook := range lib.storage.GetRawData() {
		storage.Add(idBook.Id, idBook.Book);
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

func NewLibrary() Library {
	books := []Book {
		{"1984", "Джордж Оруэлл", "Роман-антиутопия", 1948},
		{"Обломов", "Гончаров И. А", "Роман", 1859},
		{"Бедная Лиза", "Карамзин Н. М.", "Повесть", 1791},
		{"Гранатовый браслет", "Куприн А. И.", "Повесть", 1910},
	}

	simpleHasher := SimpleHasher{}

	lib := Library{storage: &MapStorage{}, generator: simpleHasher}
	lib.Init()

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
	
	polyHasher := PolynomialHasher{107_381_377, 11}
	lib.ChangeGenerator(polyHasher)
	fmt.Println("Generator function changed!")

	TestEntry("1984", &lib)
	TestEntry("Бедная Лиза", &lib)

	lib.ChangeStorage(&SliceStorage{})
	fmt.Println("Storage changed!")

	TestEntry("Обломов", &lib)
	TestEntry("Обломовъ", &lib)
	TestEntry("Гранатовый браслет", &lib)
}