package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграция схемы
	db.AutoMigrate(&Sale{})
}

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
	var flag bool = false
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "in...."})
	}

	for _, name := range users {
		if creds.Username == name.Username && creds.Password == name.Password {
			flag = true
			token, err := generateToken(creds.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
			return

		}

	}
	if !flag {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
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

type Sale struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	ARTIKUL string `json:"artikul tovara"`
	OTDEL   string `json:"name otdela "`
	DATE    string `json:"date sale"`
	COUNT   string `json:"kol-vo sale"`
}

func getSales(c *gin.Context) {
	var sales []Sale
	db.Find(&sales)
	c.JSON(http.StatusOK, sales)
}

func getSaleByID(c *gin.Context) {
	id := c.Param("id")
	var sale Sale
	if err := db.First(&sale, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "sale not found"})
		return
	}
	c.JSON(http.StatusOK, sale)
}
func createSale(c *gin.Context) {
	var newSale Sale
	if err := c.BindJSON(&newSale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	db.Create(&newSale)
	c.JSON(http.StatusCreated, newSale)
}

func updateSale(c *gin.Context) {
	id := c.Param("id")
	var updateSale Sale
	if err := c.BindJSON(&updateSale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if err := db.Model(&Sale{}).Where("id = ?", id).Updates(updateSale).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "sale not found"})
		return
	}
	c.JSON(http.StatusOK, updateSale)
}

func deleteSale(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Sale{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "sale not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sale deleted"})
}

func main() {
	initDB()
	router := gin.Default()
	router.POST("/login", login)

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/sales", getSales)
		protected.POST("/sales", createSale)
		// другие защищенные маршруты
	}

	//router.GET("/sales", getSales)
	router.GET("/sales/:id", getSaleByID)
	//router.POST("/sales", createSale)
	router.PUT("/sales/:id", updateSale)
	router.DELETE("/sales/:id", deleteSale)

	router.Run(":8080")
}
