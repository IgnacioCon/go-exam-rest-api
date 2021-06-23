package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ignaciocon/exam-rest-api/cliente"
	"github.com/ignaciocon/exam-rest-api/config"
	"github.com/ignaciocon/exam-rest-api/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func configurarRutas(app *fiber.App) {
	app.Get("/api/NutriNET/Cliente", cliente.ObtenerClientes)
	app.Get("/api/NutriNET/Cliente/:id", cliente.ObtenerCliente)
	app.Post("/api/NutriNET/Cliente", cliente.NuevoCliente)
	app.Put("/api/NutriNET/Cliente/:id", cliente.ModificarCliente)
	app.Delete("/api/NutriNET/Cliente/:id", cliente.EliminarCliente)
}

func initDatabase() {
	var err error

	database.DBConn, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{})

	if err != nil {
		panic("Error al conectar con la base de datos")
	}

	fmt.Println("Se logro conectar a la base de datos")

	if !database.DBConn.Migrator().HasTable(&cliente.Cliente{}) {
		database.DBConn.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&cliente.Cliente{})
	}

	database.DBConn.AutoMigrate(&cliente.Cliente{})
	fmt.Println("Base de datos migrada")
}

func main() {
	app := fiber.New()
	initDatabase()

	configurarRutas(app)

	app.Listen(":3000")

}
