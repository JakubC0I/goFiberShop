package module

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var P = struct {
	Index        string
	Dummy        string
	AddRecord    string
	Login        string
	Register     string
	Secret       string
	Product      string
	AddComment   string
	ViewComments string
	Cart         string
	Deliv        string
	Deliver      string
}{
	Index:        "",
	Dummy:        "dummy",
	AddRecord:    "addRecord",
	Login:        "login",
	Register:     "register",
	Secret:       "secret",
	Product:      "product",
	AddComment:   "addComment",
	ViewComments: "viewComments",
	Cart:         "cart",
	Deliv:        "deliv",
	Deliver:      "deliver",
}

type Item struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	HexID       string               `json:"hexID,omitempty" bson:"hexID,omitempty"`
	Name        string               `json:"name" bson:"name" validate:"required"`
	Price       float32              `json:"price" bson:"price" validate:"required"`
	Comments    []primitive.ObjectID `json:"comments,omitempty" bson:"comments,omitempty"`
	Rating      float32              `json:"rating,omitempty" bson:"rating,omitempty"`
	NumBuyed    int                  `json:"numbuyed,omitempty" bson:"numbuyed,omitempty"`
	Description string               `json:"description" bson:"description"`
	Producer    string               `json:"producer" bson:"producer"`
	Images      []string             `json:"images,omitempty" bson:"images,omitempty"`
	Quantity    int                  `json:"quantity" bson:"quantity"`
}
type Comment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Body      string             `json:"body" bson:"body"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

const (
	Owner int = 1 << iota
	Admin
	Moderator
	Editor
	Viewer //zalogowany user bez dodatkowych uprawnień
	Unregistered
)

type Account struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Username   string             `json:"username" bson:"username"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	Role       int                `json:"role" bson:"role"`
	Created_at time.Time          `json:"createdAt" bson:"createdAt"`
}

type Error struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errormessage"`
}
