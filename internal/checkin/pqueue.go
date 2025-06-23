package checkin

import (
	"container/heap"
	"time"
)

// ---------- Cola de prioridad (Heap) ----------

type LlegadaPasajero struct {
	DNI       string
	Prioridad int
	Llegada   time.Time
	Zonas     int
	Index     int
}

// Mapa de prioridades por categoría
var categoriaPrioridad = map[string]int{
	"Platino": 1,
	"Oro":     2,
	"Plata":   3,
	"Normal":  4,
}

// Mapa de prioridades por categoría
var ZonaSalida = map[string]int{
	"Platino": 1,
	"Oro":     2,
	"Plata":   3,
	"Normal":  3,
}

type PrioridadQueue []*LlegadaPasajero

func (pq PrioridadQueue) Len() int { return len(pq) }

func (pq PrioridadQueue) Less(i, j int) bool {
	if pq[i].Prioridad == pq[j].Prioridad {
		return pq[i].Llegada.Before(pq[j].Llegada)
	}
	return pq[i].Prioridad < pq[j].Prioridad
}

func (pq PrioridadQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PrioridadQueue) Push(x any) {
	n := len(*pq)
	item := x.(*LlegadaPasajero)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PrioridadQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}
func (pq *PrioridadQueue) Actualizar(item *LlegadaPasajero, prioridad int, llegada time.Time) {
	item.Prioridad = prioridad
	item.Llegada = llegada
	heap.Fix(pq, item.Index)
}
