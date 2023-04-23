package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type operationRoutes struct {
	service.Operation
}

func newOperationRoutes(g *echo.Group, operationService service.Operation) *operationRoutes {
	r := &operationRoutes{
		Operation: operationService,
	}

	g.GET("/history", r.getHistory)
	g.GET("/report-link", r.getReportLink)
	g.GET("/report-file", r.getReportFile)

	return r
}

type getHistoryInput struct {
	AccountId int    `json:"account_id" validate:"required"`
	SortType  string `json:"sort_type,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// @Summary Get history
// @Description Get history of operations
// @Tags operations
// @Accept json
// @Produce json
// @Param input body getHistoryInput true "input"
// @Success 200 {object} v1.operationRoutes.getHistory.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/operations/history [get]
func (r *operationRoutes) getHistory(c echo.Context) error {
	var input getHistoryInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	operations, err := r.Operation.OperationHistory(c.Request().Context(), service.OperationHistoryInput{
		AccountId: input.AccountId,
		SortType:  input.SortType,
		Offset:    input.Offset,
		Limit:     input.Limit,
	})
	if err != nil {
		log.Debugf("error while getting operation history: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Operations []service.OperationHistoryOutput `json:"operations"`
	}

	return c.JSON(http.StatusOK, response{
		Operations: operations,
	})
}

type getReportInput struct {
	Month int `json:"month" validate:"required"`
	Year  int `json:"year" validate:"required"`
}

// @Summary Get report link
// @Description Get link to report
// @Tags operations
// @Accept json
// @Produce json
// @Param input body getReportInput true "input"
// @Success 200 {object} v1.operationRoutes.getReportLink.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/operations/report-link [get]
func (r *operationRoutes) getReportLink(c echo.Context) error {
	var input getReportInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	link, err := r.Operation.MakeReportLink(c.Request().Context(), input.Month, input.Year)
	if err != nil {
		log.Debugf("error while getting report link: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Link string `json:"link"`
	}

	return c.JSON(http.StatusOK, response{
		Link: link,
	})
}

// @Summary Get report file
// @Description Get report file
// @Tags operations
// @Accept json
// @Produce text/csv
// @Param input body getReportInput true "input"
// @Success 200
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/operations/report-file [get]
func (r *operationRoutes) getReportFile(c echo.Context) error {
	var input getReportInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	file, err := r.Operation.MakeReportFile(c.Request().Context(), input.Month, input.Year)
	if err != nil {
		log.Debugf("error while getting report file: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.Blob(http.StatusOK, "text/csv", file)
}
