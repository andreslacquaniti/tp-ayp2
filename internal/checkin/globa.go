package checkin

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"sigoa/internal/app"
	"sigoa/internal/models"
	"sigoa/internal/utils"
	"sigoa/pkg/csvfile"
	"time"
)

var QueueLlegada *list.List // cola de llegadas FIFO
var Reservas []models.ReservaStruc
var Clientes []models.ClienteStruc

// ---------- Inicialización ----------

func Init() {
	utils.PrintLog(fmt.Sprint("🛫 Inicializando sistema de Check-in..."))
	QueueLlegada = list.New() // inicializa la cola de llegadas
	// Cargar reservas
	reservas, err := csvfile.CargaCSV[models.ReservaStruc]("reservas.txt")
	if err != nil {
		log.Fatalf("❌ Error cargando reservas: %v", err)
	}
	Reservas = reservas
	utils.PrintLog(fmt.Sprintf("✔️  %d reservas cargadas", len(reservas)))

	// Cargar clientes
	clientes, err := csvfile.CargaCSV[models.ClienteStruc]("clientes.txt")
	if err != nil {
		log.Fatalf("❌ Error cargando clientes: %v", err)
	}
	Clientes = clientes
	utils.PrintLog(fmt.Sprintf("✔️  %d clientes cargados", len(clientes)))
}

// 1. Obtener pasajeros por vuelo
func ObtenerPasajerosPorVuelo(vuelo models.VueloStruc) []models.ClienteStruc {
	utils.PrintLog(fmt.Sprintf("🔍 Buscando pasajeros del vuelo %s...", vuelo.Numero))
	var pasajeros []models.ClienteStruc
	for _, reserva := range Reservas {
		if reserva.NroVuelo == vuelo.Numero && reserva.EstadoReserva == "Confirmada" {
			for _, cliente := range Clientes {
				if cliente.DNI == reserva.DNIPasajero {
					pasajeros = append(pasajeros, cliente)
					break
				}
			}
		}
	}
	// Mezclar el orden de los pasajeros
	rand.Shuffle(len(pasajeros), func(i, j int) {
		pasajeros[i], pasajeros[j] = pasajeros[j], pasajeros[i]
	})

	utils.PrintLog(fmt.Sprintf("🧍 %d pasajeros encontrados para el vuelo %s", len(pasajeros), vuelo.Numero))
	return pasajeros
}

func SimularLlegadas(pasajeros []models.ClienteStruc) {
	for _, p := range pasajeros {
		// genero una pausa entre 1 y 5 min, por cada llegada
		time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
		llegada := app.HoraSistema
		utils.PrintLog(fmt.Sprintf("📍 Llegada: %s %s (DNI: %s) , Hora %s", p.Nombre, p.Apellido, p.DNI, llegada.Format("15:04")))

		prioridad := categoriaPrioridad[p.Categoria]
		zona := ZonaSalida[p.Categoria]
		cola := &models.LlegadaPasajero{
			DNI:       p.DNI,
			Llegada:   llegada,
			Prioridad: prioridad,
			Zonas:     zona,
		}
		QueueLlegada.PushBack(cola)
	}
}
