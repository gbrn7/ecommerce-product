package interfaces

import (
	"context"
	"ecommerce-product/internal/models"

	"github.com/labstack/echo/v4"
)

type IProductRepo interface {
	InsertNewProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, productID int, newData map[string]interface{}) error
	UpdateProductVariant(ctx context.Context, variantID int, newData map[string]interface{}) error
	DeleteProduct(ctx context.Context, productID int) error
	GetProducts(ctx context.Context, page int, limit int) ([]models.Product, error)
	GetProductDetail(ctx context.Context, productID int) (models.Product, error)
}

type IProductService interface {
	CreateProduct(ctx context.Context, req *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, productID int, req models.Product) error
	UpdateProductVariant(ctx context.Context, variantID int, req models.ProductVariants) error
	DeleteProduct(ctx context.Context, productID int) error
	GetProducts(ctx context.Context, page int, limit int) ([]models.Product, error)
	GetProductDetail(ctx context.Context, productID int) (models.Product, error)
}

type IProductAPI interface {
	CreateProduct(e echo.Context) error
	UpdateProduct(e echo.Context) error
	UpdateProductVariant(e echo.Context) error
	DeleteProduct(e echo.Context) error
	GetProducts(e echo.Context) error
	GetProductDetail(e echo.Context) error
}
