package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/todo/internal/adapter/http"
	"github.com/example/todo/internal/order"
	"github.com/example/todo/internal/payment"
	"github.com/example/todo/internal/user"
	"github.com/labstack/echo/v4"
)

type fixedIDGenerator struct {
	id string
}

func (g *fixedIDGenerator) New() string {
	return g.id
}

type fixedClock struct {
	now int64
}

func (c *fixedClock) Now() int64 {
	return c.now
}

func main() {
	clock := &fixedClock{now: time.Now().Unix()}
	ids := &fixedIDGenerator{id: "test-id"}

	orderRepo := order.NewMemoryRepo()
	createOrder := order.NewCreateOrder(orderRepo, ids, clock, nil)
	getOrder := order.NewGetOrder(orderRepo)
	getAllOrders := order.NewGetAllOrders(orderRepo)
	updateOrder := order.NewUpdateOrder(orderRepo)
	cancelOrder := order.NewCancelOrder(orderRepo, nil)
	orderHandler := order.NewHandler(createOrder, getOrder, getAllOrders, updateOrder, cancelOrder)

	paymentRepo := payment.NewMemoryRepo()
	createPayment := payment.NewCreatePayment(paymentRepo, ids, clock, nil)
	getPayment := payment.NewGetPayment(paymentRepo)
	getPaymentsByOrder := payment.NewGetPaymentsByOrder(paymentRepo)
	getAllPayments := payment.NewGetAllPayments(paymentRepo)
	updatePayment := payment.NewUpdatePayment(paymentRepo)
	refundPayment := payment.NewRefundPayment(paymentRepo, nil)
	paymentHandler := payment.NewHandler(
		createPayment, getPayment, getPaymentsByOrder, getAllPayments, updatePayment, refundPayment,
	)

	userRepo := user.NewMemoryRepo()
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	server := http.NewServer()
	server.RegisterRoutes(func(api *echo.Group) {
		orderHandler.RegisterRoutes(api.Group("/orders"))
		paymentHandler.RegisterRoutes(api.Group("/payments"))
		userHandler.RegisterRoutes(api.Group("/users"))
	})

	port := getEnv("PORT", "8080")
	addr := ":" + port

	log.Printf("Server started on %s", addr)

	go func() {
		if err := server.Start(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}