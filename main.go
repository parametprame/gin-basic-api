package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var err error

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{})

	handler := newHandler(db)

	r := gin.New()

	r.POST("/login", loginHandler)

	protected := r.Group("/", authorizationMiddleware)

	protected.GET("/books", handler.listBooksHandler)
	protected.POST("/book", handler.createBookHandler)
	protected.DELETE("/books/:id", handler.deleteBookHandler)

	r.Run()
}

type Handler struct {
	db *gorm.DB
}

func newHandler(db *gorm.DB) *Handler {
	return &Handler{db}
}

func validateToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("MySignature"), nil
	})
	return err
}

func authorizationMiddleware(ctx *gin.Context) {
	s := ctx.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	if err := validateToken(token); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (h *Handler) listBooksHandler(ctx *gin.Context) {
	books := []Book{}
	s := ctx.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	if err := validateToken(token); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please Authentication",
		})
		return
	}

	if result := h.db.Find(&books); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &books)
}

func (h *Handler) createBookHandler(ctx *gin.Context) {
	book := Book{}

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := h.db.Create(&book); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, &book)
}

func (h *Handler) deleteBookHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if result := h.db.Delete(&Book{}, id); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})

		return
	}

	ctx.Status(http.StatusNoContent)
}

func loginHandler(c *gin.Context) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	})

	ss, err := token.SignedString([]byte("MySignature"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}
