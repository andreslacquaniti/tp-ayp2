package main

import (
	"fmt"
	"sigoa/internal/app"
	"sigoa/internal/carga"
	"sigoa/internal/checkin"
	"sigoa/internal/embarque"
	"sigoa/internal/guardar"
	"sigoa/internal/models"
	"sigoa/internal/utils"
	"sigoa/internal/vuelo"
	"time"
)

func main() {

	app.Init()     // Inicializa la aplicación y carga los datos necesarios
	checkin.Init() // Carga las reservas y clientes para el sistema de check-in
	carga.GetInstance()
	guardar.Init() // Inicializa el registro de vuelos

	horizonte := fmt.Sprintf("horizonte_%d.huff", time.Now().UnixNano())
	vuelo.CalcularHorizonte(horizonte)

	go IniciarControlDeVuelos()
	// Obtengo el listado de vuelos
	vuelos := vuelo.GetVuelos()
	if len(vuelos) == 0 {
		utils.PrintLog("❌ No se encontraron vuelos.")
		return
	}

	for _, v := range vuelos {
		app.Wg.Add(1)
		go func() {
			defer app.Wg.Done()
			vueloDetalle := vuelo.GetVuelo(v.Numero)
			// Simula la espera hasta que el check-in esté abierto y no haya pendiente de salida
			for {
				if vueloDetalle.GetEstado() == "CheckIn" {
					utils.PrintLog(fmt.Sprintf("✅ Check-in abierto para el vuelo %s", vueloDetalle.Vuelo.Numero))
					go simulaCheckIn(vueloDetalle.Vuelo)
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()

		//ControlEmbaque
		app.Wg.Add(1)
		go func() {
			defer app.Wg.Done()
			vueloDetalle := vuelo.GetVuelo(v.Numero)
			// Simula la espera hasta que el check-in esté abierto y no haya pendiente de salida
			for {
				if vueloDetalle.GetEstado() == "Embarque" {
					utils.PrintLog(fmt.Sprintf("✈️ Embarque: Iniciando  %s", vueloDetalle.Vuelo.Numero))
					go genearEmbarque(vueloDetalle.Vuelo)
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
		}()

	}

	// Espera a que todas las goroutines terminen
	app.Wg.Wait()
}

func IniciarControlDeVuelos() {
	for {
		for _, v := range vuelo.GetVuelos() {
			instance := vuelo.GetVuelo(v.Numero)
			instance.ActualizarEstado()
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func genearEmbarque(vuelo models.VueloStruc) {
	app.Wg.Add(1)
	defer app.Wg.Done()
	genearEmbarque := embarque.NewEmbarque(vuelo)
	genearEmbarque.ProcesarEmbarque()
}

func simulaCheckIn(vuelo models.VueloStruc) {
	app.Wg.Add(1)
	defer app.Wg.Done()
	pasajeros := checkin.ObtenerPasajerosPorVuelo(vuelo)
	checkin.SimularLlegadas(pasajeros)
	checkin.StartMostrador()
}
