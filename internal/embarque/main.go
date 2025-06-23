// Archivo: internal/embarque/embarque.go
package embarque

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sigoa/internal/models"
	"sigoa/pkg/csvfile"
	"sigoa/pkg/huffman"
	"sort"
	"time"
)

type EmbarqueApp struct {
	Configuracion []models.ConfiguracionAsientoStruc
	Reservas      []models.ReservaStruc
	Clientes      []models.ClienteStruc
	vuelo         models.VueloStruc
}

var singleInstance *EmbarqueApp

func GetInstance(vuelo models.VueloStruc) *EmbarqueApp {
	if singleInstance == nil {
		singleInstance = &EmbarqueApp{
			vuelo: vuelo,
		}
	}
	return singleInstance
}

func (app *EmbarqueApp) Inicializar() {
	// Cargar configuraciones de asientos
	config, err := csvfile.CargaCSV[models.ConfiguracionAsientoStruc]("configuracion_asientos.txt")
	if err != nil {
		log.Fatalf("Error cargando configuracion de asientos: %v", err)
	}
	app.Configuracion = config

	// Cargar reservas
	reservas, err := csvfile.CargaCSV[models.ReservaStruc]("reservas.txt")
	if err != nil {
		log.Fatalf("Error cargando reservas: %v", err)
	}
	app.Reservas = reservas

	// Cargar clientes
	clientes, err := csvfile.CargaCSV[models.ClienteStruc]("clientes.txt")
	if err != nil {
		log.Fatalf("Error cargando clientes: %v", err)
	}
	app.Clientes = clientes
}

// Estructura para representar el embarque
type PasajeroEmbarque struct {
	DNI       string
	Zona      int
	Nombre    string
	Apellido  string
	Categoria string
}

func (app *EmbarqueApp) EjecutarEmbarque(codigoAeronave string) {
	var pasajeros []PasajeroEmbarque

	for _, r := range app.Reservas {
		if r.EstadoReserva == "Checkeado" {
			cliente := buscarCliente(app.Clientes, r.DNIPasajero)
			zona := obtenerZona(app.Configuracion, codigoAeronave, cliente.DNI)
			pasajeros = append(pasajeros, PasajeroEmbarque{
				DNI:       cliente.DNI,
				Nombre:    cliente.Nombre,
				Apellido:  cliente.Apellido,
				Categoria: cliente.Categoria,
				Zona:      zona,
			})
		}
	}

	// Ordenar por categoría y luego por zona (Platino > Oro > Plata > Normal)
	prioridades := map[string]int{"Platino": 1, "Oro": 2, "Plata": 3, "Normal": 4}
	sort.Slice(pasajeros, func(i, j int) bool {
		if prioridades[pasajeros[i].Categoria] != prioridades[pasajeros[j].Categoria] {
			return prioridades[pasajeros[i].Categoria] < prioridades[pasajeros[j].Categoria]
		}
		return pasajeros[i].Zona < pasajeros[j].Zona
	})

	// Construir la salida de embarque
	output := "Embarque de Pasajeros - " + time.Now().Format("2006-01-02 15:04:05") + "\n\n"
	for _, p := range pasajeros {
		output += fmt.Sprintf("✔ %s %s - DNI: %s - Zona: %d - Categoria: %s\n", p.Nombre, p.Apellido, p.DNI, p.Zona, p.Categoria)
	}

	// Codificar y guardar en archivo usando Huffman
	encoded := huffman.HuffmanEncode(output)
	filename := filepath.Join("output", time.Now().Format("20060102_150405")+".out")
	os.WriteFile(filename, encoded, 0644)

	fmt.Println("✔ Embarque realizado. Archivo generado:", filename)
}

func buscarCliente(lista []models.ClienteStruc, dni string) models.ClienteStruc {
	for _, c := range lista {
		if c.DNI == dni {
			return c
		}
	}
	return models.ClienteStruc{}
}

func obtenerZona(configs []models.ConfiguracionAsientoStruc, codigo, dni string) int {
	for _, c := range configs {
		if c.CodigoAeronave == codigo {
			return c.Zona
		}
	}
	return 0 // zona por defecto
}
