package user

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createReq struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

type updateReq struct {
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Name  *string `json:"name,omitempty"`
}

type userResp struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func toResp(u *User) userResp {
	return userResp{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (h *Handler) Create(c echo.Context) error {
	var req createReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	u, err := h.svc.Create(c.Request().Context(), req.Email, req.Name)
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusCreated, toResp(u))
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	u, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(u))
}

func (h *Handler) GetByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email required"})
	}
	u, err := h.svc.GetByEmail(c.Request().Context(), email)
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(u))
}

func (h *Handler) GetAll(c echo.Context) error {
	users, err := h.svc.GetAll(c.Request().Context())
	if err != nil {
		return h.handleErr(c, err)
	}
	resp := make([]userResp, len(users))
	for i, u := range users {
		resp[i] = toResp(u)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req updateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	in := UpdateInput{Email: req.Email, Name: req.Name}
	u, err := h.svc.Update(c.Request().Context(), id, in)
	if err != nil {
		return h.handleErr(c, err)
	}
	return c.JSON(http.StatusOK, toResp(u))
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return h.handleErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) handleErr(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrInvalidEmail):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrEmailExists):
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
	g.GET("/by-email", h.GetByEmail)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}