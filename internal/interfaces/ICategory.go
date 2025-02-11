package interfaces

import (
	"context"
	"ecommerce-product/internal/models"

	"github.com/labstack/echo/v4"
)

type ICategoryRepo interface {
	InsertNewCategory(ctx context.Context, category *models.ProductCategory) error
	UpdateCategory(ctx context.Context, categoryID int, newData map[string]interface{}) error
	DeleteCategory(ctx context.Context, categoryID int) error
	GetCategories(ctx context.Context) ([]models.ProductCategory, error)
}

type ICategoryService interface {
	CreateCategory(ctx context.Context, req *models.ProductCategory) (*models.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, categoryID int, req models.ProductCategory) error
	DeleteCategory(ctx context.Context, categoryID int) error
	GetCategories(ctx context.Context) ([]models.ProductCategory, error)
}

type ICategoryAPI interface {
	CreateCategory(e echo.Context) error
	UpdateProductCategory(e echo.Context) error
	DeleteCategory(e echo.Context) error
	GetCategories(e echo.Context) error
}
