package cart

import (
	"context"
	"fmt"

	"github.com/ducklawrence/go-ecom/types"
)

func getCartItemsIDs(items []types.CartCheckoutItem) ([]int, error) {
	productIDs := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}
	return productIDs, nil
}

func (h *Handler) createOrder(
	ctx context.Context,
	products []types.Product,
	items []types.CartCheckoutItem,
	userID int,
) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// check if all products are actually in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce quantity of products in our db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		if err := h.productStore.UpdateProduct(ctx, product); err != nil {
			return 0, 0, err
		}
	}

	// create order
	orderID, err := h.store.CreateOrder(ctx, types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some random address",
	})
	if err != nil {
		return 0, 0, err
	}

	// create order items
	for _, item := range items {
		if err := h.store.CreateOrderItem(ctx, types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		}); err != nil {
			return 0, 0, err
		}
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("prdoduct %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartCheckoutItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		products := products[item.ProductID]
		total += products.Price * float64(item.Quantity)
	}

	return total
}
