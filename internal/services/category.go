package services

import (
	"context"
	"ecommerce-product/internal/interfaces"
	"ecommerce-product/internal/models"
	"encoding/json"

	"github.com/pkg/errors"
)

type CategoryService struct {
	CategoryRepo interfaces.ICategoryRepo
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *models.ProductCategory) (*models.ProductCategory, error) {
	err := s.CategoryRepo.InsertNewCategory(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert new category")
	}

	resp := req
	return resp, nil
}

func (s *CategoryService) UpdateProductCategory(ctx context.Context, categoryID int, req models.ProductCategory) error {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request")
	}

	newData := map[string]interface{}{}
	err = json.Unmarshal(jsonReq, &newData)

	if err != nil {
		return errors.Wrap(err, "failed to unmarshaling to map")
	}

	err = s.CategoryRepo.UpdateCategory(ctx, categoryID, newData)
	if err != nil {
		return errors.Wrap(err, "failed to update product category")
	}

	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, categoryID int) error {
	return s.CategoryRepo.DeleteCategory(ctx, categoryID)
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]models.ProductCategory, error) {
	return s.CategoryRepo.GetCategories(ctx)
}
