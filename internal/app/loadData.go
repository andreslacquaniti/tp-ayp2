package main

import (
	"fmt"
	"log"
	"sigoa/internal/models"
	"sigoa/pkg/csvfile"
)

func main() {
	// Aeronaves
	aeronaves, err := csvfile.CargaCSV[models.AeronaveStruc]("aeronaves.txt")
	if err != nil {
		log.Fatalf("Error cargando aeronaves: %v", err)
	}
	fmt.Printf("Aeronaves cargadas: %+v\n\n", aeronaves)

	// Aeropuertos
	aeropuertos, err := csvfile.CargaCSV[models.AeropuertoStruc]("aeropuertos.txt")
	if err != nil {
		log.Fatalf("Error cargando aeropuertos: %v", err)
	}
	fmt.Printf("Aeropuertos cargados: %+v\n\n", aeropuertos)

	// Cargas
	cargas, err := csvfile.CargaCSV[models.CargaStruc]("cargas.txt")
	if err != nil {
		log.Fatalf("Error cargando cargas: %v", err)
	}
	fmt.Printf("Cargas cargadas: %+v\n\n", cargas)

	// Clientes
	clientes, err := csvfile.CargaCSV[models.ClienteStruc]("clientes.txt")
	if err != nil {
		log.Fatalf("Error cargando clientes: %v", err)
	}
	fmt.Printf("Clientes cargados: %+v\n\n", clientes)

	// Configuración de Asientos
	configAsientos, err := csvfile.CargaCSV[models.ConfiguracionAsientoStruc]("configuracion_asientos.txt")
	if err != nil {
		log.Fatalf("Error cargando configuración de asientos: %v", err)
	}
	fmt.Printf("Configuración de asientos cargada: %+v\n\n", configAsientos)

	// Edificios
	edificios, err := csvfile.CargaCSV[models.EdificioStruc]("edificios.txt")
	if err != nil {
		log.Fatalf("Error cargando edificios: %v", err)
	}
	fmt.Printf("Edificios cargados: %+v\n\n", edificios)

	// Reservas
	reservas, err := csvfile.CargaCSV[models.ReservaStruc]("reservas.txt")
	if err != nil {
		log.Fatalf("Error cargando reservas: %v", err)
	}
	fmt.Printf("Reservas cargadas: %+v\n\n", reservas)

	// Vuelos
	vuelos, err := csvfile.CargaCSV[models.VueloStruc]("vuelos.txt")
	if err != nil {
		log.Fatalf("Error cargando vuelos: %v", err)
	}
	fmt.Printf("Vuelos cargados: %+v\n\n", vuelos)
}
