package v1

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type productRoutes struct {
	productService service.Product
}

func newProductRoutes(g *echo.Group, productService service.Product) *productRoutes {
	r := &productRoutes{
		productService: productService,
	}

	g.POST("/create", r.create)
	g.GET("/", r.getById)

	return r
}

type productCreateInput struct {
	Name string `json:"name" validate:"required"`
}

// @Summary Create product
// @Description Create product
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} v1.productRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/products/create [post]
func (r *productRoutes) create(c echo.Context) error {
	var input productCreateInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.productService.CreateProduct(c.Request().Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

type getByIdInput struct {
	Id int `json:"id" validate:"required"`
}

// @Summary Get product by id
// @Description Get product by id
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} v1.productRoutes.getById.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/products/getById [get]
func (r *productRoutes) getById(c echo.Context) error {
	var input getByIdInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	product, err := r.productService.GetProductById(c.Request().Context(), input.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Product entity.Product `json:"product"`
	}

	return c.JSON(http.StatusOK, response{
		Product: product,
	})
}
