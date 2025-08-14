package cart

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ducklawrence/go-ecom/service/auth"
	"github.com/ducklawrence/go-ecom/service/order"
	"github.com/ducklawrence/go-ecom/service/product"
	"github.com/ducklawrence/go-ecom/service/user"
	"github.com/ducklawrence/go-ecom/types"
	"github.com/ducklawrence/go-ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

type HandlerTx struct {
	db *sql.DB
}

func NewHandler(
	store types.OrderStore,
	productStore types.ProductStore,
	userStore types.UserStore,
) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
		userStore:    userStore,
	}
}

func NewHandlerTx(db *sql.DB) *HandlerTx {
	return &HandlerTx{
		db: db,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(
		"/cart/checkout",
		auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *HandlerTx) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(
		"/cart/checkout-tx",
		auth.WithJWTAuth(h.handleCheckoutTx, user.NewStore(h.db))).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get userID from bearer
	userID := auth.GetUserIDFromContext(r.Context())

	// get JSON payload
	var payload types.CartCheckoutPayLoad
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get products
	productIDs, err := getCartItemsIDs(payload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.productStore.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ctx, products, payload.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}

func (h *HandlerTx) handleCheckoutTx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := utils.WithTransaction(ctx, h.db, func(tx *sql.Tx) *utils.MyError {
		orderStoreTx := order.NewStore(tx)
		productStoreTx := product.NewStore(tx)
		userStoreTx := user.NewStore(tx)

		cartHandlerTx := NewHandler(orderStoreTx, productStoreTx, userStoreTx)

		return cartHandlerTx.handleCheckoutInternalTx(ctx, w, r)
	})

	if err != nil {
		utils.WriteError(w, err.StatusCode, err)
	}
}

func (h *Handler) handleCheckoutInternalTx(ctx context.Context, w http.ResponseWriter, r *http.Request) *utils.MyError {
	// get userID from bearer
	userID := auth.GetUserIDFromContext(r.Context())

	// get JSON payload
	var payload types.CartCheckoutPayLoad
	if err := utils.ParseJSON(r, &payload); err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return &utils.MyError{
			Msg:        fmt.Errorf("invalid payload %v", errors).Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// get products
	productIDs, err := getCartItemsIDs(payload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	products, err := h.productStore.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	orderID, totalPrice, err := h.createOrder(ctx, products, payload.Items, userID)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})

	return nil
}
