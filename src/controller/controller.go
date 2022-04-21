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
	"math"
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
	resu := make(chan []module.Item)
	wg.Add(1)
	go top(c, resu)
	if cook == "no_cookie" {
		p := "index"
		fmt.Println(p)
		err := c.Render(p, fiber.Map{
			"Title": module.P.Index,
			"Path":  "/" + p,
			"Role":  32,
			"Top20": <-resu,
		})
		return err
	} else {
		token := jwtToken(c)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role := int(claims["Role"].(float64))
			username := claims["Username"].(string)
			p := "index"
			fmt.Println(p)
			err := c.Render(p, fiber.Map{
				"Title":    module.P.Index,
				"Path":     "/" + p,
				"Role":     role,
				"Username": username,
				"Top20":    <-resu,
			})
			return err
		} else {
			c.Redirect("http://localhost:3000/login")
		}
	}
	wg.Wait()
	return nil

}

//From editor
func AddRecord(c *fiber.Ctx, role int, user string, username string) error {
	// var items struct {
	// ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	// Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	// Description string             `json:"description,omitempty" bson:"description,omitempty"`
	// Price       float32            `json:"price,omitempty" bson:"price,omitempty"`
	// Producer    string             `json:"producer,omitempty" bson:"producer,omitempty"`
	// Images      []string           `json:"images,omitempty" bson:"images,omitempty"`
	// }
	var items module.Item
	if role <= 8 {
		fmt.Println(role)
		b := c.Request().Body()
		r := bytes.NewReader(b)
		json.NewDecoder(r).Decode(&items)
		// fmt.Println(items)
		items.Rating = 0
		if items.Name == "" || items.Producer == "" || items.Description == "" || items.Price == 0 {
			me := m{false, "All fields required"}
			c.JSON(&me)
		} else {
			wg.Add(2)
			go func() {
				id := strconv.FormatInt(time.Now().Unix()+rand.Int63(), 16)
				for k, v := range items.Images {
					name := make(chan string)
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
							sid = sid + "png"
							name <- sid
							out, _ := os.Create("views/statics/images/" + sid)
							png.Encode(out, imgPng)
						case "image/jpeg":
							imgJpeg, _ := jpeg.Decode(<-reader)
							sid = sid + "jpeg"
							name <- sid
							out, _ := os.Create("views/statics/images/" + sid)
							err := jpeg.Encode(out, imgJpeg, &jpeg.Options{
								Quality: 100,
							})
							if err != nil {
								log.Fatal(err)

							}
						}
						wg.Done()
					}()
					items.Images[k] = <-name
				}
				math.Round(float64(items.Price))
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

func SendHTML(c *fiber.Ctx) error {
	cook := c.Cookies("authentication", "no_cookie")
	p := html(c)
	// normalnie powinno być wysyłane na podstawie najlepszych ocen (elasticsearch)
	if cook == "no_cookie" {
		fmt.Println(p)
		err := c.Render(p[1:], fiber.Map{
			"Title": module.P.Index,
			"Path":  p,
			"Role":  32,
		})
		return err
	} else {
		token := jwtToken(c)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role := int(claims["Role"].(float64))
			username := claims["Username"].(string)
			fmt.Println(p)
			err := c.Render(p[1:], fiber.Map{
				"Title":    module.P.Index,
				"Path":     p,
				"Role":     role,
				"Username": username,
			})
			return err
		} else {
			c.Redirect("http://localhost:3000/login")
		}
	}
	return nil

}

func SendHTMLroles(c *fiber.Ctx, role int, user string, username string) error {
	path := html(c)
	err := c.Render(path[1:], fiber.Map{
		"Title":    strings.ToUpper(path[1:]),
		"Path":     path,
		"Role":     role,
		"Username": username,
		"User":     user,
	})
	return err
}

func Product(c *fiber.Ctx, role int, uID string, username string) error {
	var items module.Item
	param := c.Params("id")
	fmt.Println(param)
	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		fmt.Println(err)
	}
	resu := item.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	resu.Decode(&items)
	items.HexID = items.ID.Hex()
	path := html(c)
	err = c.Render(path[1:], fiber.Map{
		"Title":    strings.ToUpper(path[1:]),
		"Path":     path,
		"Role":     role,
		"Username": username,
		"Result":   items,
	})
	return err
}

