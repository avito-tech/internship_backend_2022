package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/stats/internal/domain"
	"net/http"
)

type getReportRequest struct {
	ServiceID string `form:"service_id" binding:"required"`
	Month     int    `form:"month" binding:"required"`
	Year      int    `form:"year" binding:"required"`
}

type getReportResponse struct {
	Reports []*domain.ReportByServiceID `json:"reports"`
}

func (rH RouterHandler) getReport(c *gin.Context) {
	var req getReportRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reports, err := rH.ucService.GetReportByServiceID(c, serviceID, req.Month, req.Year)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, getReportResponse{
		Reports: reports,
	})
}
