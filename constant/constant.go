package constant

const (
	APIVersion = "v1"

	BadRequestMessage = "request not fulfilled"

	//schedular constants
	HealthCheckRoute = "/health"
	Database         = "EthnicElegance"

	// email verification routes
	VerifyEmailRoute = "/verify-email"
	VerifyOtpRoute   = "/verify-otp"
	ResendEmailRoute = "/resend-email"

	// user related routes
	UserRegisterRoute = "/user-register"
	UserLoginRoute    = "/login"
	UserLogoutRoute   = "/logout"

	// product and adminroutes
	RegisterProductRoute    = "/product-register"
	ListProductRoute        = "/products"
	ListProductRouteAdmin   = "/list-products-admin"
	ListOrders              = "/list-orders"
	UpdateOrder             = "/update-order"
	ListCategoryRoute       = "/list-category"
	ListSingleProductRoute  = "/product/:id"
	SearchProductRoute      = "/search"
	UpdateProductRoute      = "/update-product/:id"
	UpdateStockRoute        = "/update-stock/:id"
	DeleteProductRoute      = "/delete-product/:id"
	AddToCartRoute          = "/cart"
	AddAddressRoute         = "/address"
	EditAddressRoute        = "/address"
	EditNameRoute           = "/name"
	GetSingleUserRoute      = "/user/:id"
	UpdateUser              = "/update-user"
	CheckoutRoute           = "/checkout"
	AddToFavoriteRoute      = "/favorite"
	RemoveFromFavoriteRoute = "/remove-favorite"
	ListFavoriteRoute       = "/favorite"
	GetAllUserRoute         = "/users"
	BlockUserRoute          = "/block-user"
	UnblockUserRoute        = "/unblock-user"
	AddCategoryRoute        = "/category"
	UpdateCategoryRoute     = "/category/:id"
	DeleteCategoryRoute     = "/category/:id"
	AddCouponRoute          = "/coupon"
	DeleteCouponRoute       = "/coupon/:id"
	ListCouponRoute         = "/coupon"
	RemoveFromCartRoute     = "/cart/remove"
	UpdateCart              = "/cart/update"
	ListCartRoute           = "/cart"
	EmptyCartRoute          = "/cart/all"
	ApplyCouponRoute        = "/cart/coupon"
	GetProductLinkRoute     = "/product-link/:id"
	AddOfferRoute           = "/offer"
	UpdateOfferRoute        = "/offer/:id"
	GiveRatingRoute         = "/rating/:id"
	CommentOnProductRoute   = "/comment/:id"
)

const (
	NormalUser = "user"
	AdminUser  = "admin"
)

const (
	// time slot for otp validation
	OtpValidation = 60
)

// collections
const (
	VerificationsCollection = "verifications"
	UsersCollection         = "users"
	ProductCollection       = "products"
	AddressCollection       = "user_addresses"
	CartCollection          = "user_cart"
	CategoryCollection      = "categories"
	CouponCollection        = "coupons"
	OfferCollection         = "offers"
	CartItemCollection      = "cart_items"
	OrderCollection         = "orders"
)

// messages
const (
	AlreadyRegisterWithThisEmail = "already register with this email"
	EmailIsNotVerified           = "your email is not verified please verify your email"
	EmailValidationError         = "wrong email entered"
	EamilExists                  = "email already exists"
	OtpValidationError           = "wrong otp entered"
	OtpExpiredValidationError    = "otp expired"
	AlreadyVerifiedError         = "already verified"
	OptAlreadySentError          = "otp already sent to email"
	NotRegisteredUser            = "you are not a register user"
	PasswordNotMatchedError      = "email and password do not match"
	NotAuthorizedUserError       = "you are not authorized"
	NoProductAvaliable           = "no product avaliable"
	UserDoesNotExists            = "user not exists"
	AddressNotExists             = "address not exists. please add one address"
)
