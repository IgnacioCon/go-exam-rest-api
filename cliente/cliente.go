package cliente

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ignaciocon/exam-rest-api/database"
	"golang.org/x/crypto/bcrypt"
)

//Cliente contiene los campos del modelo a guardar
type Cliente struct {
	ClienteID          int       `json:"Cliente_ID" gorm:"primaryKey"`
	NombreUsuario      string    `json:"Nombre_Usuario" gorm:"uniqueIndex"`
	Contraseña         string    `json:"Contraseña"`
	Nombre             string    `json:"Nombre"`
	Apellidos          string    `json:"Apellidos"`
	CorreoElectronico  string    `json:"Correo_Electronico" gorm:"uniqueIndex"`
	Edad               int       `json:"Edad"`
	Estatura           float64   `json:"Estatura" gorm:"precision:4,2"`
	Peso               float64   `json:"Peso" gorm:"precision:4,2"`
	IMC                float64   `json:"IMC" gorm:"precision:3,2"`
	GEB                float64   `json:"GEB" gorm:"precision:10,2"`
	ETA                float64   `json:"ETA" gorm:"precision:10,2"`
	FechaCreacion      time.Time `json:"Fecha_Creacion" gorm:"autoCreateTime"`
	FechaActualizacion time.Time `json:"Fecha_Actualizacion" gorm:"autoUpdateTime" `
}

//ObtenerClientes Metodo GET ruta "/NutriNET/Cliente" regresa todos los clientes
func ObtenerClientes(c *fiber.Ctx) error {
	db := database.DBConn

	var clientes []Cliente

	db.Find(&clientes)

	respuesta := fiber.Map{}

	if len(clientes) == 0 {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: no se encontraron clientes."
		return c.Status(fiber.StatusNotFound).JSON(respuesta)
	}

	respuesta["Clientes"] = clientes
	respuesta["Cve_Error"] = 0
	respuesta["Cve_Mensaje"] = "Clientes encontrados con exito."

	return c.Status(fiber.StatusOK).JSON(respuesta)
}

//ObtenerCliente Metodo GET ruta "/NutriNet/Cliente/:id" regresa un cliente por ID
func ObtenerCliente(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var cliente Cliente

	db.Find(&cliente, id)

	respuesta := fiber.Map{}

	if cliente.ClienteID == 0 {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: cliente no se encuentra en la base de datos."
		return c.Status(fiber.StatusNotFound).JSON(respuesta)
	}

	respuesta["Cliente"] = cliente
	respuesta["Cve_Error"] = 0
	respuesta["Cve_Mensaje"] = "Cliente encontrado con exito."
	return c.Status(fiber.StatusOK).JSON(respuesta)
}

//NuevoCliente Metodo POST ruta "/api/NutriNET/Cliente" Crear un nuevo cliente con la informacion enviada al API en formato JSON
func NuevoCliente(c *fiber.Ctx) error {
	db := database.DBConn

	var cliente Cliente
	err := c.BodyParser(&cliente)

	respuesta := fiber.Map{}
	if err != nil {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: No se logro leer la informacion."
		return c.Status(fiber.StatusConflict).JSON(respuesta)
	}

	//encriptar contraseña
	encriptada, err := encriptarContraseña(cliente.Contraseña)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//guardar nueva contraseña encriptada
	cliente.Contraseña = encriptada

	resultado := db.Create(&cliente)

	if resultado.Error != nil {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: Nombre de usuario o Correo Electronico ya estan registrados."
		return c.Status(fiber.StatusConflict).JSON(respuesta)
	}

	respuesta["Cliente"] = cliente
	respuesta["Cve_Error"] = 0
	respuesta["Cve_Mensaje"] = "Cliente registrado con exito."

	return c.Status(fiber.StatusCreated).JSON(respuesta)
}

//ModificarCliente metodo PUT ruta '/NutriNET/Cliente/:id' Modifica campos de un usuario por ID
func ModificarCliente(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var datosNuevos Cliente

	err := c.BodyParser(&datosNuevos)

	respuesta := fiber.Map{}
	if err != nil {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: No se logro leer la informacion."
		return c.Status(fiber.ErrBadRequest.Code).JSON(respuesta)
	}

	var cliente Cliente

	db.Find(&cliente, id)

	if cliente.ClienteID == 0 {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: Cliente no se encuentra en la base de datos."
		return c.Status(fiber.StatusNotFound).JSON(respuesta)
	}

	//revisar si hay cambio en la contraseña
	if datosNuevos.Contraseña != "" {
		if compararContraseña(cliente.Contraseña, datosNuevos.Contraseña) {
			respuesta["Cve_Error"] = -1
			respuesta["Cve_Mensaje"] = "Error: Contraseña ingresada deber ser diferente a la anterior."
			return c.Status(fiber.StatusConflict).JSON(respuesta)
		}

		encriptada, err := encriptarContraseña(datosNuevos.Contraseña)

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		datosNuevos.Contraseña = encriptada
	}

	resultado := db.Model(&cliente).Updates(&datosNuevos)

	if resultado.Error != nil {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: Cliente no se actualizo correctamente."
		return c.Status(fiber.StatusNotModified).JSON(respuesta)
	}

	respuesta["Cliente"] = cliente
	respuesta["Cve_Error"] = 0
	respuesta["Cve_Mensaje"] = "Cliente actualizado correctamente."

	return c.Status(fiber.StatusOK).JSON(respuesta)
}

//EliminarCliente Metodo DELETE ruta "/NutriNET/Cliente/:id" Si el usuario con ID existe,
//Se elimina de la base de datos
func EliminarCliente(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var cliente Cliente
	db.Find(&cliente, id)

	respuesta := fiber.Map{}
	if cliente.ClienteID == 0 {
		respuesta["Cve_Error"] = -1
		respuesta["Cve_Mensaje"] = "Error: Ningun cliente con ese ID."
		return c.Status(fiber.StatusNotFound).JSON(respuesta)
	}

	db.Delete(&cliente)

	respuesta["Cve_Error"] = 0
	respuesta["Cve_Mensaje"] = "Cliente eliminado con exito."

	return c.Status(fiber.StatusOK).JSON(respuesta)
}

//encriptar contraseña con un hash y salt
func encriptarContraseña(contraseña string) (string, error) {

	encriptada, err := bcrypt.GenerateFromPassword([]byte(contraseña), bcrypt.MinCost)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(encriptada), nil
}

//CompararContraseña para verificar si la contraseña es la correcta en caso de ingresar
//O para verificar si se ha modificado la constraseña al modificar campos de un cliente
func compararContraseña(contraseñaEncriptada, contraseñaSinEncriptar string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(contraseñaEncriptada), []byte(contraseñaSinEncriptar))
	if err != nil {
		return false
	}

	return true
}
