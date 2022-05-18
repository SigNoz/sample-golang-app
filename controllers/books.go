package controllers

import (
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
	var books []models.Book
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("This is a sample event", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("sampleAttribute", "Test")))
	models.DB.WithContext(c.Request.Context()).Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.WithContext(c.Request.Context()).Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {
	// Validate input
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.WithContext(c.Request.Context()).Create(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.WithContext(c.Request.Context()).Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.WithContext(c.Request.Context()).Model(&book).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.WithContext(c.Request.Context()).Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
