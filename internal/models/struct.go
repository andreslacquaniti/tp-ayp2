package models

import (
	"time"
)

// AeronaveStruc representa datos de "aeronaves.txt"
type AeronaveStruc struct {
	// Ejemplo: Matricula (e.g., LV085)
	Matricula string `csv:"Matricula"`
	// Ejemplo: Asientos (e.g., 180)
	Asientos int `csv:"Asientos"`
	// Ejemplo: CapacidadCarga (e.g., 2000)
	CapacidadCarga float64 `csv:"CapacidadCarga"`
	// Ejemplo: VolumenCarga (e.g., 100)
	VolumenCarga float64 `csv:"VolumenCarga"`
}

// AeropuertoStruc representa datos de "aeropuertos.txt"
type AeropuertoStruc struct {
	// Ejemplo: PROVINCIA (e.g., Buenos Aires)
	Provincia string `csv:"PROVINCIA"`
	// Ejemplo: CIUDAD (e.g., Bahía Blanca)
	Ciudad string `csv:"CIUDAD"`
	// Ejemplo: NOMBRE AEROPUERTO (e.g., Comandante Espora)
	NombreAeropuerto string `csv:"NOMBRE AEROPUERTO"`
	// Ejemplo: COD_IATA (e.g., BHI)
	CodIATA string `csv:"COD_IATA"`
}

// CargaStruc representa datos de "cargas.txt"
type CargaStruc struct {
	// Ejemplo: Destino (e.g., MDQ)
	Destino string `csv:"Destino"`
	// Ejemplo: Peso (e.g., 1500)
	Peso float64 `csv:"Peso"`
	// Ejemplo: Volumen (e.g., 2.5)
	Volumen float64 `csv:"Volumen"`
}

// ClienteStruc representa datos de "clientes.txt"
type ClienteStruc struct {
	// Ejemplo: Nombre (e.g., Juan)
	Nombre string `csv:"Nombre"`
	// Ejemplo: Apellido (e.g., Pérez)
	Apellido string `csv:"Apellido"`
	// Ejemplo: DNI (e.g., 12345678)
	DNI string `csv:"DNI"`
	// Ejemplo: Categoria (e.g., Platino)
	Categoria string `csv:"Categoria"`
}

// ConfiguracionAsientoStruc representa datos de "configuracion_asientos.txt"
type ConfiguracionAsientoStruc struct {
	// Ejemplo: CódigoAeronave (e.g., LV085)
	CodigoAeronave string `csv:"CódigoAeronave"`
	// Ejemplo: Zona (e.g., 1)
	Zona int `csv:"Zona"`
	// Ejemplo: AsientoInicial (e.g., 1)
	AsientoInicial int `csv:"AsientoInicial"`
	// Ejemplo: AsientoFinal (e.g., 20)
	AsientoFinal int `csv:"AsientoFinal"`
}

// EdificioStruc representa datos de "edificios.txt"
type EdificioStruc struct {
	// Ejemplo: xi (e.g., 100)
	Xi float64 `csv:"xi"`
	// Ejemplo: altura (e.g., 200)
	Altura float64 `csv:"altura"`
	// Ejemplo: xf (e.g., 250)
	Xf float64 `csv:"xf"`
}

// ReservaStruc representa datos de "reservas.txt"
type ReservaStruc struct {
	// Ejemplo: CodReserva (e.g., 1)
	CodReserva int `csv:"CodReserva"`
	// Ejemplo: DNIPasajero (e.g., 12345678)
	DNIPasajero string `csv:"DNIPasajero"`
	// Ejemplo: NroVuelo (e.g., 1001)
	NroVuelo string `csv:"NroVuelo"`
	// Ejemplo: FechaReserva (e.g., 2023-06-15)
	FechaReserva string `csv:"FechaReserva"`
	// Ejemplo: EstadoReserva (e.g., Confirmada)
	EstadoReserva string `csv:"EstadoReserva"`
}

// VueloStruc representa datos de "vuelos.txt"
type VueloStruc struct {
	// Ejemplo: numero (e.g., AA123)
	Numero string `csv:"numero"`
	// Ejemplo: fecha_hora (e.g., 2023-10-01 08:00:00)
	FechaHora time.Time `csv:"fecha_hora"`
	// Ejemplo: destino (e.g., MDQ)
	Destino string `csv:"destino"`
	// Ejemplo: aeronave (e.g., LV085)
	Aeronave string `csv:"aeronave"`
}
