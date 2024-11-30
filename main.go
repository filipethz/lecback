package main

import (
	"lechrome/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Habilitar CORS para permitir requisições de qualquer origem (pode ser ajustado conforme necessário)
	e.Use(middleware.CORS())

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Rotas
	e.POST("/purchase", handler.CreatePurchase)
	e.GET("/purchases", handler.GetPurchases)

	// Iniciar o servidor
	e.Logger.Fatal(e.Start(":8080"))
}
