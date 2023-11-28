package router

import (
	"net/http"

	"github.com/jykuo-love-shiritori/twp/db"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type orderDetail struct {
	OrderInfo db.SellerOrderCheckRow       `json:"order_info"`
	Products  []db.SellerGetOrderDetailRow `json:"products"`
}
type productDetail struct {
	ProductInfo db.SellerGetProductDetailRow `json:"product_info"`
	Tags        []db.ProductTag              `json:"tags"`
}
type couponDetail struct {
	CouponInfo db.SellerGetCouponDetailRow `json:"coupon_info"`
	Tags       []db.CouponTag              `json:"tags"`
}

// @Summary		Seller get shop info
// @Description	Get shop info, includes user picture, name, description.
// @Tags			Seller, Shop
// @Produce		json
// @success		200	{object}	db.Shop
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/info [get]
func sellerGetShopInfo(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"
		shopInfo, err := pg.Queries.SellerGetInfo(c.Request().Context(), username)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update seller infomation")
		}
		return c.JSON(http.StatusOK, shopInfo)

	}
}

// @Summary		Seller edit shop info
// @Description	Edit shop name, description, visibility.
// @Tags			Seller, Shop
// @Param			image_id	body	string	true	"update image UUID"
// @Param			name		body	string	true	"update shop name"	minlength(6)
// @Param			description	body	string	true	"update description"
// @Param			enabled		body	bool	true	"update enabled status"
// @Produce		json
// @success		200	{object}	db.SellerUpdateInfoRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/info [patch]
func sellerEditInfo(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerUpdateInfoParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")
		}
		param.SellerName = username
		shopInfo, err := pg.Queries.SellerUpdateInfo(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update seller infomation")
		}
		return c.JSON(http.StatusOK, shopInfo)
	}
}

// @Summary		Seller get available tag
// @Description	Get all available tags for shop.
// @Tags			Seller, Shop, Tag
// @Param			name	body	string	true	"search tagname start with"	minlength(1)
// @Produce		json
// @success		200	{object}	db.SellerSearchTagRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/tag [get]
func sellerGetTag(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {

	return func(c echo.Context) error {
		var username string = "user1"
		var tagPerPage int32 = 20

		var param db.SellerSearchTagParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		if param.Name == "" || hasSpecialChars(param.Name) {
			logger.Error("tag name invaild")
			return echo.NewHTTPError(http.StatusBadRequest, "tag name invaild")
		}
		param.SellerName = username
		param.Name = "^" + param.Name
		param.Limit = tagPerPage
		tags, err := pg.Queries.SellerSearchTag(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to search tag")
		}
		return c.JSON(http.StatusOK, tags)
	}
}

// @Summary		Seller add tag
// @Description	Add tag for shop.
// @Tags			Seller, Shop, Tag
// @Accept			json
// @Param			name	body	string	true	"insert tag"	minlength(1)
// @Produce		json
// @success		200	{object}	db.SellerInsertTagRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		409	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/tag [post]
func sellerAddTag(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.HaveTagNameParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		if param.Name == "" || hasSpecialChars(param.Name) {
			logger.Error("tag name invaild")
			return echo.NewHTTPError(http.StatusBadRequest, "tag name invaild")
		}
		param.SellerName = username
		have, err := pg.Queries.HaveTagName(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to check tag uniquity")
		}
		if have {
			logger.Error("Conflict (tag name have to be unique")
			return echo.NewHTTPError(http.StatusConflict, "Conflict (tag name have to be unique")
		}
		tag, err := pg.Queries.SellerInsertTag(c.Request().Context(), db.SellerInsertTagParams(param))
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to add tag")
		}
		return c.JSON(http.StatusOK, tag)
	}
}

// @Summary		Seller get shop coupon
// @Description	Get all coupons for shop.
// @Tags			Seller, Shop, Coupon
// @Param			offset	body	int	true	"offset page"	minimum(0)
// @Produce		json
// @success		200	{array}		db.SellerGetCouponRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/coupon [get]
func sellerGetShopCoupon(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"
		var couponPerPage int32 = 20

		var param db.SellerGetCouponParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		param.Limit = couponPerPage
		param.Offset = param.Offset * couponPerPage
		coupons, err := pg.Queries.SellerGetCoupon(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get seller coupon")
		}
		return c.JSON(http.StatusOK, coupons)
	}
}

