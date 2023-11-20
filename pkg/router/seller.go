package router

import (
	"net/http"

	"github.com/jykuo-love-shiritori/twp/db"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// @Summary Seller get shop info
// @Description Get shop info, includes user picture, name, description.
// @Tags Seller, Shop
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller [get]
func sellerGetShopInfo(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {

	// Perform the query and handle errors
	// shopInfo, err := db.Queries.GetSellerInfo(context.Background(), 0)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return func(c echo.Context) error {

		// type shopInfomation struct {
		// 	name        string      `json:"name"`
		// 	image_id    pgtype.UUID `json:"image_id"`
		// 	description string      `json:"description"`
		// 	enabled     bool        `json:"enabled"`
		// }
		// result := shopInfomation{
		// 	shopInfo.Name,
		// 	shopInfo.ImageID,
		// 	shopInfo.Description,
		// 	shopInfo.Enabled,
		// }
		// return c.JSON(http.StatusOK, result)
		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller edit shop info
// @Description Edit shop name, description, visibility.
// @Tags Seller, Shop
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller [patch]
func sellerEditInfo(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get available tag
// @Description Get all available tags for shop.
// @Tags Seller, Shop, Tag
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/tag [get]
func sellerGetTag(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	//a void the let the wildcard
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller add tag
// @Description Add tag for shop.
// @Tags Seller, Shop, Tag
// @Accept json
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/tag [post]
func sellerAddTag(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get shop coupon
// @Description Get all coupons for shop.
// @Tags Seller, Shop, Coupon
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/coupon [get]
func sellerGetShopCoupon(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get coupon detail
// @Description Get coupon detail by ID for shop.
// @Tags Seller, Shop, Coupon
// @Produce json
// @Param id path int true "Coupon ID"
// @Success 200
// @Failure 401
// @Router /seller/coupon/{id} [get]
func sellerGetCouponDetail(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller add coupon
// @Description Add coupon for shop.
// @Tags Seller, Shop, Coupon
// @Accept json
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/coupon [post]
func sellerAddCoupon(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller edit coupon
// @Description Edit coupon for shop.
// @Tags Seller, Shop, Coupon
// @Accept json
// @Produce json
// @Param id path int true "Coupon ID"
// @Success 200
// @Failure 401
// @Router /seller/coupon/{id} [patch]
func sellerEditCoupon(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller delete coupon
// @Description Delete coupon for shop.
// @Tags Seller, Shop, Coupon
// @Accept json
// @Produce json
// @Success 200
// @Failure 401
func sellerDeleteCoupon(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get order
// @Description Get all orders for shop.
// @Tags Seller, Shop, Order
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/order [get]
func sellerGetOrder(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get order detail
// @Description Get order detail by ID for shop.
// @Tags Seller, Shop, Order
// @Produce json
// @Param id path int true "Order ID"
// @Success 200
// @Failure 401
// @Router /seller/order/{id} [get]
func sellerGetOrderDetail(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get report
// @Description Get all available reports for shop.
// @Tags Seller, Shop, Report
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/report [get]
func sellerGetReport(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller get report detail
// @Description Get report detail by year and month for shop.
// @Tags Seller, Shop, Report
// @Produce json
// @Param year path int true "Year"
// @Param month path int true "Month"
// @Success 200
// @Failure 401
// @Router /seller/report/{year}/{month} [get]
func sellerGetReportDetail(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller add product
// @Description Add product for shop.
// @Tags Seller, Shop, Product
// @Accept json
// @Produce json
// @Success 200
// @Failure 401
// @Router /seller/product [post]
func sellerAddProduct(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller upload product image
// @Description Upload product image for shop.
// @Tags Seller, Shop, Product
// @Accept png,jpeg,gif
// @Produce json
// @Param id path int true "Product ID"
// @Param img formData file true "image to upload"
// @Success 200
// @Failure 401
// @Router /seller/product/{id}/upload [post]
func sellerUploadProductImage(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller edit product
// @Description Edit product for shop.
// @Tags Seller, Shop, Product
// @Accept json
// @Produce json
// @Param id query int true "Product ID"
// @Success 200
// @Failure 401
// @Router /seller/product/{id} [patch]
func sellerEditProduct(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Seller delete product
// @Description Delete product for shop.
// @Tags Seller, Shop, Product
// @Accept json
// @Produce json
// @Param id query int true "Product ID"
// @Success 200
// @Failure 401
// @Router /seller/product/{id} [delete]
func sellerDeleteProduct(db *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}
