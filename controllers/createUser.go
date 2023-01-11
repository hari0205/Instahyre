package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ini "example.com/Instahyre/teleapi/init"
	model "example.com/Instahyre/teleapi/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {
	// Getting req body
	var body struct {
		Email    string
		Password string
		Name     string
		Phone_no string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read request body",
		})
		return
	}

	// Hash and store password
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate hashed password",
		})

		return
	}

	//Create user
	userr := model.UserData{Name: body.Name, Password: string(hashedpass), Email: body.Email, Phone_no: body.Phone_no}
	//fmt.Println(userr)
	result := ini.DB.Create(&userr)
	// fmt.Println(result)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERror creating user ",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"name":    body.Name,
	})
}

func Login(c *gin.Context) {

	var body struct {
		Email   string
		Pasword string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read request body",
		})
		return
	}

	var user model.UserData
	ini.DB.First(&user, "email= ? ", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist or Invalid Credentials",
		})
		return
	}

	// Compare Hashed pass with sent pass
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Pasword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password mismatch",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// Generating tokenstring
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authentication", tokenString, 3600*24, "/", " ", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	if user.(model.UserData).ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

func SearchByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	var users []model.UserData
	ini.DB.Select([]string{"Name", "Phone_no", "Is_spam", "Reported_count"}).Find(&users, "name= ?", name)

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Record does not exist",
		})

		return
	}

	// Need to loop through range and display all the users' records for the sake of Convenience display only 1st record
	// for v := range users { // Display all the users}
	c.JSON(http.StatusOK, gin.H{
		"name":        users[0].Name,
		"Ph.no":       users[0].Phone_no,
		"Spam_Report": users[0].Reported_count,
	})
}

func SearchByNumber(c *gin.Context) {
	number := c.Param("number")
	fmt.Println(number)
	if number == " " {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	var user model.UserData
	ini.DB.Select([]string{"Name", "Phone_no", "Is_spam", "Reported_count"}).Find(&user, "Phone_no= ?", number)

	c.JSON(http.StatusOK, gin.H{
		"name":        user.Name,
		"Ph.no":       user.Phone_no,
		"Spam_Report": user.Reported_count,
	})
}

func Report(c *gin.Context) {
	number := c.Param("number")
	fmt.Println(number)
	if number == " " {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	res := ini.DB.Where(model.UserData{Phone_no: number}).Omit("Email").Attrs(model.UserData{Name: "Unknown", Is_spam: true, Reported_count: 1}).FirstOrCreate(&model.UserData{})
	fmt.Println(res.RowsAffected)
	if res.RowsAffected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Record not found. Created first report successfully",
		})

		return
	}

	var user model.UserData
	ini.DB.Model(&user).Select("Is_spam").Where("Phone_no= ?", number).Update("Is_spam", true)
	ini.DB.Model(&user).Select("Reported_count").Where("Phone_no= ?", number).UpdateColumn("Reported_count", gorm.Expr("Reported_count + ?", 1))

	u := model.UserData{}

	result := ini.DB.Select("Reported_count").Where("Phone_no= ?", number).First(&u)

	log.Printf("%v", result)

	c.JSON(http.StatusOK, gin.H{
		"Reported_count": u.Reported_count,
	})
}
