package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goMongoFiber/src/module"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type m struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errormsg"`
}

func Index(c *fiber.Ctx) error {
	cook := c.Cookies("authentication", "no_cookie")
	// normalnie powinno być wysyłane na podstawie najlepszych ocen (elasticsearch)
	item.Find(context.Background(), bson.D{"", ""})
	if cook == "no_cookie" {
		p := "index"
		fmt.Println(p)
		err := c.Render(p, fiber.Map{
			"Title": module.P.Index,
			"Path":  "/" + p,
			"Role":  32,
		})
		return err
	} else {
		token := jwtToken(c)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role := int(claims["Role"].(float64))
			p := "index"
			fmt.Println(p)
			err := c.Render(p, fiber.Map{
				"Title": module.P.Index,
				"Path":  "/" + p,
				"Role":  role,
			})
			return err
		} else {
			c.Redirect("http://localhost:3000/login")
		}
	}

	return nil
}

//From editor
func AddRecord(c *fiber.Ctx, role int) error {
	var items struct {
		ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Name        string             `json:"name,omitempty" bson:"name,omitempty"`
		Description string             `json:"description,omitempty" bson:"description,omitempty"`
		Price       float32            `json:"price,omitempty" bson:"price,omitempty"`
		Producer    string             `json:"producer,omitempty" bson:"producer,omitempty"`
		Images      []string           `json:"images,omitempty" bson:"images,omitempty"`
	}
	if role <= 8 {
		fmt.Println(role)
		b := c.Request().Body()
		r := bytes.NewReader(b)
		json.NewDecoder(r).Decode(&items)
		// fmt.Println(items)
		if items.Name == "" || items.Producer == "" || items.Description == "" || items.Price == 0 {
			me := m{false, "All field required"}
			c.JSON(&me)
		} else {
			wg.Add(2)
			go func() {
				id := strconv.FormatInt(time.Now().Unix()+rand.Int63(), 16)
				for k, v := range items.Images {
					coI := strings.Index(v, ",")
					vB := []byte(v)
					value := vB[coI+1:]
					sid := fmt.Sprintf("%v_%v.", id, k)
					reader := make(chan *bytes.Reader)
					wg.Add(2)
					go func() {
						unbased, err := base64.StdEncoding.DecodeString(string(value))
						if err != nil {
							panic(err)
						}
						r := bytes.NewReader(unbased)
						reader <- r
						wg.Done()
					}()
					go func() {
						switch strings.TrimSuffix(v[5:coI], ";base64") {
						case "image/png":
							imgPng, _ := png.Decode(<-reader)
							out, _ := os.Create("views/statics/images/" + sid + "png")
							png.Encode(out, imgPng)
						case "image/jpeg":
							imgJpeg, _ := jpeg.Decode(<-reader)
							out, _ := os.Create("views/statics/images/" + sid + "jpeg")
							err := jpeg.Encode(out, imgJpeg, &jpeg.Options{
								Quality: 100,
							})
							if err != nil {
								log.Fatal(err)

							}
						}
						wg.Done()
					}()
					fmt.Println(sid)
					items.Images[k] = sid
				}
				item.InsertOne(context.Background(), &items)
				wg.Wait()
				wg.Done()
			}()
			go func() {
				me := m{true, "Item added"}
				c.JSON(&me)
				wg.Done()
			}()
			wg.Wait()
		}
	} else {
		m := m{Success: false, ErrorMsg: "Permission denied"}
		c.JSON(&m)
	}

	return nil
}

func search() {
	//This function should be connected with elasticsearch to perform search
}
