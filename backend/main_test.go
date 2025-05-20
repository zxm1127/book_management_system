package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-management-backend/models"
	"book-management-backend/storage"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test server with registered handlers
func setupTestServer(t *testing.T) *hertz.Hertz {
	h := server.New()

	h.POST("/books", addBookHandler)
	h.GET("/books", getAllBooksHandler)
	h.PUT("/books/:id", updateBookHandler)
	h.DELETE("/books/:id", deleteBookHandler)

	// Clear storage before each test
	storage.ClearBooks()

	return h
}

// Add a ClearBooks function to your storage package for testing
// You'll need to manually add this to backend/storage/storage.go

func TestAddBookHandler(t *testing.T) {
	h := setupTestServer(t)

	// Prepare test book data
	newBook := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "1234567890",
	}

	// Marshal book to JSON
	bookJSON, _ := json.Marshal(newBook)

	// Create a new HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a new RequestContext and run the handler
	c := h.NewContext(w, req)
	addBookHandler(c.Context(), c)

	// Assert the response status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Assert the response body (optional, but good practice)
	var responseBook models.Book
	json.Unmarshal(w.Body.Bytes(), &responseBook)
	assert.NotEmpty(t, responseBook.ID)
	assert.Equal(t, newBook.Title, responseBook.Title)
	assert.Equal(t, newBook.Author, responseBook.Author)
	assert.Equal(t, newBook.ISBN, responseBook.ISBN)

	// Verify the book was added to storage
	addedBook, found := storage.GetBook(responseBook.ID)
	assert.True(t, found)
	assert.Equal(t, responseBook, addedBook)
}

// TODO: Add tests for getAllBooksHandler, updateBookHandler, and deleteBookHandler

func TestGetAllBooksHandler(t *testing.T) {
	h := setupTestServer(t)

	// Add some books for testing
	book1 := models.Book{ID: "1", Title: "Book 1", Author: "Author 1", ISBN: "ISBN 1"}
	book2 := models.Book{ID: "2", Title: "Book 2", Author: "Author 2", ISBN: "ISBN 2"}
	storage.AddBook(book1)
	storage.AddBook(book2)

	// Create a new HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)

	// Create a new RequestContext and run the handler
	c := h.NewContext(w, req)
	getAllBooksHandler(c.Context(), c)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	var responseBooks []models.Book
	json.Unmarshal(w.Body.Bytes(), &responseBooks)
	assert.Len(t, responseBooks, 2)
	// You might want to add more specific assertions about the content of the books
}

func TestUpdateBookHandler(t *testing.T) {
	h := setupTestServer(t)

	// Add a book to update
	existingBook := models.Book{ID: "1", Title: "Old Title", Author: "Old Author", ISBN: "Old ISBN"}
	storage.AddBook(existingBook)

	// Prepare updated book data
	updatedBookData := models.Book{ID: "1", Title: "New Title", Author: "New Author", ISBN: "New ISBN"}
	updatedBookJSON, _ := json.Marshal(updatedBookData)

	// Create a new HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer(updatedBookJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a new RequestContext and run the handler
	c := h.NewContext(w, req)
	updateBookHandler(c.Context(), c)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the book was updated in storage
	bookInStorage, found := storage.GetBook("1")
	assert.True(t, found)
	assert.Equal(t, updatedBookData.Title, bookInStorage.Title)
	assert.Equal(t, updatedBookData.Author, bookInStorage.Author)
	assert.Equal(t, updatedBookData.ISBN, bookInStorage.ISBN)
}

func TestDeleteBookHandler(t *testing.T) {
	h := setupTestServer(t)

	// Add a book to delete
	existingBook := models.Book{ID: "1", Title: "Book to Delete", Author: "Author", ISBN: "ISBN"}
	storage.AddBook(existingBook)

	// Verify the book exists before deletion
	_, foundBefore := storage.GetBook("1")
	assert.True(t, foundBefore)

	// Create a new HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/1", nil)

	// Create a new RequestContext and run the handler
	c := h.NewContext(w, req)
	deleteBookHandler(c.Context(), c)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the book was deleted from storage
	_, foundAfter := storage.GetBook("1")
	assert.False(t, foundAfter)
}
