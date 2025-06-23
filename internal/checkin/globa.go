package checkin

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"sigoa/internal/app"
	"sigoa/internal/models"
	"sigoa/pkg/csvfile"
	"time"
)

var RegistroVuelo *app.RegistroVueloStruc
var QueueLlegada *list.List // cola de llegadas FIFO
var Reservas []models.ReservaStruc
var Clientes []models.ClienteStruc

// ---------- Inicialización ----------

func Init() {
	fmt.Println("🛫 Inicializando sistema de Check-in...")
	QueueLlegada = list.New() // inicializa la cola de llegadas
	// Cargar reservas
	reservas, err := csvfile.CargaCSV[models.ReservaStruc]("reservas.txt")
	if err != nil {
		log.Fatalf("❌ Error cargando reservas: %v", err)
	}
	Reservas = reservas
	fmt.Printf("✔️  %d reservas cargadas\n", len(reservas))

	// Cargar clientes
	clientes, err := csvfile.CargaCSV[models.ClienteStruc]("clientes.txt")
	if err != nil {
		log.Fatalf("❌ Error cargando clientes: %v", err)
	}
	Clientes = clientes
	fmt.Printf("✔️  %d clientes cargados\n", len(clientes))
}

// 1. Obtener pasajeros por vuelo
func ObtenerPasajerosPorVuelo(vuelo models.VueloStruc) []models.ClienteStruc {
	fmt.Printf("🔍 Buscando pasajeros del vuelo %s...\n", vuelo.Numero)
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

	fmt.Printf("🧍 %d pasajeros encontrados para el vuelo %s\n", len(pasajeros), vuelo.Numero)
	return pasajeros
}

func SimularLlegadas(pasajeros []models.ClienteStruc) {
	for _, p := range pasajeros {
		// genero una pausa entre 1 y 15 min, por cada check-in
		//	time.Sleep(time.Duration(rand.Intn(401)+100) * time.Millisecond)
		time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
		llegada := app.HoraSistema
		fmt.Printf("📍 %s %s (DNI: %s) llegó a las %s\n", p.Nombre, p.Apellido, p.DNI, llegada.Format("15:04"))
		QueueLlegada.PushBack(p.DNI)
	}
}
