-- name: GetOrderHistory :many

SELECT
    O."id",
    s."name",
    s."image_id" AS "shop_image_id",
    O."image_id" AS "thumbnail_id",
    "shipment",
    "total_price",
    "status",
    "created_at"
FROM
    "order_history" AS O,
    "user" AS U,
    "shop" AS S
WHERE
    U."username" = $1
    AND U."id" = O."user_id"
    AND O."shop_id" = S."id"

ORDER BY "created_at" ASC OFFSET $2 LIMIT $3;

-- name: GetOrderInfo :one

SELECT
    O."id",
    s."name",
    s."image_id",
    "shipment",
    "total_price",
    "status",
    "created_at", (
        "subtotal" + "shipment" - "total_price"
    ) AS "discount"
FROM
    "order_history" AS O,
    "order_detail" AS D,
    "product_archive" AS P,
    "user" AS U,
    "shop" AS S, (
        SELECT
            SUM(P."price" * D."quantity") AS "subtotal"
        FROM
            "order_detail" AS D,
            "product_archive" AS P
        WHERE
            D."order_id" = $1
            AND D."product_id" = P."id"
            AND D."product_version" = P."version"
    ) AS T
WHERE
    U."username" = $2
    AND O."id" = $1;

-- name: GetOrderDetail :many

SELECT
    O."product_id",
    P."name",
    P."description",
    P."price",
    P."image_id",
    O."quantity"
FROM
    "order_detail" AS O,
    "product_archive" AS P
WHERE
    O."order_id" = $1
    AND O."product_id" = P."id"
    AND O."product_version" = P."version";

-- name: GetCart :many

SELECT
    C."id",
    S."seller_name",
    S."image_id",
    S."name"
FROM
    "cart" AS C,
    "user" AS U,
    "shop" AS S
WHERE
    U."username" = $1
    AND U."id" = C."user_id"
    AND C."shop_id" = S."id";

-- name: GetProductInCart :many

SELECT
    "product_id",
    "name",
    "image_id",
    "price",
    "quantity"
FROM
    "cart_product" AS C,
    "product" AS P
WHERE
    "cart_id" = $1
    AND C."product_id" = P."id";

-- name: UpdateProductInCart :one

UPDATE "cart_product"
SET "quantity" = $3
FROM "user" AS U, "cart" AS C
WHERE
    U."username" = $4
    AND U."id" = C."user_id"
    AND "cart_id" = $1
    AND "product_id" = $2 RETURNING "quantity";

-- name: DeleteProductInCart :execrows

WITH valid_cart AS (
        SELECT C."id"
        FROM "cart" C
            JOIN "user" u ON u."id" = C."user_id"
        WHERE
            u."username" = $3
            AND C."id" = $1
    ),
    deleted_products AS (
        DELETE FROM
            "cart_product" CP
        WHERE "cart_id" = (
                SELECT "id"
                FROM
                    valid_cart
            )
            AND CP."product_id" = $2 RETURNING 1
    ),
    remaining_products AS (
        SELECT COUNT(*) AS count
        FROM "cart_product"
        WHERE "cart_id" = (
                SELECT "id"
                FROM valid_cart
            )
    )
DELETE FROM "cart" AS 🛒
WHERE 🛒."id" = $1 AND (
        SELECT count
        FROM remaining_products
    ) = 0;

-- i hope this works ☠️

-- name: AddProductToCart :one

WITH valid_product AS (
        SELECT P."id", S."id"
        FROM
            "product" P,
            "shop" S
        WHERE
            P."shop_id" = S."id"
            AND P."id" = $3
            AND P."enabled" = TRUE
    ),
    -- check product enabled ⬆️
    new_cart AS (
        INSERT INTO
            "cart" ("user_id", "shop_id")
        SELECT
            U."id",
            S."shop_id"
        FROM
            "user" AS U,
            "shop" AS S,
            "product" AS P
        WHERE
            U."username" = $1
            AND S."id" = P."shop_id"
            AND NOT EXISTS (
                SELECT 1
                FROM
                    "cart" AS C
                WHERE
                    C."user_id" = U."id"
                    AND C."shop_id" = S."shop_id"
            ) RETURNING "id"
    ),
    -- create new cart if not exists ⬆️
    existing_cart_product AS (
        UPDATE
            "cart_product" AS CP
        SET
            "quantity" = "quantity" + $2
        FROM
            "cart" AS C,
            "user" AS U
        WHERE
            U."username" = $1
            AND C."user_id" = U."id"
            AND C."id" = CP."cart_id"
            AND CP."product_id" = (
                SELECT "id"
                FROM
                    valid_product
            ) RETURNING 1
    ) -- if the product already exists in the cart, update the quantity ⬆️
INSERT INTO
    -- insert into the cart that have no given product ⬇️
    "cart_product" (
        "cart_id",
        "product_id",
        "quantity"
    )
SELECT
    C."id",
    valid_product."id",
    $2
FROM "cart" C, valid_product
WHERE NOT EXISTS (
        SELECT 1
        FROM
            existing_cart_product
    ) RETURNING (
        SELECT COUNT(*)
        FROM
            "cart_product" CP,
            "cart" C,
            "user" U
        WHERE
            CP."cart_id" = C."id"
            AND U."id" = C."user_id"
            AND U."username" = $1
    );

-- returning the number of products in any cart for US-SC-2 in SRS ⬆️
