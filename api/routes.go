package api

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine, h *Handler) *gin.Engine {

	r.POST("/login", loginHandler)

	protected := r.Group("/", authorizationMiddleware)

	protected.GET("/books", h.listBooksHandler)
	protected.POST("/book", h.createBookHandler)
	protected.DELETE("/books/:id", h.deleteBookHandler)

	return r

}
