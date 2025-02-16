package repository

import (
	"context"
	"ecommerce-product/constants"
	"ecommerce-product/helpers"
	"ecommerce-product/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductRepo struct {
	DB    *gorm.DB
	Redis *redis.ClusterClient
}

func (r *ProductRepo) InsertNewProduct(ctx context.Context, product *models.Product) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(product).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		go func() {
			ctx := context.Background()
			jsonData, err := json.Marshal(product)
			if err != nil {
				helpers.Logger.Warn("failed to marshal the product for cache: ", err)
				return
			}

			if err := r.Redis.Del(ctx, constants.RedisKeyProducts).Err(); err != nil {
				helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
			}
			if err := r.Redis.Set(ctx, fmt.Sprintf(constants.RedisKeyProductsDetail, product.ID), string(jsonData), time.Hour*24).Err(); err != nil {
				helpers.Logger.Warn("failed to insert redis with key: ", fmt.Sprintf(constants.RedisKeyProductsDetail, product.ID), err)
			}
		}()
	}

	return err
}

func (r ProductRepo) UpdateProduct(ctx context.Context, productID int, newData map[string]interface{}) error {
	err := r.DB.Model(&models.Product{}).Where("id = ?", productID).Updates(newData).Error
	if err != nil {
		return err
	}

	go func() {
		ctx := context.Background()

		if err := r.Redis.Del(ctx, constants.RedisKeyProducts).Err(); err != nil {
			helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
		}
		if err := r.Redis.Del(ctx, fmt.Sprintf(constants.RedisKeyProductsDetail, productID)).Err(); err != nil {
			helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
		}
	}()

	return nil
}

func (r ProductRepo) UpdateProductVariant(ctx context.Context, variantID int, newData map[string]interface{}) error {
	err := r.DB.Model(&models.ProductVariants{}).Where("id = ?", variantID).Updates(newData).Error
	if err != nil {
		return err
	}

	go func() {
		ctx := context.Background()

		variant := &models.ProductVariants{}
		err := r.DB.Where("id = ?", variantID).Last(variant).Error
		if err != nil {
			helpers.Logger.Warn("failed to get product: ", err)
			return
		}

		if err := r.Redis.Del(ctx, constants.RedisKeyProducts).Err(); err != nil {
			helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
		}
		if err := r.Redis.Del(ctx, fmt.Sprintf(constants.RedisKeyProductsDetail, variant.ProductID)).Err(); err != nil {
			helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
		}
	}()

	return nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, productID int) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Exec("DELETE FROM products WHERE id = ?", productID).Error
		if err != nil {
			return err
		}
		err = tx.Exec("DELETE FROM product_variants WHERE product_id = ?", productID).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		go func() {
			ctx := context.Background()

			if err := r.Redis.Del(ctx, constants.RedisKeyProducts).Err(); err != nil {
				helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
			}
			if err := r.Redis.Del(ctx, fmt.Sprintf(constants.RedisKeyProductsDetail, productID)).Err(); err != nil {
				helpers.Logger.Warn("failed to delete redis with key: ", constants.RedisKeyProducts, err)
			}

		}()

	}

	return err
}

func (r *ProductRepo) GetProducts(ctx context.Context, page int, limit int) ([]models.Product, error) {
	offset := (page - 1) * limit
	products := []models.Product{}

	// Get redis first
	productStr, err := r.Redis.Get(ctx, constants.RedisKeyProducts).Result()
	if err == nil && productStr != "" {
		result := []models.Product{}
		if err := json.Unmarshal([]byte(productStr), &products); err != nil {
			helpers.Logger.Warn("failed to unmarshal the redis response get products")
		}

		if page > 0 && limit > 0 {
			for i := offset; i < len(products); i++ {
				if i == offset+limit {
					break
				}
				result = append(result, products[i])
			}
		} else {
			helpers.Logger.Info("successfuly get product from redis cache")
			return products, nil
		}

		helpers.Logger.Info("successfuly get product from redis cache")
		return result, nil
	}

	sql := r.DB.Preload("ProductVariants")
	if limit > 0 && page > 0 {
		sql = sql.Limit(limit).Offset(offset)
	}
	err = sql.Find(&products).Error
	if err != nil {
		return nil, err
	}

	go func() {
		ctx := context.Background()

		cacheProduct := []models.Product{}
		err := sql.Find(&cacheProduct).Error
		if err != nil {
			helpers.Logger.Warn("failed to get product to cache", err)
			return
		}

		jsonCacheProduct, err := json.Marshal(cacheProduct)
		if err != nil {
			helpers.Logger.Warn("failed to marshal the product to cache", err)
			return
		}
		err = r.Redis.Set(ctx, constants.RedisKeyProducts, string(jsonCacheProduct), time.Hour*24).Err()
		if err != nil {
			helpers.Logger.Warn("failed to set product to cache", err)
			return
		}

	}()

	return products, nil
}

func (r *ProductRepo) GetProductDetail(ctx context.Context, productID int) (models.Product, error) {

	product := models.Product{}

	// Get redis first
	productStr, err := r.Redis.Get(ctx, fmt.Sprintf(constants.RedisKeyProductsDetail, productID)).Result()
	if err == nil && productStr != "" {
		if err := json.Unmarshal([]byte(productStr), &product); err != nil {
			helpers.Logger.Warn("failed to unmarshal the redis response get products")
		}

		helpers.Logger.Info("successfuly get product detail from redis cache")
		return product, nil
	}

	err = r.DB.Preload("ProductVariants").Where("id = ?", productID).First(&product).Error
	if err != nil {
		return models.Product{}, err
	}

	go func() {
		ctx := context.Background()

		jsonCacheProduct, err := json.Marshal(product)
		if err != nil {
			helpers.Logger.Warn("failed to marshal the product detail to cache", err)
			return
		}
		err = r.Redis.Set(ctx, fmt.Sprintf(constants.RedisKeyProductsDetail, productID), string(jsonCacheProduct), time.Hour*24).Err()
		if err != nil {
			helpers.Logger.Warn("failed to set product detail to cache", err)
			return
		}

	}()

	return product, nil
}
