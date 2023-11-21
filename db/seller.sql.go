// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: seller.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteCoupon = `-- name: DeleteCoupon :exec

DELETE FROM "coupon" c
WHERE c."id" = $2 AND "shop_id" = (
        SELECT s."id"
        FROM "shop" s
        WHERE
            s."seller_name" = $1
            AND s."enabled" = true
    )
`

type DeleteCouponParams struct {
	SellerName string `json:"seller_name" params:"seller_name"`
	ID         int32  `json:"id" params:"id"`
}

func (q *Queries) DeleteCoupon(ctx context.Context, arg DeleteCouponParams) error {
	_, err := q.db.Exec(ctx, deleteCoupon, arg.SellerName, arg.ID)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec

UPDATE "product" p
SET
    "enabled" = false,
    "edit_date" = NOW()
WHERE "shop_id" = (
        SELECT s."id"
        FROM "shop" s
        WHERE
            s."seller_name" = $1
            AND s."enabled" = true
    )
    AND p."id" = $2
`

type DeleteProductParams struct {
	SellerName string `json:"seller_name" params:"seller_name"`
	ID         int32  `json:"id" params:"id"`
}

func (q *Queries) DeleteProduct(ctx context.Context, arg DeleteProductParams) error {
	_, err := q.db.Exec(ctx, deleteProduct, arg.SellerName, arg.ID)
	return err
}

const getSellerInfo = `-- name: GetSellerInfo :one

SELECT s.id, s.seller_name, s.image_id, s.name, s.description, s.enabled
FROM "user" u
    JOIN "shop" s ON u.username = s.seller_name
WHERE u.id = $1
`

func (q *Queries) GetSellerInfo(ctx context.Context, id int32) (Shop, error) {
	row := q.db.QueryRow(ctx, getSellerInfo, id)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.SellerName,
		&i.ImageID,
		&i.Name,
		&i.Description,
		&i.Enabled,
	)
	return i, err
}

const insertTag = `-- name: InsertTag :one

INSERT INTO
    "tag" ("shop_id", "name")
VALUES ( (
            SELECT s."id"
            FROM "shop" s
            WHERE
                s."seller_name" = $1
                AND s."enabled" = true
        ),
        $2
    ) RETURNING ("id", "name")
`

type InsertTagParams struct {
	SellerName string `json:"seller_name" params:"seller_name"`
	Name       string `json:"name"`
}

func (q *Queries) InsertTag(ctx context.Context, arg InsertTagParams) (interface{}, error) {
	row := q.db.QueryRow(ctx, insertTag, arg.SellerName, arg.Name)
	var column_1 interface{}
	err := row.Scan(&column_1)
	return column_1, err
}

const orderDetail = `-- name: OrderDetail :one

SELECT order_id, product_id, product_version, quantity, id, version, name, description, price, image_id
FROM "order_detail"
    LEFT JOIN "product_archive" ON "order_detail"."product_id" = "product"."id" AND "order_detail"."version" = "product"."version"
WHERE "order_id" = $1
`

type OrderDetailRow struct {
	OrderID        int32          `json:"order_id"`
	ProductID      int32          `json:"product_id"`
	ProductVersion int32          `json:"product_version"`
	Quantity       int32          `json:"quantity"`
	ID             pgtype.Int4    `json:"id"`
	Version        pgtype.Int4    `json:"version"`
	Name           pgtype.Text    `json:"name"`
	Description    pgtype.Text    `json:"description"`
	Price          pgtype.Numeric `json:"price"`
	ImageID        pgtype.UUID    `json:"image_id"`
}

func (q *Queries) OrderDetail(ctx context.Context, orderID int32) (OrderDetailRow, error) {
	row := q.db.QueryRow(ctx, orderDetail, orderID)
	var i OrderDetailRow
	err := row.Scan(
		&i.OrderID,
		&i.ProductID,
		&i.ProductVersion,
		&i.Quantity,
		&i.ID,
		&i.Version,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.ImageID,
	)
	return i, err
}

const searchTag = `-- name: SearchTag :many

SELECT t."id", t."name"
FROM "tag" t
    JOIN "shop" s ON "shop_id" = s.id
    JOIN "user" u ON s.seller_name = u.username
WHERE u.id = $1 AND t."name" ~* $2
ORDER BY LENGTH(t."name")
LIMIT 10
`

type SearchTagParams struct {
	ID   int32  `json:"id" params:"id"`
	Name string `json:"name"`
}

