package cmd

import (
	"ecommerce-product/external"
	"ecommerce-product/helpers"
	"ecommerce-product/internal/api"
	"ecommerce-product/internal/interfaces"
	"ecommerce-product/internal/repository"
	"ecommerce-product/internal/services"

	"github.com/labstack/echo/v4"
)

func ServeHTTP() {
	d := dependencyInject()

	e := echo.New()
	e.GET("/healthcheck", d.HealthcheckAPI.Healthcheck)

	productV1 := e.Group("/products/v1")
	productV1.POST("", d.ProductAPI.CreateProduct, d.MiddlewareValidateAuth)
	productV1.PUT("/:id", d.ProductAPI.UpdateProduct, d.MiddlewareValidateAuth)
	productV1.PUT("/variant/:id", d.ProductAPI.UpdateProductVariant, d.MiddlewareValidateAuth)
	productV1.DELETE("/:id", d.ProductAPI.DeleteProduct, d.MiddlewareValidateAuth)
	productV1.GET("/list", d.ProductAPI.GetProducts)
	productV1.GET("/:id", d.ProductAPI.GetProductDetail)

	categoryV1 := e.Group("products/v1/category")
	categoryV1.POST("", d.CategoryAPI.CreateCategory, d.MiddlewareValidateAuth)
	categoryV1.PUT("/:id", d.CategoryAPI.UpdateProductCategory, d.MiddlewareValidateAuth)
	categoryV1.DELETE("/:id", d.CategoryAPI.DeleteCategory, d.MiddlewareValidateAuth)
	categoryV1.GET("", d.CategoryAPI.GetCategories)

	e.Start(":" + helpers.GetEnv("PORT", "9000"))
}

type Dependency struct {
	External       interfaces.IExternal
	HealthcheckAPI *api.HealthcheckAPI

	ProductAPI  interfaces.IProductAPI
	CategoryAPI interfaces.ICategoryAPI
}

func dependencyInject() Dependency {
	external := &external.External{}

	productRepo := &repository.ProductRepo{
		DB:    helpers.DB,
		Redis: helpers.RedisClient,
	}

	categoryRepo := &repository.CategoryRepo{
		DB: helpers.DB,
	}

	productSvc := &services.ProductService{
		ProductRepo: productRepo,
	}

	categorySvc := &services.CategoryService{
		CategoryRepo: categoryRepo,
	}

	productAPI := &api.ProductAPI{
		ProductService: productSvc,
	}

	categoryAPI := &api.CategoryAPI{
		CategoryService: categorySvc,
	}

	return Dependency{
		External:       external,
		HealthcheckAPI: &api.HealthcheckAPI{},
		ProductAPI:     productAPI,
		CategoryAPI:    categoryAPI,
	}
}
