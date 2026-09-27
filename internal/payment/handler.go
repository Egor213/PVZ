package payment

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	createPayment     *CreatePayment
	getPayment        *GetPayment
	getPaymentsByOrder *GetPaymentsByOrder
	getAllPayments    *GetAllPayments
	updatePayment     *UpdatePayment
	refundPayment     *RefundPayment
}

func NewHandler(
	createPayment *CreatePayment,
	getPayment *GetPayment,
	getPaymentsByOrder *GetPaymentsByOrder,
	getAllPayments *GetAllPayments,
	updatePayment *UpdatePayment,
	refundPayment *RefundPayment,
) *Handler {
	return &Handler{
		createPayment:      createPayment,
		getPayment:         getPayment,
		getPaymentsByOrder: getPaymentsByOrder,
		getAllPayments:     getAllPayments,
		updatePayment:      updatePayment,
		refundPayment:      refundPayment,
	}
}

type createReq struct {
	OrderID string `json:"order_id" validate:"required"`
	Amount  int64  `json:"amount" validate:"required,gt=0"`
}

type updateReq struct {
	Status *Status `json:"status,omitempty"`
}

type paymentResp struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	Amount    int64  `json:"amount"`
	Status    Status `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func toResp(p interface{}) paymentResp {
	switch v := p.(type) {
	case CreatePaymentOutput:
		return paymentResp{ID: v.ID, OrderID: v.OrderID, Amount: v.Amount, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	case GetPaymentOutput:
		return paymentResp{ID: v.ID, OrderID: v.OrderID, Amount: v.Amount, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	case UpdatePaymentOutput:
		return paymentResp{ID: v.ID, OrderID: v.OrderID, Amount: v.Amount, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	case RefundPaymentOutput:
		return paymentResp{ID: v.ID, OrderID: v.OrderID, Amount: v.Amount, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	}
	return paymentResp{}
}

func (h *Handler) Create(c echo.Context) error {
	var req createReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	out, err := h.createPayment.Execute(c.Request().Context(), CreatePaymentInput{
		OrderID: req.OrderID,
		Amount:  req.Amount,
	})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusCreated, toResp(out))
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	out, err := h.getPayment.Execute(c.Request().Context(), GetPaymentInput{ID: id})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(out))
}

func (h *Handler) GetByOrderID(c echo.Context) error {
	orderID := c.Param("order_id")
	out, err := h.getPaymentsByOrder.Execute(c.Request().Context(), GetPaymentsByOrderInput{OrderID: orderID})
	if err != nil {
		return h.handleErr(c, err)
	}
	resp := make([]paymentResp, len(out.Payments))
	for i, p := range out.Payments {
		resp[i] = toResp(p)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAll(c echo.Context) error {
	out, err := h.getAllPayments.Execute(c.Request().Context(), GetAllPaymentsInput{})
	if err != nil {
		return h.handleErr(c, err)
	}
	resp := make([]paymentResp, len(out.Payments))
	for i, p := range out.Payments {
		resp[i] = toResp(p)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req updateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	out, err := h.updatePayment.Execute(c.Request().Context(), UpdatePaymentInput{
		ID:     id,
		Status: req.Status,
	})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(out))
}

func (h *Handler) Refund(c echo.Context) error {
	id := c.Param("id")
	out, err := h.refundPayment.Execute(c.Request().Context(), RefundPaymentInput{ID: id})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(out))
}

func (h *Handler) handleErr(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrInvalidAmount):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidStatus):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrCannotRefund):
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.GetAll)
	g.GET("/:id", h.Get)
	g.GET("/order/:order_id", h.GetByOrderID)
	g.PUT("/:id", h.Update)
	g.POST("/:id/refund", h.Refund)
}