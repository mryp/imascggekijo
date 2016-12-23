package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	//ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "imascggekijo v1.0")
	})
	apiGroup := e.Group("/api")
	apiGroup.POST("/regist", RegistHandler)
	apiGroup.POST("/select", SelectHandler)
	apiGroup.GET("/image/:number", ImageHandler)

	//開始
	e.Logger.Fatal(e.Start(":8001"))
}
