// Archivo: internal/embarque/embarque.go
package embarque

import (
	"fmt"
	"log"
	"sigoa/internal/carga"
	"sigoa/internal/checkin"
	"sigoa/internal/guardar"
	"sigoa/internal/models"
	"sigoa/internal/utils"
	"sigoa/internal/vuelo"
	"sigoa/pkg/csvfile"
)

type EmbarqueApp struct {
	Configuracion []models.ConfiguracionAsientoStruc
	vuelo         models.VueloStruc
}

// var singleInstance *EmbarqueApp

func NewEmbarque(vuelo models.VueloStruc) *EmbarqueApp {
	singleInstance := &EmbarqueApp{
		vuelo: vuelo,
	}
	singleInstance.inicializar()
	return singleInstance
}

func (c *EmbarqueApp) inicializar() {
	// Cargar configuraciones de asientos
	config, err := csvfile.CargaCSV[models.ConfiguracionAsientoStruc]("configuracion_asientos.txt")
	if err != nil {
		log.Fatalf("Error cargando configuracion de asientos: %v", err)
	}
	c.Configuracion = config
}

func (c *EmbarqueApp) ProcesarEmbarque() {
	// Apertura de Embarque.
	for checkin.Pqueue[c.vuelo.Numero] == nil {
	}

	// Subierndo Pasajeros en la zona asignada xxx
	pasajeros := checkin.Pqueue[c.vuelo.Numero]
	for _, pasajero := range *pasajeros {
		utils.PrintLog(fmt.Sprintf("✔ Vuelo %s , Embarque de Pasajero %s, en la Zona %d:", c.vuelo.Numero, pasajero.DNI, pasajero.Zonas))
	}

	// Cargando Carga
	carga.GetInstance().ProcesarCarga(c.vuelo)

	// Calculando Altura para el vuelo
	for {
		esSeguro, altura := vuelo.GetVuelo(c.vuelo.Numero).VuelosSeguro()
		if esSeguro {
			utils.PrintLog(fmt.Sprintf("✈️ Salida: Vuelo Seguro %s, Altura: %.0f metros", c.vuelo.Numero, float64(altura)))
			break
		} else {
			utils.PrintLog(fmt.Sprintf("⚠️ Salida: Vuelo INSEGURO %s, Altura : %.0f metros, esperando condiciones seguras...", c.vuelo.Numero, float64(altura)))
		}
	}

	for vuelo.GetVuelo(c.vuelo.Numero).GetEstado() != "Despegue" {
	}

	utils.PrintLog(fmt.Sprintf("✈️ Embarque: Vuelo %s CERRADO ", c.vuelo.Numero))

	for vuelo.GetVuelo(c.vuelo.Numero).GetEstado() != "Despegue" {
	}

	utils.PrintLog(fmt.Sprintf("✈️ Salida: Vuelo %s en Despegue...", c.vuelo.Numero))
	//delete(checkin.Pqueue, app.vuelo.Numero)

	guardar.GuardarRegistroVueloEnJson(c.vuelo.Numero)
}
