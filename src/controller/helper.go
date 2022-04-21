package controller

import (
	"context"
	"fmt"
	"goMongoFiber/src/module"
	"goMongoFiber/src/secrets"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbString string = "mongodb+srv://jakub_user:ck3YJrce9rtuPRdj@cluster0.vdtkf.mongodb.net/ChromebookDB?retryWrites=true&w=majority"

// var dbString string = "mongodb://admin:123456@localhost:27017/?maxPoolSize=20&w=majority"
var user *mongo.Collection
var item *mongo.Collection
var comment *mongo.Collection

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbString))
	if err != nil {
		panic(err)
	}
	//już działa
	user = client.Database("goFiber").Collection("User")
	item = client.Database("goFiber").Collection("Item")
	comment = client.Database("goFiber").Collection("Comment")
}

//Trzeba jeszcze sprawdzić hasło przez bcrypt
func createJWT(c *fiber.Ctx, userID primitive.ObjectID, username string, role int) error {
	type MyClaims struct {
		Role     int
		Username string
		jwt.StandardClaims
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaims{
		Role:     role,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + int64(time.Hour)*3,
			Issuer:    userID.Hex(),
		},
	})
	token, err := claims.SignedString([]byte(secrets.Secrets.SignedKey))
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "authentication",
		HTTPOnly: true,
		Value:    token,
		MaxAge:   3600 * 3,
	})
	return nil
}

func jwtToken(c *fiber.Ctx) *jwt.Token {
	var token jwt.Token
	cookie := c.Cookies("authentication", "no_cookie")
	if cookie == "no_cookie" {
		c.Redirect(module.Address + "/login")
	} else {
		tkn, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secrets.Secrets.SignedKey), nil
		})
		if err != nil {
			c.Redirect(module.Address + "/login")
			fmt.Println(err)
		}
		token = *tkn
	}
	return &token
}

func top(c *fiber.Ctx, ch chan []module.Item) {
	var results []module.Item
	cur, err := item.Find(context.Background(), bson.D{},
		options.Find().SetSort(bson.D{{Key: "rating", Value: -1}, {Key: "_id", Value: 1}}).SetLimit(20))

	for cur.Next(context.Background()) {
		var res module.Item
		err := cur.Decode(&res)
		if err != nil {
			panic(err)
		}
		res.HexID = res.ID.Hex()
		results = append(results, res)
	}
	if err != nil {
		fmt.Println(err)
	}
	ch <- results
	wg.Done()
}

func html(c *fiber.Ctx) string {
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
	return path
}
