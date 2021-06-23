CREATE TABLE `clientes` (
  `cliente_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `nombre_usuario` varchar(100),
  `contraseña` varchar(100),
  `nombre` varchar(100),
  `apellidos` varchar(100),
  `correo_electronico` varchar(100),
  `edad` int(5),
  `estatura` decimal(4, 2),
  `peso` decimal(4, 2),
  `imc` decimal(3, 2),
  `geb` decimal(10, 2),
  `eta` decimal(10, 2),
  `fecha_creacion` date,
  `fecha_actualizacion` date,
  PRIMARY KEY (`cliente_id`),UNIQUE INDEX idx_clientes_nombre_usuario (`nombre_usuario`),UNIQUE INDEX idx_clientes_correo_electronico (`correo_electronico`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT;