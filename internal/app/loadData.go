package main

import (
	"log"
	"sigoa/internal/models"
	"sigoa/pkg/csvfile"
	"sigoa/internal/app" // Importar el paquete app
)

func main() {
	// Aeronaves
	aeronaves, err := csvfile.CargaCSV[models.AeronaveStruc]("aeronaves.txt")
	if err != nil {
		log.Fatalf("Error cargando aeronaves: %v", err)
	}
	app.Aeronaves = aeronaves // Asignar los datos cargados a la variable global en app

	// Aeropuertos
	aeropuertos, err := csvfile.CargaCSV[models.AeropuertoStruc]("aeropuertos.txt")
	if err != nil {
		log.Fatalf("Error cargando aeropuertos: %v", err)
	}
	app.Aeropuertos = aeropuertos // Asignar los datos cargados a la variable global en app

	// Cargas
	cargas, err := csvfile.CargaCSV[models.CargaStruc]("cargas.txt")
	if err != nil {
		log.Fatalf("Error cargando cargas: %v", err)
	}
	app.Cargas = cargas // Asignar los datos cargados a la variable global en app

	// Clientes
	clientes, err := csvfile.CargaCSV[models.ClienteStruc]("clientes.txt")
	if err != nil {
		log.Fatalf("Error cargando clientes: %v", err)
	}
	app.Clientes = clientes // Asignar los datos cargados a la variable global en app

	// Configuración de Asientos
	configAsientos, err := csvfile.CargaCSV[models.ConfiguracionAsientoStruc]("configuracion_asientos.txt")
	if err != nil {
		log.Fatalf("Error cargando configuración de asientos: %v", err)
	}
	app.ConfiguracionAsientos = configAsientos // Asignar los datos cargados a la variable global en app

	// Edificios
	edificios, err := csvfile.CargaCSV[models.EdificioStruc]("edificios.txt")
	if err != nil {
		log.Fatalf("Error cargando edificios: %v", err)
	}
	app.Edificios = edificios // Asignar los datos cargados a la variable global en app

	// Vuelos
	vuelos, err := csvfile.CargaCSV[models.VueloStruc]("vuelos.txt")
	if err != nil {
		log.Fatalf("Error cargando vuelos: %v", err)
	}
	fmt.Printf("Vuelos cargados: %+v\n\n", vuelos)
	app.Vuelos = vuelos // Asignar los datos cargados a la variable global en app

	// Reservas - Depende de Vuelos y Clientes, cargar después
	reservas, err := csvfile.CargaCSV[models.ReservaStruc]("reservas.txt")
	if err != nil {
		log.Fatalf("Error cargando reservas: %v", err)
	}
	app.Reservas = reservas // Asignar los datos cargados a la variable global en app

	// Aquí podrías iniciar la simulación principal una vez que los datos estén cargados
	// app.RunSimulation()
}