type SearchTagRow struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) SearchTag(ctx context.Context, arg SearchTagParams) ([]SearchTagRow, error) {
	rows, err := q.db.Query(ctx, searchTag, arg.ID, arg.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchTagRow
	for rows.Next() {
		var i SearchTagRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const sellerGetCoupon = `-- name: SellerGetCoupon :many

SELECT
    c."id",
    c."type",
    c."shop_id",
    c."name",
    c."discount",
    c."expire_date"
FROM "coupon" c
    JOIN "shop" s ON c."shop_id" = s.id
WHERE s.seller_name = $1
ORDER BY "start_date" DESC
LIMIT $2
OFFSET $3
`

type SellerGetCouponParams struct {
	SellerName string `json:"seller_name" params:"seller_name"`
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
}

type SellerGetCouponRow struct {
	ID         int32              `json:"id" params:"id"`
	Type       CouponType         `json:"type"`
	ShopID     int32              `json:"shop_id"`
	Name       string             `json:"name"`
	Discount   pgtype.Numeric     `json:"discount"`
	ExpireDate pgtype.Timestamptz `json:"expire_date"`
}

func (q *Queries) SellerGetCoupon(ctx context.Context, arg SellerGetCouponParams) ([]SellerGetCouponRow, error) {
	rows, err := q.db.Query(ctx, sellerGetCoupon, arg.SellerName, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SellerGetCouponRow
	for rows.Next() {
		var i SellerGetCouponRow
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.ShopID,
			&i.Name,
			&i.Discount,
			&i.ExpireDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const sellerGetCouponDetail = `-- name: SellerGetCouponDetail :one

SELECT c.id, c.type, c.shop_id, c.name, c.description, c.discount, c.start_date, c.expire_date
FROM "coupon" c
    JOIN "shop" s ON c."shop_id" = s.id
WHERE
    c."id" = $1
    and s."seller_name" = $2
`

type SellerGetCouponDetailParams struct {
	ID         int32  `json:"id" params:"id"`
	SellerName string `json:"seller_name" params:"seller_name"`
}

func (q *Queries) SellerGetCouponDetail(ctx context.Context, arg SellerGetCouponDetailParams) (Coupon, error) {
	row := q.db.QueryRow(ctx, sellerGetCouponDetail, arg.ID, arg.SellerName)
	var i Coupon
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.ShopID,
		&i.Name,
		&i.Description,
		&i.Discount,
		&i.StartDate,
		&i.ExpireDate,
	)
	return i, err
}

const sellerGetOrder = `-- name: SellerGetOrder :many

SELECT
    "id",
    "shop_id",
    "shipment",
    "total_price",
    "status",
    "created_at"
FROM "order_history"
WHERE "shop_id" = (
        SELECT s."id"
        FROM "shop" s
        WHERE
            s."seller_name" = $1
            AND s."enabled" = true
    )
ORDER BY "created_at" DESC
LIMIT $2
OFFSET $3
`

type SellerGetOrderParams struct {
	SellerName string `json:"seller_name" params:"seller_name"`
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
}

type SellerGetOrderRow struct {
	ID         int32              `json:"id"`
	ShopID     int32              `json:"shop_id"`
	Shipment   int32              `json:"shipment"`
	TotalPrice int32              `json:"total_price"`
	Status     OrderStatus        `json:"status"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) SellerGetOrder(ctx context.Context, arg SellerGetOrderParams) ([]SellerGetOrderRow, error) {
	rows, err := q.db.Query(ctx, sellerGetOrder, arg.SellerName, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SellerGetOrderRow
	for rows.Next() {
		var i SellerGetOrderRow
		if err := rows.Scan(
			&i.ID,
			&i.ShopID,
			&i.Shipment,
			&i.TotalPrice,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const sellerGetReport = `-- name: SellerGetReport :many



INSERT INTO
    "product"(
        "version",
        "shop_id",
        "name",
        "description",
        "price",
        "image_id",
        "exp_date"
    )
VALUES (
        0, (
            SELECT s."id"
            FROM "shop" s
            WHERE
                s."seller_name" = $1
                AND s."enabled" = true
        ),
        $2,
        $3,
        $4,
        $5,
        $6
    ) RETURNING "id"
`

type SellerGetReportParams struct {
	SellerName  string             `json:"seller_name" params:"seller_name"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       pgtype.Numeric     `json:"price"`
	ImageID     pgtype.UUID        `json:"image_id"`
	ExpDate     pgtype.Timestamptz `json:"exp_date"`
}

func (q *Queries) SellerGetReport(ctx context.Context, arg SellerGetReportParams) ([]int32, error) {
	rows, err := q.db.Query(ctx, sellerGetReport,
		arg.SellerName,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.ImageID,
		arg.ExpDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const sellerInsertCoupon = `-- name: SellerInsertCoupon :one

INSERT INTO
    "coupon" (
        "type",
        "shop_id",
        "description",
        "discount",
        "start_date",
        "expire_date"
    )
VALUES (
        $2, (
            SELECT s."id"
            FROM "shop" s
            WHERE
                s."seller_name" = $1
                AND s."enabled" = true
        ),
        $3,
        $4,
        NOW(),
        $5
    ) RETURNING "id"
`

type SellerInsertCouponParams struct {
	SellerName  string             `json:"seller_name" params:"seller_name"`
	Type        CouponType         `json:"type"`
	Description string             `json:"description"`
	Discount    pgtype.Numeric     `json:"discount"`
	ExpireDate  pgtype.Timestamptz `json:"expire_date"`
}

func (q *Queries) SellerInsertCoupon(ctx context.Context, arg SellerInsertCouponParams) (int32, error) {
	row := q.db.QueryRow(ctx, sellerInsertCoupon,
		arg.SellerName,
		arg.Type,
		arg.Description,
		arg.Discount,
		arg.ExpireDate,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const updateCouponInfo = `-- name: UpdateCouponInfo :exec

UPDATE "coupon" c
SET
    "type" = COALESCE($3, "type"),
    "description" = COALESCE($4, "description"),
    "discount" = COALESCE($4, "discount"),
    "start_date" = COALESCE($4, "start_date"),
    "expire_date" = COALESCE($4, "expire_date")
WHERE c."id" = $2 AND "shop_id" = (
        SELECT s."id"
        FROM "shop" s
        WHERE
            s."seller_name" = $1
            AND s."enabled" = true
    )
`

type UpdateCouponInfoParams struct {
	SellerName  string     `json:"seller_name" params:"seller_name"`
	ID          int32      `json:"id" params:"id"`
	Type        CouponType `json:"type"`
	Description string     `json:"description"`
}

func (q *Queries) UpdateCouponInfo(ctx context.Context, arg UpdateCouponInfoParams) error {
	_, err := q.db.Exec(ctx, updateCouponInfo,
		arg.SellerName,
		arg.ID,
		arg.Type,
		arg.Description,
	)
	return err
}

const updateProductInfo = `-- name: UpdateProductInfo :exec

UPDATE "product" p
SET
    "name" = COALESCE($3, "name"),
    "description" = COALESCE($4, "description"),
    "price" = COALESCE($5, "price"),
    "image_id" = COALESCE($6, "image_id"),
    "exp_date" = COALESCE($7, "exp_date"),
    "description" = COALESCE($8, "description"),
    "edit_date" = NOW(),
    "version" = "version" + 1
WHERE "shop_id" = (
        SELECT s."id"
        FROM "shop" s
        WHERE
            s."seller_name" = $1
            AND s."enabled" = true
    )
    AND p."id" = $2
`

type UpdateProductInfoParams struct {
	SellerName    string             `json:"seller_name" params:"seller_name"`
	ID            int32              `json:"id" params:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Price         pgtype.Numeric     `json:"price"`
	ImageID       pgtype.UUID        `json:"image_id"`
	ExpDate       pgtype.Timestamptz `json:"exp_date"`
	Description_2 string             `json:"description_2"`
}

func (q *Queries) UpdateProductInfo(ctx context.Context, arg UpdateProductInfoParams) error {
	_, err := q.db.Exec(ctx, updateProductInfo,
		arg.SellerName,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.ImageID,
		arg.ExpDate,
		arg.Description_2,
	)
	return err
}

const updateSellerInfo = `-- name: UpdateSellerInfo :exec

UPDATE "shop"
SET
    "image_id" = COALESCE($2, "image_id"),
    "name" = COALESCE($3, "name"),
    "description" = COALESCE($4, "description"),
    "enabled" = COALESCE($5, "enabled")
WHERE "seller_name" IN (
        SELECT "username"
        FROM "user" u
        WHERE u.id = $1
    )
`

type UpdateSellerInfoParams struct {
	ID          int32       `json:"id" params:"id"`
	ImageID     pgtype.UUID `json:"image_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled"`
}

func (q *Queries) UpdateSellerInfo(ctx context.Context, arg UpdateSellerInfoParams) error {
	_, err := q.db.Exec(ctx, updateSellerInfo,
		arg.ID,
		arg.ImageID,
		arg.Name,
		arg.Description,
		arg.Enabled,
	)
	return err
}
