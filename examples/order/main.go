package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/sklyar/go-transact"
	"github.com/sklyar/go-transact/adapters/txstd"
	"github.com/sklyar/go-transact/txsql"
)

type OrderRepository interface {
	Create(ctx context.Context, customerID int) (int, error)
	AddProduct(ctx context.Context, orderID, productID int) error
}

type InventoryRepository interface {
	GetProductQuantity(ctx context.Context, productID int) (int, error)
	DecrementProductQuantity(ctx context.Context, productID int) error
}

type OrderService struct {
	orderRepo     OrderRepository
	inventoryRepo InventoryRepository

	txManager *transact.Manager
}

func NewOrderService(
	txManager *transact.Manager,
	orderRepo OrderRepository,
	inventoryRepo InventoryRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		txManager:     txManager,
	}
}

func main() {
	ctx := context.Background()

	sqlDB, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("unable to open database: %w", err))
	}

	txManager, db, err := transact.NewManager(txstd.Wrap(sqlDB))
	if err != nil {
		panic(fmt.Errorf("unable to create transaction manager: %w", err))
	}

	orderRepo := &orderRepository{db: db}
	inventoryRepo := &inventoryRepository{db: db}

	orderService := NewOrderService(txManager, orderRepo, inventoryRepo)

	// sample customer and products
	customerID := 1
	products := []int{1, 2, 3}

	orderID, err := orderService.Create(ctx, customerID, products)
	if err != nil {
		fmt.Printf("Failed to create order: %v", err)
		return
	}

	fmt.Printf("Order '%d' created successfully", orderID)
}

// Create creates a new order and adds products to it
func (s *OrderService) Create(ctx context.Context, customerID int, products []int) (int, error) {
	ctx, tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderID, err := s.orderRepo.Create(ctx, customerID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	for _, productID := range products {
		quantity, err := s.inventoryRepo.GetProductQuantity(ctx, productID)
		if err != nil {
			return 0, fmt.Errorf("failed to get product quantity: %w", err)
		}
		if quantity < 1 {
			return 0, fmt.Errorf("not enough quantity for product %d", productID)
		}

		if err := s.orderRepo.AddProduct(ctx, orderID, productID); err != nil {
			return 0, fmt.Errorf("failed to add product to order: %w", err)
		}

		if err := s.inventoryRepo.DecrementProductQuantity(ctx, productID); err != nil {
			return 0, fmt.Errorf("failed to update inventory: %w", err)
		}
	}

	if _, err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit: %w", err)
	}

	return orderID, nil
}

type orderRepository struct {
	db txsql.DB
}

// Create creates a new order and returns its ID
func (s *orderRepository) Create(ctx context.Context, customerID int) (int, error) {
	var orderID int
	query := "INSERT INTO orders (customer_id) VALUES ($1) RETURNING id"
	row := s.db.QueryRow(ctx, query, customerID)
	err := row.Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// AddProduct adds a product to the order
func (s *orderRepository) AddProduct(ctx context.Context, orderID, productID int) error {
	query := "INSERT INTO order_products (order_id, product_id) VALUES ($1, $2)"
	_, err := s.db.Exec(ctx, query, orderID, productID)
	if err != nil {
		return err
	}

	return nil
}

type inventoryRepository struct {
	db txsql.DB
}

func (r *inventoryRepository) GetProductQuantity(ctx context.Context, productID int) (int, error) {
	var quantity int
	query := "SELECT quantity FROM inventory WHERE product_id = $1"
	row := r.db.QueryRow(ctx, query, productID)
	err := row.Scan(&quantity)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func (r *inventoryRepository) DecrementProductQuantity(ctx context.Context, productID int) error {
	query := "UPDATE inventory SET quantity = quantity - 1 WHERE product_id = $1"
	_, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return err
	}

	return nil
}
