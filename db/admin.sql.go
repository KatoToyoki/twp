// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: admin.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addCoupon = `-- name: AddCoupon :one

INSERT INTO
    "coupon" (
        "type",
        "scope",
        "shop_id",
        "name",
        "description",
        "discount",
        "start_date",
        "expire_date"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "id",
    "type",
    "scope",
    "name",
    "description",
    "discount",
    "start_date",
    "expire_date"
`

type AddCouponParams struct {
	Type        CouponType         `json:"type"`
	Scope       CouponScope        `json:"scope"`
	ShopID      pgtype.Int4        `json:"-"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Discount    pgtype.Numeric     `json:"discount" swaggertype:"number"`
	StartDate   pgtype.Timestamptz `json:"start_date" swaggertype:"string"`
	ExpireDate  pgtype.Timestamptz `json:"expire_date" swaggertype:"string"`
}

type AddCouponRow struct {
	ID          int32              `json:"id" param:"coupon_id"`
	Type        CouponType         `json:"type"`
	Scope       CouponScope        `json:"scope"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Discount    pgtype.Numeric     `json:"discount" swaggertype:"number"`
	StartDate   pgtype.Timestamptz `json:"start_date" swaggertype:"string"`
	ExpireDate  pgtype.Timestamptz `json:"expire_date" swaggertype:"string"`
}

func (q *Queries) AddCoupon(ctx context.Context, arg AddCouponParams) (AddCouponRow, error) {
	row := q.db.QueryRow(ctx, addCoupon,
		arg.Type,
		arg.Scope,
		arg.ShopID,
		arg.Name,
		arg.Description,
		arg.Discount,
		arg.StartDate,
		arg.ExpireDate,
	)
	var i AddCouponRow
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Scope,
		&i.Name,
		&i.Description,
		&i.Discount,
		&i.StartDate,
		&i.ExpireDate,
	)
	return i, err
}

const couponExists = `-- name: CouponExists :one

SELECT EXISTS( SELECT 1 FROM "coupon" WHERE "id" = $1 )
`

func (q *Queries) CouponExists(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRow(ctx, couponExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const deleteCoupon = `-- name: DeleteCoupon :execrows

WITH _ AS (
        DELETE FROM
            "cart_coupon"
        WHERE "coupon_id" = $1
    )
DELETE FROM "coupon"
WHERE "id" = $1
`

func (q *Queries) DeleteCoupon(ctx context.Context, id int32) (int64, error) {
	result, err := q.db.Exec(ctx, deleteCoupon, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const disableProductsFromShop = `-- name: DisableProductsFromShop :execrows

UPDATE "product" AS p SET p."enabled" = FALSE WHERE p."shop_id" = $1
`

func (q *Queries) DisableProductsFromShop(ctx context.Context, shopID int32) (int64, error) {
	result, err := q.db.Exec(ctx, disableProductsFromShop, shopID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const disableShop = `-- name: DisableShop :execrows


WITH disable_shop AS (
        UPDATE "shop" AS s
        SET s."enabled" = FALSE
        WHERE
            s."seller_name" = $1 RETURNING s."id"
    )
UPDATE "product" AS p
SET p."enabled" = FALSE
WHERE p."shop_id" = (
        SELECT "id"
        FROM disable_shop
    )
`

// there are some sql 🪄 happening here
func (q *Queries) DisableShop(ctx context.Context, sellerName string) (int64, error) {
	result, err := q.db.Exec(ctx, disableShop, sellerName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const disableUser = `-- name: DisableUser :execrows

WITH disabled_user AS (
        UPDATE "user"
        SET "enabled" = FALSE
        WHERE
            "username" = $1 RETURNING "username"
    ),
    disabled_shop AS (
        UPDATE "shop"
        SET "enabled" = FALSE
        WHERE "seller_name" = (
                SELECT
                    "username"
                FROM
                    disabled_user
            ) RETURNING "id"
    )
UPDATE "product"
SET "enabled" = FALSE
WHERE "shop_id" = (
        SELECT "id"
        FROM disabled_shop
    )
`

func (q *Queries) DisableUser(ctx context.Context, username string) (int64, error) {
	result, err := q.db.Exec(ctx, disableUser, username)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const editCoupon = `-- name: EditCoupon :execrows


UPDATE "coupon"
SET
    "type" = COALESCE($2, "type"),
    "name" = COALESCE($3, "name"),
    "description" = COALESCE($4, "description"),
    "discount" = COALESCE($5, "discount"),
    "start_date" = COALESCE($6, "start_date"),
    "expire_date" = COALESCE($7, "expire_date")
WHERE
    "id" = $1 RETURNING "id",
    "type",
    "scope",
    "name",
    "description",
    "discount",
    "start_date",
    "expire_date"
`

type EditCouponParams struct {
	ID          int32              `json:"id" param:"coupon_id"`
	Type        CouponType         `json:"type"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Discount    pgtype.Numeric     `json:"discount" swaggertype:"number"`
	StartDate   pgtype.Timestamptz `json:"start_date" swaggertype:"string"`
	ExpireDate  pgtype.Timestamptz `json:"expire_date" swaggertype:"string"`
}

