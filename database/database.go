package database

import (
	"gorm.io/gorm"
)

//DBConn conexion global para utilizar en la aplicacion
var DBConn *gorm.DB
