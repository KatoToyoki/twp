// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: buyer.sql

package db

import (
	"context"
)

const buyerCart = `-- name: BuyerCart :many

SELECT id, user_id, shop_id FROM "cart" WHERE "user_id" = $1
`

func (q *Queries) BuyerCart(ctx context.Context, userID int32) ([]Cart, error) {
	rows, err := q.db.Query(ctx, buyerCart, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cart
	for rows.Next() {
		var i Cart
		if err := rows.Scan(&i.ID, &i.UserID, &i.ShopID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const buyerGetOrder = `-- name: BuyerGetOrder :many

SELECT id, user_id, shop_id, shipment, total_price, status, created_at FROM "order_history" WHERE "id" = $1 and "user_id" = $2
`

type BuyerGetOrderParams struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"`
}

func (q *Queries) BuyerGetOrder(ctx context.Context, arg BuyerGetOrderParams) ([]OrderHistory, error) {
	rows, err := q.db.Query(ctx, buyerGetOrder, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrderHistory
	for rows.Next() {
		var i OrderHistory
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
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
