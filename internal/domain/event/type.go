package event

type Type string

const (
	ProductViewed    Type = "product_viewed"
	SearchPerformed  Type = "search_performed"
	ProductPurchased Type = "product_purchased"
	ProductAddedCart Type = "product_added_to_cart"
)
