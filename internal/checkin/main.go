package checkin

import (
	"container/heap"
	"fmt"
)

// ---------- Tipos y estructuras ----------

var Pqueue = make(map[string]*PrioridadQueue)

// ---------- Inicialización ----------
func IniciarCola(nroVuelo string, lp *LlegadaPasajero) {
	if Pqueue[nroVuelo] == nil {
		Pqueue[nroVuelo] = &PrioridadQueue{}
		heap.Init(Pqueue[nroVuelo])
	}
	heap.Push(Pqueue[nroVuelo], lp)

}

// ---------- Funciones ----------

// 5. Mostrar lista de espera
func MostrarListaEspera(lista []string) {
	fmt.Println("\n📋 Lista de espera:")
	if len(lista) == 0 {
		fmt.Println("✅ No hay pasajeros en lista de espera.")
	} else {
		for _, dni := range lista {
			fmt.Println(" - DNI:", dni)
		}
	}
}
