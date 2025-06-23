package app

import (
	"fmt"
	"log"
	"sigoa/internal/models"
	"sync"
	"time"
)

var HoraSistema time.Time
var Wg sync.WaitGroup

var RegistroFinal *RegistroVueloStruc

// RegistroVuelo agrupa toda la información de un vuelo tras el embarque.
type RegistroVueloStruc struct {
	NumeroVuelo          string                // número de vuelo
	SalidaProgramada     time.Time             // fecha y hora de salida programada
	PartidaReal          time.Time             // fecha y hora de partida real
	PasajerosEmbarcados  []models.ClienteStruc // lista de pasajeros que embarcaron
	PasajerosNoPresentes []models.ClienteStruc // lista de pasajeros que no se presentaron
	ListaEspera          []models.ClienteStruc // lista de espera (y si finalmente embarcaron)
	//EquipajesDespachados []Equipaje            // bultos despachados como equipaje
	// PaquetesCarga        []PaqueteCarga        // paquetes de carga embarcados
}

func init() {
	var err error
	HoraSistema, err = time.Parse("2006-01-02 15:04:05", "2025-06-23 07:30:00")
	if err != nil {
		log.Fatalf("Error al parsear la hora inicial: %v", err)
	}
}

func Init() {
	Wg.Add(1)
	go func() {
		defer Wg.Done()
		const incremento = 10 * time.Minute // 10 minutos por cada segundo real
		const duracionMax = 12 * time.Hour  // Duración máxima de simulación
		tiempoFinal := HoraSistema.Add(duracionMax)

		fmt.Println("⏱️ HoraSistema inicial:", HoraSistema.Format("2006-01-02 15:04:05"))
		for HoraSistema.Before(tiempoFinal) {
			time.Sleep(1 * time.Second) // Esperar un segundo real
			HoraSistema = HoraSistema.Add(incremento)
			fmt.Println("⏰ HoraSistema:", HoraSistema.Format("2006-01-02 15:04:05"))
		}

		fmt.Println("✅ Simulación completa. HoraSistema final:", HoraSistema.Format("2006-01-02 15:04:05"))
	}()
}