func AddComment(c *fiber.Ctx, role int, user string, username string) error {
	if role <= 16 {
		id := c.Params("id")
		var comm module.Comment
		c.BodyParser(&comm)
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println(err)
		}
		if err != nil {
			fmt.Println(err)
		}
		comm.CreatedAt = time.Now()
		comm.Username = username

		resu, err := comment.InsertOne(context.Background(), &comm)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		resuID := resu.InsertedID.(primitive.ObjectID)

		item.UpdateOne(context.Background(), bson.M{"_id": oid}, bson.M{"$addToSet": bson.M{"comments": resuID}})
		m := module.Error{Success: true, ErrorMessage: "Comment has been added!"}
		c.JSON(&m)
	} else {
		m := module.Error{Success: false, ErrorMessage: "Please login to create a comment"}
		c.JSON(&m)
	}
	return nil
}

func ViewComments(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(err.Error())
	}
	// s := c.Params("skip")
	// skip, err := strconv.Atoi(s)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// l := c.Params("limit")
	// limit, err := strconv.Atoi(l)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	var it module.Item
	res := item.FindOne(context.Background(), bson.M{"_id": oid})
	res.Decode(&it)
	if it.Comments == nil {
		m := module.Error{
			Success:      false,
			ErrorMessage: "No comments",
		}
		c.JSON(&m)
	} else {
		var commentIds []primitive.ObjectID
		commentIds = append(commentIds, it.Comments...)
		fmt.Println(commentIds)
		results, err := comment.Find(context.Background(), bson.M{"_id": bson.M{"$in": commentIds}})
		if err != nil {
			fmt.Println(err)
		}

		var comm module.Comment
		var comms []module.Comment

		for results.Next(context.Background()) {
			results.Decode(&comm)
			comms = append(comms, comm)
		}
		cs := struct {
			Comments []module.Comment `json:"comments"`
		}{Comments: comms}
		fmt.Println(cs)
		c.JSON(&cs)
	}
	return nil
}

//Tutaj należałoby przekierować użytkownika do strony payU, odebrać odpowiedź i wysłać adres do Inpost/DPD/DHL/Poczta
func Deliver(c *fiber.Ctx) error {
	//Przy takim serwisie jak allegro powinien być jeszcze widoczny seller aby wysłać do odpowiedniego konta zamówienie
	//W naszym wypadku można stworzyć mapę producentów (Enum) i na podstawie tego wysyłać requesty na odpowiednie ścieżki
	//Wszystko zależy od modelu sprzedaży sklepu

	var body map[string]interface{}
	c.BodyParser(&body)
	var quantity []int
	var price []float64
	var id []primitive.ObjectID
	var bill float64
	for _, v := range body["products"].([]interface{}) {
		value, _ := v.(map[string]interface{})
		q := value["quantity"].(float64)
		p := math.Round(float64(value["price"].(float64)*100)) / 100
		quantity = append(quantity, int(q))
		price = append(price, p)
		r, err := primitive.ObjectIDFromHex(value["product"].(string))
		if err != nil {
			fmt.Println(err)
		}
		id = append(id, r)
		bill += q * p
	}

	resu, err := item.Find(context.Background(), bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resu.Current.Elements())
	var it module.Item
	var its []module.Item
	for resu.Next(context.Background()) {
		resu.Decode(&it)
		its = append(its, it)
	}
	fmt.Println(bill, its)

	//W tym miejscu podłączyć PayU, odnotować płatność, po czym wysłać JSON "do magazynu"

	return nil
}
