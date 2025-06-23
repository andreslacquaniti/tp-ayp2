package main

import (
	"fmt"
	"sigoa/internal/app"
	"sigoa/internal/checkin"
	"sigoa/internal/models"
	"sigoa/internal/vuelo"
	"time"
)

func main() {

	app.Init()     // Inicializa la aplicación y carga los datos necesarios
	checkin.Init() // Carga las reservas y clientes para el sistema de check-in

	go IniciarControlDeVuelos()
	// Obtengo el listado de vuelos
	vuelos := vuelo.GetVuelos()
	if len(vuelos) == 0 {
		fmt.Println("❌ No se encontraron vuelos.")
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
					fmt.Printf("✅ Check-in abierto para el vuelo %s\n", vueloDetalle.Vuelo.Numero)
					go simulaCheckIn(vueloDetalle.Vuelo)
					break
				}
				time.Sleep(100 * time.Millisecond)
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
					fmt.Printf("✈️ Embarque: Iniciando  %s\n", vueloDetalle.Vuelo.Numero)
					go genearEmbarque(vueloDetalle.Vuelo)
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

	}

	// Espera a que todas las goroutines terminen
	app.Wg.Wait()
}

func IniciarControlDeVuelos() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		for _, v := range vuelo.GetVuelos() {
			instance := vuelo.GetVuelo(v.Numero)
			instance.ActualizarEstado()
			//fmt.Printf("Estado del vuelo %s: %s\n", v.Numero, instance.GetEstado())
		}
	}
}

func genearEmbarque(vuelo models.VueloStruc) {
	app.Wg.Add(1)
	defer app.Wg.Done()
	// genearEmbarque := embarque.GetInstance(vuelo)
	// genearEmbarque.EjecutarEmbarque()
}

func simulaCheckIn(vuelo models.VueloStruc) {
	app.Wg.Add(1)
	defer app.Wg.Done()
	pasajeros := checkin.ObtenerPasajerosPorVuelo(vuelo)
	checkin.SimularLlegadas(pasajeros)
	checkin.Mostrador()
}
