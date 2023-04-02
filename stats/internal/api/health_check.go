package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radovsky1/internship_backend_2022/stats/internal/api/v1/formatter"
	"net/http"
)

type healthCheckHandler struct {
	db *pgxpool.Pool
}

func newHealthCheckHandler(db *pgxpool.Pool) *healthCheckHandler {
	return &healthCheckHandler{
		db: db,
	}
}

func (h *healthCheckHandler) Health(c *gin.Context) {
	if err := h.db.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, formatter.NewError("error ping DB", http.StatusInternalServerError))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (*healthCheckHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