// @Summary		Seller get coupon detail
// @Description	Get coupon detail by ID for shop.
// @Tags			Seller, Shop, Coupon
// @Produce		json
// @Param			id	path		int	true	"Coupon ID"
// @success		200	{object}	couponDetail
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/coupon/{id} [get]
func sellerGetCouponDetail(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		var username string = "user1"

		var param db.SellerGetCouponDetailParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		var result couponDetail
		var err error
		param.SellerName = username
		result.CouponInfo, err = pg.Queries.SellerGetCouponDetail(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get coupon detail")
		}
		result.Tags, err = pg.Queries.SellerGetCouponTag(c.Request().Context(), db.SellerGetCouponTagParams{SellerName: param.SellerName, CouponID: param.ID})
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get coupon tags")
		}

		return c.JSON(http.StatusOK, result)
	}
}

// @Summary		Seller add coupon
// @Description	Add coupon for shop.
// @Tags			Seller, Shop, Coupon
// @Param			type		body	string	true	"Coupon type"	Enums('percentage', 'fixed', 'shipping')
// @Param			name		body	string	true	"name of coupon"
// @Param			description	body	string	true	"description of coupon"
// @Param			discount	body	number	false	"discount perscent"
// @Param			start_date	body	time	true	"start date"
// @Param			expire_date	body	time	true	"expire date"
// @Accept			json
// @Produce		json
// @success		200	{object}	db.Coupon
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/coupon [post]
func sellerAddCoupon(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerInsertCouponParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		coupon, err := pg.Queries.SellerInsertCoupon(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to add coupon")
		}
		return c.JSON(http.StatusOK, coupon)
	}
}

// @Summary		Seller add coupon tag
// @Description	Add tag on coupon
// @Tags			Seller, Shop, Coupon,Tag
// @Accept			json
// @Param			id		path	string	true	"coupon id"
// @Param			tag_id	body	int		true	"add tag id"
// @Produce		json
// @success		200	{object}	db.CouponTag
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/coupon/{id}/tag [post]
func sellerAddCouponTag(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerInsertCouponTagParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")
		}
		param.SellerName = username
		couponTag, err := pg.Queries.SellerInsertCouponTag(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to add coupon tag")
		}
		return c.JSON(http.StatusOK, couponTag)
	}
}

// @Summary		Seller edit coupon
// @Description	Edit coupon for shop.
// @Tags			Seller, Shop, Coupon
// @Accept			json
// @Produce		json
// @Param			id			path		int		true	"Coupon ID"
// @Param			type		body		string	true	"Coupon type"	Enums('percentage', 'fixed', 'shipping')
// @Param			name		body		string	true	"name of coupon"
// @Param			description	body		string	true	"description of coupon"
// @Param			discount	body		number	false	"discount perscent"
// @Param			start_date	body		time	true	"start date"
// @Param			expire_date	body		time	true	"expire date"
// @success		200			{object}	db.Coupon
// @Failure		400			{object}	echo.HTTPError
// @Failure		500			{object}	echo.HTTPError
// @Router			/seller/coupon/{id} [patch]
func sellerEditCoupon(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.UpdateCouponInfoParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		coupon, err := pg.Queries.UpdateCouponInfo(c.Request().Context(), param)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update coupon")
		}
		return c.JSON(http.StatusOK, coupon)
	}
}

// @Summary		Seller delete coupon
// @Description	Delete coupon for shop.
// @Tags			Seller, Shop, Coupon
// @Param			id	path	int	true	"Coupon ID"
// @Accept			json
// @Produce		json
// @Success		200	{string}	string	"success"
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/coupon/{id} [delete]
func sellerDeleteCoupon(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerDeleteCouponParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		effectRow, err := pg.Queries.SellerDeleteCoupon(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete coupon")
		}
		if effectRow == 0 {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusNotFound, "Not Found (Coupon)")
		}
		return c.JSON(http.StatusOK, "success")
	}
}

// @Summary		Seller delete coupon tag
// @Description	Delete coupon for shop.
// @Tags			Seller, Shop, Coupon,Tag
// @Param			id		path	string	true	"coupon id"
// @Param			tag_id	body	int		true	"add tag id"
// @Accept			json
// @Produce		json
// @Success		200	{string}	string	"success"
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/coupon/{id}/tag [delete]
func sellerDeleteCouponTag(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerDeleteCouponTagParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		effectRow, err := pg.Queries.SellerDeleteCouponTag(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete coupon tag")
		}
		if effectRow == 0 {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusNotFound, "Not Found (coupon_id or tag_id)")
		}
		return c.JSON(http.StatusOK, "success")
	}
}

