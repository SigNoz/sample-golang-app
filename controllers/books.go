package controllers

import (
	"log/slog"
	"net/http"

	"github.com/SigNoz/sample-golang-app/models"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// GET /books
// Find all books
func FindBooks(c *gin.Context) {
	ctx := c.Request.Context()
	slog.InfoContext(ctx, "list books request",
		slog.String("http.method", c.Request.Method),
		slog.String("http.path", c.Request.URL.Path),
	)
	var books []models.Book
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("This is a sample event", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("sampleAttribute", "Test")))
	models.DB.WithContext(ctx).Find(&books)
	slog.InfoContext(ctx, "list books success",
		slog.Int("http.status", http.StatusOK),
		slog.Int("books.count", len(books)),
	)
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	ctx := c.Request.Context()
	bookID := c.Param("id")
	slog.InfoContext(ctx, "get book request",
		slog.String("http.method", c.Request.Method),
		slog.String("http.path", c.Request.URL.Path),
		slog.String("book.id", bookID),
	)
	var book models.Book
	if err := models.DB.WithContext(ctx).Where("id = ?", bookID).First(&book).Error; err != nil {
		slog.WarnContext(ctx, "book not found",
			slog.String("book.id", bookID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	slog.InfoContext(ctx, "get book success",
		slog.Int("http.status", http.StatusOK),
		slog.String("book.id", bookID),
		slog.String("book.title", book.Title),
	)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {
	ctx := c.Request.Context()
	slog.InfoContext(ctx, "create book request",
		slog.String("http.method", c.Request.Method),
		slog.String("http.path", c.Request.URL.Path),
	)
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.WarnContext(ctx, "create book validation failed",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.WithContext(ctx).Create(&book)
	slog.InfoContext(ctx, "create book success",
		slog.Int("http.status", http.StatusOK),
		slog.Any("book.id", book.ID),
		slog.String("book.title", book.Title),
		slog.String("book.author", book.Author),
	)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	ctx := c.Request.Context()
	bookID := c.Param("id")
	slog.InfoContext(ctx, "update book request",
		slog.String("http.method", c.Request.Method),
		slog.String("http.path", c.Request.URL.Path),
		slog.String("book.id", bookID),
	)
	var book models.Book
	if err := models.DB.WithContext(ctx).Where("id = ?", bookID).First(&book).Error; err != nil {
		slog.WarnContext(ctx, "update book: record not found",
			slog.String("book.id", bookID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.WarnContext(ctx, "update book validation failed",
			slog.String("book.id", bookID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.WithContext(ctx).Model(&book).Updates(input)
	slog.InfoContext(ctx, "update book success",
		slog.Int("http.status", http.StatusOK),
		slog.String("book.id", bookID),
		slog.String("book.title", book.Title),
		slog.String("book.author", book.Author),
	)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()
	bookID := c.Param("id")
	slog.InfoContext(ctx, "delete book request",
		slog.String("http.method", c.Request.Method),
		slog.String("http.path", c.Request.URL.Path),
		slog.String("book.id", bookID),
	)
	var book models.Book
	if err := models.DB.WithContext(ctx).Where("id = ?", bookID).First(&book).Error; err != nil {
		slog.WarnContext(ctx, "delete book: record not found",
			slog.String("book.id", bookID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&book)
	slog.InfoContext(ctx, "delete book success",
		slog.Int("http.status", http.StatusOK),
		slog.String("book.id", bookID),
	)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
