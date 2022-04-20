package main

import (
	"goMongoFiber/src/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("./views", ".html")
	engine.Templates.ParseGlob("./views/partials/*")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./views/statics/scripts/unsecure/")
	app.Static("/", "./views/statics/images/")
	router.Router(app)
	log.Fatal(app.Listen(":3000"))
}

//TODO:

//Trzeba stworzyć stronę, z wybieraniem dowozu (kurier, inpost), wpisanie adresu (przy okazji trzeba się dowiedzieć jak to dobrze zabezpieczyć)
//Podpiąć pay-U
//Elasticsearch and searchbar
//Cart
//Showing mostly buyed items on index page (elasticsearch?) it is needed to create dummy data to items
//Comments
//Front-End
