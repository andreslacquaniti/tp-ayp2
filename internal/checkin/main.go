package checkin

import (
	"container/heap"
	"fmt"
	"sigoa/internal/models"
	"sigoa/internal/utils"
)

// ---------- Tipos y estructuras ----------

var Pqueue = make(map[string]*PrioridadQueue)

// ---------- Inicialización ----------
func IniciarCola(nroVuelo string, lp *models.LlegadaPasajero) {
	if Pqueue[nroVuelo] == nil {
		Pqueue[nroVuelo] = &PrioridadQueue{}
		heap.Init(Pqueue[nroVuelo])

	}
	heap.Push(Pqueue[nroVuelo], lp)

}

// ---------- Funciones ----------

// 5. Mostrar lista de espera
func MostrarListaEspera(lista []string) {
	utils.PrintLog("📋 Lista de espera:")
	if len(lista) == 0 {
		utils.PrintLog("✅ No hay pasajeros en lista de espera.")
	} else {
		for _, dni := range lista {
			utils.PrintLog(fmt.Sprint(" - DNI:", dni))
		}
	}
}
