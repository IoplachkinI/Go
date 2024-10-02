package storage

import (
	"task1/book"
	"task1/generator"
)

type Storage interface {
	Init()
	Search(id int) (book.Book, bool)
	Add(id int, book book.Book)
	Rehash(generator generator.IdGenerator)
	GetRawData() []book.IdBook
}

type MapStorage struct {
	storage map[int]book.Book
}

func (ms *MapStorage) Init() {
	ms.storage = make(map[int]book.Book)
}

func (ms *MapStorage) Add(id int, book book.Book) {
	ms.storage[id] = book
}

func (ms *MapStorage) Search(id int) (book.Book, bool) {
	book, ok := ms.storage[id]
	return book, ok
}

func (ms *MapStorage) Rehash (generator generator.IdGenerator) {
	newStorage := make(map[int]book.Book)
	for _, book := range ms.storage {
		newStorage[generator.GenerateId(book.Name)] = book
	}
	ms.storage = newStorage
}

func (ms* MapStorage) GetRawData() []book.IdBook {
	tmpSlice := make([]book.IdBook, 0, 0)
	for id, b := range ms.storage {
		tmpSlice = append(tmpSlice, book.IdBook{Id:id, Book:b})
	}
	return tmpSlice
}

type SliceStorage struct {
	storage [] book.IdBook
}

func (ss *SliceStorage) Init() {
	ss.storage = make([]book.IdBook, 0, 0)
}

func (ss *SliceStorage) Add(id int, b book.Book) {
	ss.storage = append(ss.storage, book.IdBook{Id:id, Book:b})
}

func (ss *SliceStorage) Search(id int) (book.Book, bool) {
	for _, idBook := range ss.storage {
		if (idBook.Id == id) {
			return idBook.Book, true;
		}
	}
	return book.Book{}, false;
}

func (ss *SliceStorage) Rehash(generator generator.IdGenerator) {
	newStorage := make([]book.IdBook, 0, 0)
	for _, idBook := range ss.storage {
		newStorage = append(newStorage, 
			book.IdBook{Id: generator.GenerateId(idBook.Book.Name), Book:idBook.Book})
	}
	ss.storage = newStorage
}

func (ss *SliceStorage) GetRawData() []book.IdBook {
	return ss.storage
}