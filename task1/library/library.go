package library

import (
	"task1/book"
	"task1/storage"
	"task1/generator"
)

type Library struct {
	storage storage.Storage
	generator generator.IdGenerator;
}

func (lib* Library) Init(storage storage.Storage, generator generator.IdGenerator) {
	lib.storage = storage
	lib.storage.Init()
	lib.generator = generator
}

func (lib* Library) Search(name string) (book.Book, bool) {
	book, ok := lib.storage.Search(lib.generator.GenerateId(name))
	return book, ok
}

func (lib* Library) AddBook(book book.Book) {
	lib.storage.Add(lib.generator.GenerateId(book.Name), book)
}

func (lib* Library) ChangeGenerator(generator generator.IdGenerator) {
	lib.generator = generator
	lib.storage.Rehash(generator)
}


// Needs empty Storage object
func (lib* Library) ChangeStorage(storage storage.Storage) {
	for _, idBook := range lib.storage.GetRawData() {
		storage.Add(idBook.Id, idBook.Book);
	}
	lib.storage = storage
}