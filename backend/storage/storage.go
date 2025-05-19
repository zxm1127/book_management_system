package storage

import (
	"sync"

	"book-management-backend/models"
)

var (
	books = make(map[string]models.Book)
	mu    sync.Mutex
)

func AddBook(book models.Book) {
	mu.Lock()
	defer mu.Unlock()
	books[book.ID] = book
}

func GetBook(id string) (models.Book, bool) {
	mu.Lock()
	defer mu.Unlock()
	book, ok := books[id]
	return book, ok
}

func GetAllBooks() []models.Book {
	mu.Lock()
	defer mu.Unlock()
	var bookList []models.Book
	for _, book := range books {
		bookList = append(bookList, book)
	}
	return bookList
}

func UpdateBook(book models.Book) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := books[book.ID]; ok {
		books[book.ID] = book
		return true
	}
	return false
}

func DeleteBook(id string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := books[id]; ok {
		delete(books, id)
		return true
	}
	return false
} 
