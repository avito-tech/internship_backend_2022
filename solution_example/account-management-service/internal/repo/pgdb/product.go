package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/pkg/postgres"
	"context"
	"fmt"
)

type ProductRepo struct {
	*postgres.Postgres
}

func NewProductRepo(pg *postgres.Postgres) *ProductRepo {
	return &ProductRepo{pg}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, name string) (int, error) {
	sql, args, _ := r.Builder.
		Insert("products").
		Columns("name").
		Values(name).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ProductRepo.CreateProduct - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *ProductRepo) GetProductById(ctx context.Context, id int) (entity.Product, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("products").
		Where("id = ?", id).
		ToSql()

	var product entity.Product
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&product.Id,
		&product.Name,
	)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepo.GetProductById - r.Pool.QueryRow: %v", err)
	}

	return product, nil
}

func (r *ProductRepo) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("products").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo.GetAllProducts - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("ProductRepo.GetAllProducts - rows.Scan: %v", err)
		}
		products = append(products, product)
	}

	return products, nil
}
