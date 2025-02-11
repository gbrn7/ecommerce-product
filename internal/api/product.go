package api

import (
	"ecommerce-product/constants"
	"ecommerce-product/helpers"
	"ecommerce-product/internal/interfaces"
	"ecommerce-product/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductAPI struct {
	ProductService interfaces.IProductService
}

func (api *ProductAPI) CreateProduct(e echo.Context) error {
	var (
		log = helpers.Logger
	)
	req := models.Product{}
	if err := e.Bind(&req); err != nil {
		log.Error("failed to parse request, ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request, ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	resp, err := api.ProductService.CreateProduct(e.Request().Context(), &req)

	if err != nil {
		log.Error("failed to create product, ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, resp)

}

func (api *ProductAPI) UpdateProduct(e echo.Context) error {
	var (
		log          = helpers.Logger
		productIDstr = e.Param("id")
	)

	productID, err := strconv.Atoi(productIDstr)

	fmt.Println("productIDstr")
	fmt.Println(productIDstr)
	if err != nil || productID == 0 {
		log.Error("failed to get productID : ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	req := models.Product{}
	if err := e.Bind(&req); err != nil {
		log.Error("failed to parse request, ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	err = api.ProductService.UpdateProduct(e.Request().Context(), productID, req)

	if err != nil {
		log.Error("failed to update product, ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, nil)

}

func (api *ProductAPI) UpdateProductVariant(e echo.Context) error {
	var (
		log          = helpers.Logger
		variantIDstr = e.Param("id")
	)

	variantID, err := strconv.Atoi(variantIDstr)
	if err != nil || variantID == 0 {
		log.Error("failed to get variantID", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	req := models.ProductVariants{}
	if err := e.Bind(&req); err != nil {
		log.Error("failed to parse request, ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	err = api.ProductService.UpdateProductVariant(e.Request().Context(), variantID, req)

	if err != nil {
		log.Error("failed to update product variant, ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, nil)

}

func (api *ProductAPI) DeleteProduct(e echo.Context) error {
	var (
		log          = helpers.Logger
		productIDstr = e.Param("id")
	)

	productID, err := strconv.Atoi(productIDstr)
	if err != nil || productID == 0 {
		log.Error("failed to get productID", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	err = api.ProductService.DeleteProduct(e.Request().Context(), productID)

	if err != nil {
		log.Error("failed to delete product, ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, nil)

}

func (api *ProductAPI) GetProducts(e echo.Context) error {
	var (
		log      = helpers.Logger
		pageStr  = e.QueryParam("page")
		limitStr = e.QueryParam("limit")
	)

	page, err := strconv.Atoi(pageStr)
	if err != nil && pageStr != "" {
		log.Error("failed to get page", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		log.Error("failed to get limit", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	resp, err := api.ProductService.GetProducts(e.Request().Context(), page, limit)

	if err != nil {
		log.Error("failed to get products, ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, resp)

}

func (api *ProductAPI) GetProductDetail(e echo.Context) error {
	var (
		log          = helpers.Logger
		productIDstr = e.Param("id")
	)

	productID, err := strconv.Atoi(productIDstr)
	if err != nil || productID == 0 {
		log.Error("failed to get productID", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	resp, err := api.ProductService.GetProductDetail(e.Request().Context(), productID)

	if err != nil {
		log.Error("failed to get product detail, ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, resp)
}
