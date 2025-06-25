package app

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var HoraSistema time.Time
var Wg sync.WaitGroup

func init() {
	var err error

	HoraSistema, err = time.Parse("2006-01-02 15:04:05", "2025-06-23 05:30:00")
	if err != nil {
		log.Fatalf("Error al parsear la hora inicial: %v", err)
	}
}

func Init() {
	Wg.Add(1)
	go func() {
		defer Wg.Done()
		const incremento = 1 * time.Minute         // Incrementa 1 minuto en la HoraSistema
		const duracionReal = 50 * time.Millisecond // Cada 50 milisegundos de tiempo real
		const duracionMax = 28 * time.Hour         // Duración máxima de simulación

		tiempoFinalSimulacion := HoraSistema.Add(duracionMax)

		//fmt.Println("⏱️ HoraSistema inicial:", HoraSistema.Format("2006-01-02 15:04:05"))

		// Registrar el tiempo real de inicio de la simulación
		inicioTiempoReal := time.Now()

		for HoraSistema.Before(tiempoFinalSimulacion) {
			time.Sleep(duracionReal) // Esperar 50 milisegundos reales
			HoraSistema = HoraSistema.Add(incremento)
			// if HoraSistema.Minute()%10 == 0 { // Imprimir cada 10 minutos simulados
			// 	fmt.Sprint("⏰ HoraSistema:", HoraSistema.Format("2006-01-02 15:04:05"))
			// }
		}

		// Registrar el tiempo real de fin de la simulación
		finTiempoReal := time.Now()
		duracionSimulacionReal := finTiempoReal.Sub(inicioTiempoReal)

		fmt.Println("✅ Simulación completa. HoraSistema final:", HoraSistema.Format("2006-01-02 15:04:05"))
		fmt.Printf("⏳ La simulación duró un tiempo real de: %s \n", duracionSimulacionReal.String())
		os.Exit(0)

	}()
}