// @Summary		Seller get order
// @Description	Get all orders for shop.
// @Tags			Seller, Shop, Order
// @Param			offset	body	int	true	"offset page"	minimum(0)
// @Produce		json
// @Success		200	{object}	db.SellerGetOrderRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/order [get]
func sellerGetOrder(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		var username string = "user1"
		var orderPerPage int32 = 20

		var param db.SellerGetOrderParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		param.Limit = orderPerPage
		param.Offset = param.Offset * orderPerPage
		orders, err := pg.Queries.SellerGetOrder(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get order")
		}
		return c.JSON(http.StatusOK, orders)
	}
}

// @Summary		Seller get order detail
// @Description	Get order detail by ID for shop.
// @Tags			Seller, Shop, Order
// @Produce		json
// @Param			id		path		int	true	"Order ID"
// @Param			offset	body		int	true	"offset page"	minimum(0)
// @Success		200		{object}	orderDetail
// @Failure		400		{object}	echo.HTTPError
// @Failure		404		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router			/seller/order/{id} [get]
func sellerGetOrderDetail(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"
		var orderPerPage int32 = 20

		var param db.SellerGetOrderDetailParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		param.Limit = orderPerPage
		param.Offset = orderPerPage * param.Offset
		var result orderDetail
		var err error
		result.OrderInfo, err = pg.Queries.SellerOrderCheck(c.Request().Context(), db.SellerOrderCheckParams{SellerName: param.SellerName, ID: param.OrderID})
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to find order")
		}
		result.Products, err = pg.Queries.SellerGetOrderDetail(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get order detail product")
		}
		return c.JSON(http.StatusOK, result)
	}
}

// @Summary		Seller update order status
// @Description	seller update orders status.
// @Tags			Seller, Shop, Order
// @Param			id				path		int		true	"Order ID"
// @Param			current_status	body		string	true	"order status"	Enums('pending','paid','shipped','delivered','cancelled')
// @Param			set_status		body		string	true	"order status"	Enums('pending','paid','shipped','delivered','cancelled')
// @Success		200				{object}	db.OrderHistory
// @Failure		400				{object}	echo.HTTPError
// @Failure		500				{object}	echo.HTTPError
// @Router			/seller/order [patch]
func sellerUpdateOrderStatus(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerUpdateOrderStatusParams
		if err := c.Bind(&param); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username

		// shop can only a prove the status traction {paid > shipped ,shipped > delivered}
		// paid > shipped > delivered > (canelled || finished)
		if !((param.CurrentStatus == "paid" && param.SetStatus == "shipped") || (param.CurrentStatus == "shipped" && param.SetStatus == "delivered")) {
			logger.Error("invaild status chage")
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request (invaild status chage)")
		}
		order, err := pg.Queries.SellerUpdateOrderStatus(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update order status")
		}
		return c.JSON(http.StatusOK, order)
	}
}

// @Summary		Seller get report
// @Description	Get all available reports for shop.
// @Tags			Seller, Shop, Report
// @Produce		json
// @Success		200
// @Failure		401
// @Router			/seller/report [get]
func sellerGetReport(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary		Seller get report detail
// @Description	Get report detail by year and month for shop.
// @Tags			Seller, Shop, Report
// @Produce		json
// @Param			year	path	int	true	"Year"
// @Param			month	path	int	true	"Month"
// @Success		200
// @Failure		401
// @Router			/seller/report/{year}/{month} [get]
func sellerGetReportDetail(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary		Seller get product
// @Description	Delete product for shop.
// @Tags			Seller, Shop, Product
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Product ID"
// @Success		200	{object}	productDetail
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/product/{id} [get]
func sellerGetProductDetail(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerGetProductDetailParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		var result productDetail
		var err error
		param.SellerName = username
		result.ProductInfo, err = pg.Queries.SellerGetProductDetail(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get product detail")
		}
		result.Tags, err = pg.Queries.SellerGetProductTag(c.Request().Context(), db.SellerGetProductTagParams{SellerName: param.SellerName, ProductID: param.ID})
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get product tag")
		}
		return c.JSON(http.StatusOK, result)
	}
}

