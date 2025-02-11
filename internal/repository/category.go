package repository

import (
	"context"
	"ecommerce-product/internal/models"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	DB *gorm.DB
}

func (r *CategoryRepo) InsertNewCategory(ctx context.Context, category *models.ProductCategory) error {
	return r.DB.Create(category).Error
}

func (r CategoryRepo) UpdateCategory(ctx context.Context, categoryID int, newData map[string]interface{}) error {
	err := r.DB.Model(&models.ProductCategory{}).Where("id = ?", categoryID).Updates(newData).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, categoryID int) error {
	return r.DB.Exec("DELETE FROM product_categories WHERE id = ?", categoryID).Error
}

func (r *CategoryRepo) GetCategories(ctx context.Context) ([]models.ProductCategory, error) {
	var (
		resp []models.ProductCategory
		err  error
	)

	err = r.DB.Find(&resp).Error
	return resp, err
}
