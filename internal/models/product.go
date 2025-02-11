package models

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/lib/pq"
)

type Product struct {
	ID              int               `json:"id,omitempty"`
	Name            string            `json:"name,omitempty" gorm:"column:name;type:varchar(255)" validate:"required"`
	Description     string            `json:"description,omitempty" gorm:"column:description;type:text" validate:"required"`
	Price           float64           `json:"price,omitempty" gorm:"column:price;type:decimal(10,2)" validate:"required"`
	Categories      pq.Int64Array     `json:"categories,omitempty" gorm:"column:categories;type:int[]" validate:"required"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	ProductVariants []ProductVariants `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
}

func (*Product) TableName() string {
	return "products"
}

func (l Product) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type ProductCategory struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty" gorm:"column:name;type:varchar(255)" validate:"required"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (*ProductCategory) TableName() string {
	return "product_categories"
}

func (l ProductCategory) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type ProductVariants struct {
	ID        int       `json:"id,omitempty"`
	ProductID int       `json:"product_id,omitempty" gorm:"column:product_id;type:int" validate:"required"`
	Color     string    `json:"color,omitempty" gorm:"column:color;type:varchar(50)" validate:"required"`
	Size      string    `json:"size,omitempty" gorm:"column:size;type:varchar(50)" validate:"required"`
	Quantity  int       `json:"quantity,omitempty" gorm:"column:quantity"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (*ProductVariants) TableName() string {
	return "product_variants"
}

func (l ProductVariants) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
