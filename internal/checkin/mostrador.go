package checkin

import (
	"fmt"
	"math/rand"
	"sigoa/internal/models"
	"sigoa/internal/vuelo"
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

func Mostrador() {

	for e := QueueLlegada.Front(); e != nil; e = e.Next() {
		// 1) Tomamos el valor como string
		dni, _ := e.Value.(string)
		// 2) Buscamos al cliente por DNI
		cliente := BuscarPasajero(dni)
		if cliente != nil {
			fmt.Printf("🛂 Mostrador: %s, %s, %s Procesando pasajero\n", cliente.Apellido, cliente.Nombre, cliente.DNI)
			ProcesarCheckin(cliente)
		} else {
			fmt.Printf("❌ Mostrador: Pasajero con DNI %s no encontrado\n", dni)
		}
	}
}

func ProcesarCheckin(pasajero *models.ClienteStruc) {
	// genero una pausa entre 1 y 15 min, por cada check-in
	time.Sleep(time.Duration(rand.Intn(401)+100) * time.Millisecond)
	for _, r := range Reservas {
		if r.DNIPasajero == pasajero.DNI {
			v := vuelo.GetVuelo(r.NroVuelo)
			if v.GetEstado() == "PreDespegue" {
				fmt.Printf("🚨 Pasajero DNI %s Llego Tarde, No puede abordar\n", pasajero.DNI)
			} else {
				prioridad := categoriaPrioridad[pasajero.Categoria]
				zona := ZonaSalida[pasajero.Categoria]
				lp := &LlegadaPasajero{
					DNI:       pasajero.DNI,
					Llegada:   time.Now(),
					Prioridad: prioridad,
					Zonas:     zona,
				}
				IniciarCola(r.NroVuelo, lp)
				fmt.Printf("✔️ Checkin: Pasajero %s, %s, %s hizo check-in exitoso. (Posición %d)\n", pasajero.Apellido, pasajero.Nombre, pasajero.DNI, lp.Index)
			}
			break
		}
	}
}
