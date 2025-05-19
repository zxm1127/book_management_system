package main

import (
	"book-management-backend/models"
	"book-management-backend/storage"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"github.com/hertz-contrib/cors"
	"time"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8888"))

	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},                   // 允许你的前端源访问
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type"},                  // 允许的请求头
		AllowCredentials: true,                                                // 允许携带 cookies 等凭证
		MaxAge:           12 * time.Hour,                                      // CORS 预检请求的缓存时间
	}))

	h.POST("/books", addBookHandler)
	h.GET("/books", getAllBooksHandler)
	h.PUT("/books/:id", updateBookHandler)
	h.DELETE("/books/:id", deleteBookHandler)

	h.Spin()
}

func addBookHandler(ctx context.Context, c *app.RequestContext) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Generate a unique ID for the new book

	book.ID = uuid.New().String()

	storage.AddBook(book)

	c.JSON(consts.StatusCreated, book)
}

func getAllBooksHandler(ctx context.Context, c *app.RequestContext) {
	books := storage.GetAllBooks()
	c.JSON(consts.StatusOK, books)
}

func updateBookHandler(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	var updatedBook models.Book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	updatedBook.ID = id // Ensure the ID from the URL is used

	if storage.UpdateBook(updatedBook) {
		c.Status(consts.StatusOK)
	} else {
		c.JSON(consts.StatusNotFound, map[string]string{"error": "Book not found"})
	}
}

func deleteBookHandler(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if storage.DeleteBook(id) {
		c.Status(consts.StatusOK)
	} else {
		c.JSON(consts.StatusNotFound, map[string]string{"error": "Book not found"})
	}
}