type EditCouponRow struct {
	ID          int32              `json:"id" param:"coupon_id"`
	Type        CouponType         `json:"type"`
	Scope       CouponScope        `json:"scope"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Discount    pgtype.Numeric     `json:"discount" swaggertype:"number"`
	StartDate   pgtype.Timestamptz `json:"start_date" swaggertype:"string"`
	ExpireDate  pgtype.Timestamptz `json:"expire_date" swaggertype:"string"`
}

// i don't feel right about this
func (q *Queries) EditCoupon(ctx context.Context, arg EditCouponParams) (int64, error) {
	result, err := q.db.Exec(ctx, editCoupon,
		arg.ID,
		arg.Type,
		arg.Name,
		arg.Description,
		arg.Discount,
		arg.StartDate,
		arg.ExpireDate,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const enabledShop = `-- name: EnabledShop :execrows

UPDATE "shop" AS s SET s."enabled" = TRUE WHERE s."seller_name" = $1
`

func (q *Queries) EnabledShop(ctx context.Context, sellerName string) (int64, error) {
	result, err := q.db.Exec(ctx, enabledShop, sellerName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAnyCoupons = `-- name: GetAnyCoupons :many

SELECT id, type, scope, shop_id, name, description, discount, start_date, expire_date FROM "coupon" ORDER BY "id" ASC LIMIT $1 OFFSET $2
`

type GetAnyCouponsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAnyCoupons(ctx context.Context, arg GetAnyCouponsParams) ([]Coupon, error) {
	rows, err := q.db.Query(ctx, getAnyCoupons, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Coupon{}
	for rows.Next() {
		var i Coupon
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Scope,
			&i.ShopID,
			&i.Name,
			&i.Description,
			&i.Discount,
			&i.StartDate,
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

const getCouponDetail = `-- name: GetCouponDetail :one

SELECT
    "id",
    "type",
    "scope",
    "name",
    "description",
    "discount",
    "start_date",
    "expire_date"
FROM "coupon"
WHERE "id" = $1
`

type GetCouponDetailRow struct {
	ID          int32              `json:"id" param:"coupon_id"`
	Type        CouponType         `json:"type"`
	Scope       CouponScope        `json:"scope"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Discount    pgtype.Numeric     `json:"discount" swaggertype:"number"`
	StartDate   pgtype.Timestamptz `json:"start_date" swaggertype:"string"`
	ExpireDate  pgtype.Timestamptz `json:"expire_date" swaggertype:"string"`
}

func (q *Queries) GetCouponDetail(ctx context.Context, id int32) (GetCouponDetailRow, error) {
	row := q.db.QueryRow(ctx, getCouponDetail, id)
	var i GetCouponDetailRow
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Scope,
		&i.Name,
		&i.Description,
		&i.Discount,
		&i.StartDate,
		&i.ExpireDate,
	)
	return i, err
}

const getShopIDBySellerName = `-- name: GetShopIDBySellerName :one

SELECT "id" FROM "shop" WHERE "seller_name" = $1
`

func (q *Queries) GetShopIDBySellerName(ctx context.Context, sellerName string) (int32, error) {
	row := q.db.QueryRow(ctx, getShopIDBySellerName, sellerName)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getUserIDByUsername = `-- name: GetUserIDByUsername :one


SELECT "id" FROM "user" WHERE "username" = $1
`

// TODO name: GetReport :many
func (q *Queries) GetUserIDByUsername(ctx context.Context, username string) (int32, error) {
	row := q.db.QueryRow(ctx, getUserIDByUsername, username)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getUsers = `-- name: GetUsers :many

SELECT
    "username",
    "name",
    "email",
    "address",
    "role",
    "credit_card",
    "enabled"
FROM "user"
ORDER BY "id" ASC
LIMIT $1
OFFSET $2
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetUsersRow struct {
	Username   string   `json:"username"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Address    string   `json:"address"`
	Role       RoleType `json:"role"`
	CreditCard []byte   `json:"credit_card"`
	Enabled    bool     `json:"enabled"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]GetUsersRow, error) {
	rows, err := q.db.Query(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUsersRow{}
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(
			&i.Username,
			&i.Name,
			&i.Email,
			&i.Address,
			&i.Role,
			&i.CreditCard,
			&i.Enabled,
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
