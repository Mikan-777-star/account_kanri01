package main

import (
	"fmt"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

var db *gorm.DB
var err error

type User struct{
	gorm.Model
	Username string `gorm: "unique"`
	Password string
	Email    string `gorm: "unique"`
}
func initDB(){
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err =gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
}

func register(c *gin.Context){
	var input struct{
		Username string `json:"username"`
        Password string `json:"password"`
        Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
	user := User{Username: input.Username, Password: string(hashPassword), Email: input.Email}
	result := db.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func login(c *gin.Context){
	var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
		return
	}
    var user User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
	}
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

} 

func main() {
    initDB()

    router := gin.Default()
    router.POST("/register", register)
    router.POST("/login", login)
    router.Run(":8080")
}