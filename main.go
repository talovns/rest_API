package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// добавлено отсюда
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []Credentials{
	{Username: "bob", Password: "123"},
	{Username: "dog", Password: "223"},
	{Username: "gob", Password: "113"},
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func login(c *gin.Context) {
	var creds Credentials

	for _, name := range users {
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusOK, nil)
		}
		if name.Username == creds.Username {

			// Здесь добавим простую проверку пароля

			if creds.Username != name.Username || creds.Password != name.Password {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}

			token, err := generateToken(creds.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		} else {
			continue
		}
	}
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// до сюда
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
	router.POST("/login", login)

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/sales", getSales)
		protected.POST("/sales", createSale)
		// другие защищенные маршруты
	}
	// Получение всех книгx
	//router.GET("/sales", getSales)

	// Получение книги по ID
	router.GET("/sales/:id", getSaleByID)

	// Создание новой книги
	//router.POST("/sales", createSale)

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
