package router

import (
	"goMongoFiber/src/controller"
	"goMongoFiber/src/module"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {

	app.Get("/", controller.Index)

	app.Get("/"+module.P.Dummy+"/*", controller.SendHTML)

	app.Get("/scripts/secure/*", controller.IsLoggedInWithRoles(controller.SecureJS))

	app.Get("/"+module.P.Login, controller.SendHTML)
	app.Post("/"+module.P.Login, controller.Login)

	app.Get("/"+module.P.Register, controller.SendHTML)
	app.Post("/"+module.P.Register, controller.Register)

	app.Get("/"+module.P.AddRecord, controller.IsLoggedInWithRoles(controller.SendHTMLroles))
	app.Post("/"+module.P.AddRecord, controller.IsLoggedInWithRoles(controller.AddRecord))
	//dodać odpowiednie wyświetlanie danego produktu
	app.Get("/"+module.P.Product+"/:id", controller.NoVerfWithRoles(controller.Product))

	app.Post("/"+module.P.AddComment+"/:id", controller.IsLoggedInWithRoles(controller.AddComment))
	app.Get("/"+module.P.ViewComments+"/:id", controller.ViewComments)

	app.Get("/"+module.P.Cart, controller.SendHTML)
}
