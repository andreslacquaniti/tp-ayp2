package checkin

import (
	"fmt"
	"math/rand"
	"sigoa/internal/app"
	"sigoa/internal/guardar"
	"sigoa/internal/models"
	"sigoa/internal/utils"
	"sigoa/internal/vuelo"
	"sync"
	"time"
)

// 2. Buscar pasajero por DNI
func BuscarPasajero(dni string) *models.ClienteStruc {
	for _, c := range Clientes {
		if c.DNI == dni {
			return &c
		}
	}
	return nil
}

func CalculaMostradores(n int) int {
	switch {
	case n >= 1 && n <= 10:
		return 1
	case n >= 11 && n <= 20:
		return 2
	case n >= 21 && n <= 30:
		return 3
	case n > 30:
		return 5
	default:
		return 0
	}
}

// StartMostrador crea tantos goroutines "Mostrador" como el primer dígito
// de la longitud de la cola, procesa esos pasajeros y espera a que terminen.
func StartMostrador() {
	for QueueLlegada.Len() > 0 {
		nMostradores := CalculaMostradores(QueueLlegada.Len())
		var wg sync.WaitGroup
		for i := range nMostradores {
			elem := QueueLlegada.Front()
			if elem == nil {
				break
			}
			QueueLlegada.Remove(elem)

			wg.Add(1)
			go func(idx int, lp *models.LlegadaPasajero) {
				defer wg.Done()
				// Llamamos a Mostrador SIN ningún wg.Done adicional
				Mostrador(idx, lp)
			}(i+1, elem.Value.(*models.LlegadaPasajero))
		}

		wg.Wait()
	}
}

func Mostrador(nRoMostrador int, cola *models.LlegadaPasajero) {
	// NO MÁS defer wg.Done() aquí
	cliente := BuscarPasajero(cola.DNI)
	if cliente == nil {
		utils.PrintLog(fmt.Sprintf("❌ Mostrador %d: Pasajero con DNI %s no encontrado",
			nRoMostrador, cola.DNI))
		return
	}

	tiempoEspera := int(app.HoraSistema.Sub(cola.Llegada).Minutes())
	utils.PrintLog(fmt.Sprintf("🛂 Mostrador %d: Pasajero %s ingresó al mostrador. Tiempo de espera: %d minutos",
		nRoMostrador, cliente.DNI, tiempoEspera))

	ProcesarCheckin(cliente)
}

func ProcesarCheckin(pasajero *models.ClienteStruc) {
	timeStart := app.HoraSistema
	// genero una pausa entre 1 y 15 min, por cada check-in
	time.Sleep(time.Duration(rand.Intn(401)+100) * time.Millisecond)
	for _, r := range Reservas {
		if r.DNIPasajero == pasajero.DNI {
			v := vuelo.GetVuelo(r.NroVuelo)
			if v.GetEstado() == "PreDespegue" || v.GetEstado() == "Despegue" {
				utils.PrintLog(fmt.Sprintf("🚨 Pasajero %s Llego Tarde, No puede abordar", pasajero.DNI))
				guardar.RegistroFinalStruc[r.NroVuelo].PasajerosNoPresentes = append(guardar.RegistroFinalStruc[r.NroVuelo].PasajerosNoPresentes, *pasajero)
			} else {
				prioridad := categoriaPrioridad[pasajero.Categoria]
				zona := ZonaSalida[pasajero.Categoria]
				lp := &models.LlegadaPasajero{
					Ticket:    utils.GeneraNroTicket(),
					DNI:       pasajero.DNI,
					Llegada:   app.HoraSistema,
					Prioridad: prioridad,
					Zonas:     zona,
				}

				IniciarCola(r.NroVuelo, lp)
				//	fmt.Println(r.NroVuelo)
				// Convert the duration to minutes (float) and then cast to int
				tiempoTraccuridoEnMinutos := int(app.HoraSistema.Sub(timeStart).Minutes())
				utils.PrintLog(fmt.Sprintf("✅ Checkin: Pasajero %s, hizo check-in exitoso. Tiempo del tramite: %d minutos", pasajero.DNI, tiempoTraccuridoEnMinutos))
				guardar.RegistroFinalStruc[r.NroVuelo].PasajerosEmbarcados = append(guardar.RegistroFinalStruc[r.NroVuelo].PasajerosEmbarcados, *lp)
			}
			break
		}
	}
}
