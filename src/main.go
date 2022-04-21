package main

import (
	"fmt"
	"goMongoFiber/src/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

//Blokada na spamowanie Register
func main() {
	dir, _ := os.Getwd()
	engine := html.New(dir+"/src/views", ".html")
	engine.Templates.ParseGlob("./views/partials/*")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	app.Static("/", dir+"/src/views/statics/scripts/unsecure/")
	app.Static("/", dir+"/src/views/statics/images/")
	router.Router(app)
	log.Fatal(app.Listen(addr))
}
func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

//TODO:

//Trzeba stworzyć stronę, z wybieraniem dowozu (kurier, inpost), wpisanie adresu (przy okazji trzeba się dowiedzieć jak to dobrze zabezpieczyć)
//Podpiąć pay-U
//Elasticsearch and searchbar
//Cart
//Showing mostly buyed items on index page (elasticsearch?) it is needed to create dummy data to items
//Comments
//Front-End
