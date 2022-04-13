package controller

import (
	"context"
	"fmt"
	"goMongoFiber/src/secrets"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbString string = "mongodb://admin:123456@localhost:27017/?maxPoolSize=20&w=majority"
var user *mongo.Collection
var item *mongo.Collection

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbString))
	if err != nil {
		panic(err)
	}
	//już działa
	user = client.Database("goFiber").Collection("User")
	item = client.Database("goFiber").Collection("Item")
}

func SendHTML(c *fiber.Ctx) error {
	p := c.Path()
	parm := c.Params("*")

	re := regexp.MustCompile("^/.*/")
	path := string(re.Find([]byte(p)))
	if path == "" {
		path = c.Path()
	}
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}
	fmt.Println(p, parm, path)
	err := c.Render(path[1:], fiber.Map{
		"Title": strings.ToUpper(path[1:]),
		"Path":  path,
		"Role":  32,
	})
	return err
}

func SendHTMLroles(c *fiber.Ctx, role int) error {
	p := c.Path()
	parm := c.Params("*")
	fmt.Println(role)
	re := regexp.MustCompile("^/.*/")
	path := string(re.Find([]byte(p)))
	if path == "" {
		path = c.Path()
	}
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}
	fmt.Println(p, parm, path)
	err := c.Render(path[1:], fiber.Map{
		"Title": strings.ToUpper(path[1:]),
		"Path":  path,
		"Role":  role,
	})
	return err
}

//Trzeba jeszcze sprawdzić hasło przez bcrypt
func createJWT(c *fiber.Ctx, userID primitive.ObjectID, role int) error {
	type MyClaims struct {
		Role int
		jwt.StandardClaims
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + int64(time.Hour)*3,
			Issuer:    userID.Hex(),
		},
	})
	token, err := claims.SignedString([]byte(secrets.Secrets.SignedKey))
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "authentication",
		HTTPOnly: true,
		Value:    token,
	})
	return nil
}

func jwtToken(c *fiber.Ctx) *jwt.Token {
	var token jwt.Token
	cookie := c.Cookies("authentication", "no_cookie")
	if cookie == "no_cookie" {
		c.Redirect("http://localhost:3000/login")
	} else {
		tkn, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secrets.Secrets.SignedKey), nil
		})
		if err != nil {
			c.Redirect("http://localhost:3000/login")
			fmt.Println(err)
		}
		token = *tkn
	}
	return &token
}