// @Summary		Seller get product
// @Description	seller get product
// @Tags			Seller, Shop, Product
// @Param			offset	body	int	true	"offset page"	minimum(0)
// @Accept			json
// @Produce		json
// @Success		200	{array}		db.SellerProductListRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/product [get]
func sellerListProduct(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"
		var orderPerPage int32 = 20

		var param db.SellerProductListParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}

		param.SellerName = username
		param.Limit = orderPerPage
		param.Offset = orderPerPage * param.Offset
		products, err := pg.Queries.SellerProductList(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get product")
		}
		return c.JSON(http.StatusOK, products)
	}
}

// @Summary		Seller add product
// @Description	Add product for shop.
// @Tags			Seller, Shop, Product
// @Param			name		body	string	true	"name of product"
// @Param			description	body	string	true	"description of product"
// @Param			price		body	number	false	"price"
// @Param			image_id	body	string	true	"image id"
// @Param			expire_date	body	time	true	"expire date"
// @Param			stock		body	int		true	"stock"
// @Param			enabled		body	time	true	"enabled"
// @Accept			json
// @Produce		json
// @Success		200	{object}	db.SellerInsertProductRow
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/product [post]
func sellerAddProduct(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerInsertProductParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		product, err := pg.Queries.SellerInsertProduct(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert product")
		}
		return c.JSON(http.StatusOK, product)
	}
}

// @Summary		Seller upload product image
// @Description	Upload product image for shop.
// @Tags			Seller, Shop, Product
// @Accept			png,jpeg,gif
// @Produce		json
// @Param			id	path		int		true	"Product ID"
// @Param			img	formData	file	true	"image to upload"
// @Success		200
// @Failure		401
// @Router			/seller/product/{id}/upload [post]
func sellerUploadProductImage(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// @Summary		Seller edit product
// @Description	Edit product for shop.
// @Tags			Seller, Shop, Product
// @Accept			json
// @Produce		json
// @Param			id			path		int		true	"Product ID"
// @Param			name		body		string	true	"name of product"
// @Param			description	body		string	true	"description of product"
// @Param			price		body		number	false	"price"
// @Param			image_id	body		string	true	"image id"
// @Param			expire_date	body		time	true	"expire date"
// @Param			stock		body		int		true	"stock"
// @Param			enabled		body		time	true	"enabled"
// @Success		200			{object}	db.SellerUpdateProductInfoRow
// @Failure		400			{object}	echo.HTTPError
// @Failure		500			{object}	echo.HTTPError
// @Router			/seller/product/{id} [patch]
func sellerEditProduct(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerUpdateProductInfoParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		product, err := pg.Queries.SellerUpdateProductInfo(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get product tag")
		}
		return c.JSON(http.StatusOK, product)
	}
}

// @Summary		Seller add product tag
// @Description	Add tag on product
// @Tags			Seller, Shop, Product,Tag
// @Accept			json
// @Param			id		path	string	true	"product id"
// @Param			tag_id	body	int		true	"add tag id"
// @Produce		json
// @Success		200	{object}	db.ProductTag
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/product/{id}/tag [post]
func sellerAddProductTag(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerInsertProductTagParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		productTag, err := pg.Queries.SellerInsertProductTag(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to add product tag")
		}
		return c.JSON(http.StatusOK, productTag)
	}
}

// @Summary		Seller delete product
// @Description	Delete product for shop.
// @Tags			Seller, Shop, Product
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Product ID"
// @Success		200	{string}	string	"success"
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/product/{id} [delete]
func sellerDeleteProduct(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerDeleteProductParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		effectRow, err := pg.Queries.SellerDeleteProduct(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete product")
		}
		if effectRow == 0 {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusNotFound, "Not Found (Product)")
		}
		return c.JSON(http.StatusOK, "success")

	}
}

// @Summary		Seller delete product tag
// @Description	Delete product for shop.
// @Tags			Seller, Shop, Coupon
// @Param			id		path	int	true	"product id"
// @Param			tag_id	body	int	true	"add tag id"
// @Accept			json
// @Produce		json
// @Success		200	{string}	string	"success"
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/seller/product/{id}/tag [delete]
func sellerDeleteProductTag(pg *db.DB, logger *zap.SugaredLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var username string = "user1"

		var param db.SellerDeleteProductTagParams
		if err := c.Bind(&param); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "paramter error")

		}
		param.SellerName = username
		effectRow, err := pg.Queries.SellerDeleteProductTag(c.Request().Context(), param)
		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete product tag")
		}
		if effectRow == 0 {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusNotFound, "Not Found (product_id or tag_id)")
		}
		return c.JSON(http.StatusOK, "success")
	}
}
