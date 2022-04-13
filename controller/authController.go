package controller

import (
	"context"
	"fmt"
	"goMongoFiber/src/module"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var b = struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	Role     int    `json:"role,omitempty"`
}{}

var wg sync.WaitGroup

func Register(c *fiber.Ctx) error {
	b := b
	err := c.BodyParser(&b)
	if err != nil {
		panic(err)
	}
	hashByte := make(chan []byte)
	wg.Add(2)
	go func() {
		hash, err := bcrypt.GenerateFromPassword([]byte(b.Password), 12)
		if err != nil {
			panic(err)
		}
		hashByte <- hash
		wg.Done()
	}()
	go func() {
		b.Role = 16
		inserted, err := user.InsertOne(context.Background(), bson.D{
			{Key: "email", Value: b.Email},
			{Key: "username", Value: b.Username},
			{Key: "role", Value: b.Role},
			{Key: "createdAt", Value: time.Now()},
			{Key: "password", Value: string(<-hashByte)},
		})
		fmt.Println(inserted.InsertedID)
		if err != nil {
			panic(err)
		}
		close(hashByte)
		wg.Done()
	}()
	wg.Wait()
	return err
}

func Login(c *fiber.Ctx) error {
	b := b
	b.Username = ""
	var a module.Account
	aPass := make(chan string)
	bPass := make(chan string)

	wg.Add(2)
	go func() {
		err := c.BodyParser(&b)
		if err != nil {
			panic(err)
		}
		res := user.FindOne(context.Background(), bson.M{"email": b.Email})
		erro := res.Decode(&a)
		aPass <- a.Password
		bPass <- b.Password
		if erro != nil {
			panic(erro)
		}
		wg.Done()
	}()

	go func() {
		erro := bcrypt.CompareHashAndPassword([]byte(<-aPass), []byte(<-bPass))
		if erro != nil {
			panic(erro)
		} else {
			err := createJWT(c, a.ID, a.Role)
			if err != nil {
				errMsg := module.Error{Success: false, ErrorMessage: err.Error()}
				c.JSON(&errMsg)
			}
			status := module.Error{Success: true, ErrorMessage: "Logged in successfully"}
			c.JSON(&status)
			wg.Done()
		}
	}()
	wg.Wait()
	return nil
}

func IsLoggedIn(next func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := jwtToken(c)
		if token.Valid {
			next(c)
		} else {
			c.Redirect("http://localhost:3000/login")
		}
		return nil
	}
}

func IsLoggedInWithRoles(next func(c *fiber.Ctx, r int) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := jwtToken(c)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role := int(claims["Role"].(float64))
			next(c, role)
		} else {
			c.Redirect("http://localhost:3000/login")
		}
		return nil
	}
}

//Jak będą variable w link to trzeba będzie zrobić regexp
func SecureJS(c *fiber.Ctx, role int) error {
	if role <= 8 {
		p := c.Path()
		r, _ := os.Getwd()
		c.SendFile(r + "/views/statics" + p)
		fmt.Println(r + "/views/statics" + p)
	} else {
		m := m{Success: false, ErrorMsg: "Not authorized"}
		c.Status(fiber.StatusNetworkAuthenticationRequired)
		c.JSON(m)
	}
	return nil
}
