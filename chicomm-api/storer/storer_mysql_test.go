package storer

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	fn(db, mock)
}

func TestMySQLStorer_CreateProduct(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))

				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WillReturnError(fmt.Errorf("failed inserting product"))

				_, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting last insert ID",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("failed to get the last insert ID")))

				_, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_GetProduct(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)

				gp, err := st.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("failed getting product"))

				_, err := st.GetProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_ListProducts(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM products").WillReturnRows(rows)

				products, err := st.ListProducts(context.Background())
				require.NoError(t, err)
				require.Len(t, products, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products").WillReturnError(fmt.Errorf("failed querying products"))

				_, err := st.ListProducts(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_UpdateProduct(t *testing.T) {
	p := &Product{
		ID:           1,
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.0,
		CountInStock: 100,
	}

	np := &Product{
		ID:           1,
		Name:         "new test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))

				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)

				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=? WHERE id=?").
					WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := st.UpdateProduct(context.Background(), np)
				require.NoError(t, err)
				require.Equal(t, int64(1), up.ID)
				require.Equal(t, np.Name, up.Name)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=? WHERE id=?").
					WillReturnError(fmt.Errorf("failed updating product"))

				_, err := st.UpdateProduct(context.Background(), np)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_DeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := st.DeleteProduct(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(1).
					WillReturnError(fmt.Errorf("error deleting product"))

				err := st.DeleteProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_CreateOrder(t *testing.T) {
	ois := []OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test.jpg",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				co, err := st.CreateOrder(context.Background(), o)
				require.NoError(t, err)
				require.Equal(t, int64(1), co.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").
					WillReturnError(fmt.Errorf("failed creating order"))
				mock.ExpectRollback()

				_, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order item",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
					WillReturnError(fmt.Errorf("failed creating order item"))
				mock.ExpectRollback()

				_, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed committing transaction",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("failed committing transaction"))

				_, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_GetOrder(t *testing.T) {
	ois := []OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test.jpg",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").
					WithArgs(1).
					WillReturnRows(orows)

				oirows := sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(1, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, 1).
					AddRow(2, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, 1)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").
					WithArgs(1).
					WillReturnRows(oirows)

				mo, err := st.GetOrder(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), mo.ID)

				for i, oi := range mo.Items {
					require.Equal(t, ois[i].Name, oi.Name)
					require.Equal(t, ois[i].Quantity, oi.Quantity)
					require.Equal(t, ois[i].Image, oi.Image)
					require.Equal(t, ois[i].Price, oi.Price)
					require.Equal(t, ois[i].ProductID, oi.ProductID)
				}

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").
					WithArgs(1).
					WillReturnError(fmt.Errorf("failed getting order"))

				_, err := st.GetOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting order item",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").
					WithArgs(1).
					WillReturnRows(orows)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").
					WithArgs(1).
					WillReturnError(fmt.Errorf("failed getting order item"))

				_, err := st.GetOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_ListOrders(t *testing.T) {
	ois := []OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test.jpg",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders").
					WillReturnRows(orows)

				oirows := sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(1, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, 1).
					AddRow(2, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, 1)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").
					WithArgs(1).
					WillReturnRows(oirows)

				mo, err := st.ListOrders(context.Background())
				require.NoError(t, err)
				require.Len(t, mo, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying orders",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM orders").
					WillReturnError(fmt.Errorf("failed querying orders"))

				_, err := st.ListOrders(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying order items",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM orders").
					WillReturnRows(orows)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").
					WithArgs(1).
					WillReturnError(fmt.Errorf("failed querying order items"))

				_, err := st.ListOrders(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestMySQLStorer_DeleteOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM orders WHERE id=?").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := st.DeleteOrder(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order item",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").
					WillReturnError(fmt.Errorf("failed deleting order item"))
				mock.ExpectRollback()

				err := st.DeleteOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM orders WHERE id=?").
					WillReturnError(fmt.Errorf("failed deleting order"))
				mock.ExpectRollback()

				err := st.DeleteOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}
