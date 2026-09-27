package order

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	createOrder   *CreateOrder
	getOrder      *GetOrder
	getAllOrders  *GetAllOrders
	updateOrder   *UpdateOrder
	cancelOrder   *CancelOrder
}

func NewHandler(
	createOrder *CreateOrder,
	getOrder *GetOrder,
	getAllOrders *GetAllOrders,
	updateOrder *UpdateOrder,
	cancelOrder *CancelOrder,
) *Handler {
	return &Handler{
		createOrder:  createOrder,
		getOrder:     getOrder,
		getAllOrders: getAllOrders,
		updateOrder:  updateOrder,
		cancelOrder:  cancelOrder,
	}
}

type createReq struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description"`
}

type updateReq struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,max=200"`
	Description *string `json:"description,omitempty"`
	Status      *Status `json:"status,omitempty"`
}

type orderResp struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func toResp(o interface{}) orderResp {
	switch v := o.(type) {
	case CreateOrderOutput:
		return orderResp{ID: v.ID, Title: v.Title, Description: v.Description, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	case GetOrderOutput:
		return orderResp{ID: v.ID, Title: v.Title, Description: v.Description, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	case UpdateOrderOutput:
		return orderResp{ID: v.ID, Title: v.Title, Description: v.Description, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	case CancelOrderOutput:
		return orderResp{ID: v.ID, Title: "", Description: "", Status: v.Status, CreatedAt: 0, UpdatedAt: v.UpdatedAt}
	}
	return orderResp{}
}

func (h *Handler) Create(c echo.Context) error {
	var req createReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	out, err := h.createOrder.Execute(c.Request().Context(), CreateOrderInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusCreated, toResp(out))
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	out, err := h.getOrder.Execute(c.Request().Context(), GetOrderInput{ID: id})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(out))
}

func (h *Handler) GetAll(c echo.Context) error {
	out, err := h.getAllOrders.Execute(c.Request().Context(), GetAllOrdersInput{})
	if err != nil {
		return h.handleErr(c, err)
	}
	resp := make([]orderResp, len(out.Orders))
	for i, o := range out.Orders {
		resp[i] = toResp(o)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req updateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	out, err := h.updateOrder.Execute(c.Request().Context(), UpdateOrderInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(out))
}

func (h *Handler) Cancel(c echo.Context) error {
	id := c.Param("id")
	out, err := h.cancelOrder.Execute(c.Request().Context(), CancelOrderInput{ID: id})
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(out))
}

func (h *Handler) handleErr(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrEmptyTitle):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrTitleTooLong):
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidStatus):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrCannotCancel):
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
	g.PUT("/:id", h.Update)
	g.POST("/:id/cancel", h.Cancel)
}