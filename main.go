package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Sale struct {
	ID        string `json:"id"`
	GOODSNAME string `json:"goodsname"`
	PRICE     string `json:"price"`
	CURRENCY  string `json:"currency"`
}

var sales = []Sale{
	{ID: "1", GOODSNAME: "drug", PRICE: "100000000000", CURRENCY: "RUB"},
	{ID: "2", GOODSNAME: "doll", PRICE: "15", CURRENCY: "USD"},
	{ID: "3", GOODSNAME: "helicopter", PRICE: "3000000", CURRENCY: "EUR"},
}

func main() {
	router := gin.Default()

	// Получение всех книг
	router.GET("/sales", getSales)

	// Получение книги по ID
	router.GET("/sales/:id", getSaleByID)

	// Создание новой книги
	router.POST("/sales", createSale)

	// Обновление существующей книги
	router.PUT("/sales/:id", updateSale)

	// Удаление книги
	router.DELETE("/sales/:id", deleteSale)

	router.Run(":8080")
}
func getSales(c *gin.Context) {
	c.JSON(http.StatusOK, sales)
}

func getSaleByID(c *gin.Context) {
	id := c.Param("id")

	for _, sale := range sales {
		if sale.ID == id {
			c.JSON(http.StatusOK, sale)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "sale not found"})
}

func createSale(c *gin.Context) {
	var newSale Sale

	if err := c.BindJSON(&newSale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	sales = append(sales, newSale)
	c.JSON(http.StatusCreated, newSale)
}
func updateSale(c *gin.Context) {
	id := c.Param("id")
	var updatedsale Sale

	if err := c.BindJSON(&updatedsale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	for i, sale := range sales {
		if sale.ID == id {
			sales[i] = updatedsale
			c.JSON(http.StatusOK, updatedsale)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "sale not found"})
}

func deleteSale(c *gin.Context) {
	id := c.Param("id")

	for i, sale := range sales {
		if sale.ID == id {
			sales = append(sales[:i], sales[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "sale deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "sale not found"})
}
