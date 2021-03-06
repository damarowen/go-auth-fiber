package controller

import (
	"auth-go-fiber/database"
	"auth-go-fiber/helpers"
	"auth-go-fiber/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func Home(c *fiber.Ctx) error {

	//TODO
	//ambil dari bearer token

	//* get token from cookies
	cookie := c.Cookies("jwt")

	//* parsing token
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Register(c *fiber.Ctx) error {

	//var data map[string]string
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := helpers.ValidateStruct(*user)
	if errors != nil {
		return c.JSON(errors)
	}

	fmt.Println(*user)

	//* convert payload password to byte
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.JSON(errors)
	}

	//* convert payload password from byte to string
	save := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(password),
	}

	fmt.Println(save)

	database.DB.Create(&save)

	return c.JSON(save)
}

func Login(c *fiber.Ctx) error {
	//var data map[string]string
	payload := new(models.CheckUser)

	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	errors := helpers.ValidateStruct(*payload)
	if errors != nil {
		return c.JSON(errors)
	}
	var user models.User

	//* search user
	database.DB.Where("email = ?", payload.Email).First(&user)

	//* if not fount
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		//* convert user id to string
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success login",
		"user":    user,
	})
}

func User(c *fiber.Ctx) error {
	//* get token from cookies
	cookie := c.Cookies("jwt")

	//* parsing token
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
