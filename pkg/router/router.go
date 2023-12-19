package router

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	"github.com/jykuo-love-shiritori/twp/db"
	_ "github.com/jykuo-love-shiritori/twp/docs"
	"github.com/jykuo-love-shiritori/twp/minio"
	"github.com/jykuo-love-shiritori/twp/pkg/auth"
	"github.com/jykuo-love-shiritori/twp/pkg/constants"
	"github.com/jykuo-love-shiritori/twp/pkg/router/admin"
	"github.com/jykuo-love-shiritori/twp/pkg/router/buyer"
	"github.com/jykuo-love-shiritori/twp/pkg/router/general"
)

//	@title			twp API
//	@version		0.o
//	@description	twp server api.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func RegisterDocs(e *echo.Echo) {
	docs := e.Group(constants.SWAGGER_PATH)
	docs.GET("/*", echoSwagger.WrapHandler)
}

func RegisterApi(e *echo.Echo, pg *db.DB, mc *minio.MC, logger *zap.SugaredLogger) {
	api := e.Group("/api")

	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "pong"})
	}, auth.IsRole(pg, logger, db.RoleTypeCustomer))

	api.GET("/delay", func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		return c.JSON(http.StatusOK, map[string]string{"message": "delay"})
	})

	api.POST("/signup", auth.Signup(pg, logger))

	api.POST("/oauth/authorize", auth.Authorize(pg, logger))
	api.POST("/oauth/token", auth.Token(pg, logger))
	api.POST("/oauth/refresh", auth.Refresh(pg, logger))
	api.POST("/oauth/logout", auth.Logout(pg, logger), auth.ValidateJwt(pg, logger))

	// admin
	api.GET("/admin/user", admin.GetUser(pg, mc, logger))
	api.DELETE("/admin/user/:username", admin.DisableUser(pg, logger))

	api.GET("/admin/coupon", admin.GetCoupon(pg, logger))
	api.GET("/admin/coupon/:id", admin.GetCouponDetail(pg, logger))
	api.POST("/admin/coupon", admin.AddCoupon(pg, logger))
	api.PATCH("/admin/coupon/:id", admin.EditCoupon(pg, logger))
	api.DELETE("/admin/coupon/:id", admin.DeleteCoupon(pg, logger))

	api.GET("/admin/report", admin.GetReport(pg, mc, logger))

	// user
	api.GET("/user/info", userGetInfo(pg, logger))
	api.PATCH("/user/info", userEditInfo(pg, logger))
	api.POST("/user/info/upload", userUploadAvatar(pg, logger))
	api.POST("/user/security/password", userEditPassword(pg, logger))

	api.GET("/user/security/credit_card", userGetCreditCard(pg, logger))
	api.PATCH("/user/security/credit_card", userUpdateCreditCard(pg, logger))

	// general
	api.GET("/shop/:seller_name", general.GetShopInfo(pg, mc, logger)) // user
	api.GET("/shop/:seller_name/coupon", general.GetShopCoupon(pg, logger))
	api.GET("/shop/:seller_name/search", general.SearchShopProduct(pg, mc, logger))

	api.GET("/tag/:id", general.GetTagInfo(pg, logger))

	api.GET("/search", general.Search(pg, mc, logger)) // search both product and shop
	api.GET("/search/shop", general.SearchShopByName(pg, mc, logger))

	api.GET("/news", general.GetNews(pg, logger))
	api.GET("/news/:id", general.GetNewsDetail(pg, logger))
	api.GET("/discover", general.GetDiscover(pg, mc, logger))
	api.GET("/popular", general.GetPopular(pg, mc, logger))

	api.GET("/product/:id", general.GetProductInfo(pg, mc, logger))

	// buyer
	api.GET("/buyer/order", buyer.GetOrderHistory(pg, mc, logger))
	api.GET("/buyer/order/:id", buyer.GetOrderDetail(pg, mc, logger))
	api.PATCH("/buyer/order/:id", buyer.UpdateOrderStatus(pg, logger))

	api.GET("/buyer/cart", buyer.GetCart(pg, mc, logger)) // include product and coupon
	api.GET("/buyer/cart/:cart_id/coupon", buyer.GetCoupon(pg, logger))
	api.POST("/buyer/cart/product/:id", buyer.AddProductToCart(pg, logger))
	api.POST("/buyer/cart/:cart_id/coupon/:coupon_id", buyer.AddCouponToCart(pg, logger))
	api.PATCH("/buyer/cart/:cart_id/product/:product_id", buyer.EditProductInCart(pg, logger))
	api.DELETE("/buyer/cart/:cart_id/product/:product_id", buyer.DeleteProductFromCart(pg, logger))
	api.DELETE("/buyer/cart/:cart_id/coupon/:coupon_id", buyer.DeleteCouponFromCart(pg, logger))

	api.GET("/buyer/cart/:cart_id/checkout", buyer.GetCheckout(pg, logger))
	api.POST("/buyer/cart/:cart_id/checkout", buyer.Checkout(pg, logger))

	// seller
	api.GET("/seller/info", sellerGetShopInfo(pg, mc, logger))
	api.PATCH("/seller/info", sellerEditInfo(pg, mc, logger))
	api.GET("/seller/tag", sellerGetTag(pg, logger))  // search available tag
	api.POST("/seller/tag", sellerAddTag(pg, logger)) // add tag for shop

	api.GET("/seller/coupon", sellerGetShopCoupon(pg, logger))
	api.GET("/seller/coupon/:id", sellerGetCouponDetail(pg, logger))
	api.POST("/seller/coupon", sellerAddCoupon(pg, logger))
	api.PATCH("/seller/coupon/:id", sellerEditCoupon(pg, logger))
	api.DELETE("/seller/coupon/:id", sellerDeleteCoupon(pg, logger))
	api.POST("/seller/coupon/:id/tag", sellerAddCouponTag(pg, logger))
	api.DELETE("/seller/coupon/:id/tag", sellerDeleteCouponTag(pg, logger))

	api.GET("/seller/order", sellerGetOrder(pg, mc, logger))
	api.GET("/seller/order/:id", sellerGetOrderDetail(pg, mc, logger))
	api.PATCH("/seller/order/:id", sellerUpdateOrderStatus(pg, logger))

	api.GET("/seller/report/:year/:month", sellerGetReportDetail(pg, mc, logger))

	api.GET("/seller/product", sellerListProduct(pg, mc, logger))
	api.POST("/seller/product", sellerAddProduct(pg, mc, logger))
	api.GET("/seller/product/:id", sellerGetProductDetail(pg, mc, logger))
	api.PATCH("/seller/product/:id", sellerEditProduct(pg, mc, logger))
	api.POST("/seller/product/:id/tag", sellerAddProductTag(pg, logger))
	api.DELETE("/seller/product/:id/tag", sellerDeleteProductTag(pg, logger))
	api.DELETE("/seller/product/:id", sellerDeleteProduct(pg, logger))

}
