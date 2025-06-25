package vuelo

import (
	"log"
	"sort"
	"sync"
	"time"

	"sigoa/internal/app"
	"sigoa/internal/models"
	"sigoa/pkg/csvfile"
)

// Despachoc es una estructura Singleton que maneja la lógica relacionada a vuelos.
type VueloApp struct {
	Vuelo models.VueloStruc

	Estado string
}

var Edificios []models.EdificioStruc

// Mapa de prioridades por categoría
var EstadoVuevlo = map[int]string{
	1: "Cerrado",
	2: "CheckIn",
	3: "Embarque",
	4: "PreDespegue",
	5: "Despegue",
}

// Mutex para evitar condiciones de carrera en concurrencia
var (
	singleInstance = make(map[string]*VueloApp)
	// tengo en memoria todo el listado de vuelos
	vuelos []models.VueloStruc
	mu     sync.Mutex
)

// GetInstance devuelve la única instancia de VueloApp por número de vuelo
func GetVuelo(nroVuelo string) *VueloApp {
	mu.Lock()
	defer mu.Unlock()
	if singleInstance[nroVuelo] == nil {
		singleInstance[nroVuelo] = &VueloApp{}
		singleInstance[nroVuelo].inicializar(nroVuelo)
	}
	return singleInstance[nroVuelo]
}

// ActualizarEstado actualiza el estado del vuelo según la HoraSistema
func (c *VueloApp) ActualizarEstado() {
	ahora := app.HoraSistema

	switch {
	case ahora.Before(c.Vuelo.FechaHora.Add(-2 * time.Hour)):
		c.Estado = EstadoVuevlo[1] // Cerrado
	case ahora.After(c.Vuelo.FechaHora.Add(-2*time.Hour)) && ahora.Before(c.Vuelo.FechaHora.Add(-1*time.Hour)):
		c.Estado = EstadoVuevlo[2] // CheckIn
	case ahora.After(c.Vuelo.FechaHora.Add(-1*time.Hour)) && ahora.Before(c.Vuelo.FechaHora.Add(-15*time.Minute)):
		c.Estado = EstadoVuevlo[3] // Embarque
	case ahora.After(c.Vuelo.FechaHora.Add(-15*time.Minute)) && ahora.Before(c.Vuelo.FechaHora):
		c.Estado = EstadoVuevlo[4] // PreDespegue
	case ahora.Equal(c.Vuelo.FechaHora) || ahora.After(c.Vuelo.FechaHora):
		c.Estado = EstadoVuevlo[5] // Despegue
	}
}

// GetEstado retorna el estado actual del vuelo
func (c *VueloApp) GetEstado() string {
	return c.Estado
}

// GetVuelos devuelve el listado de vuelos ordenado por FechaHora de salida.
// Si ya se ha cargado previamente, retorna el cache sin volver a leer el CSV.
func GetVuelos() []models.VueloStruc {
	if vuelos != nil {
		return vuelos
	}

	vs, err := csvfile.CargaCSV[models.VueloStruc]("vuelos.txt")
	if err != nil {
		log.Fatalf("Error cargando vuelos: %v", err)
	}

	// Ordenar por FechaHora ascendente
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].FechaHora.Before(vs[j].FechaHora)
	})

	vuelos = vs
	return vuelos
}

// Inicializar carga los vuelos y edificios desde archivos CSV.
func (c *VueloApp) inicializar(nroVuelo string) {
	// Cargar edificios
	edificios, err := csvfile.CargaCSV[models.EdificioStruc]("edificios.txt")
	if err != nil {
		log.Fatalf("Error cargando edificios: %v", err)
	}
	Edificios = edificios

	// Buscar el vuelo exacto
	for _, v := range vuelos {
		if v.Numero == nroVuelo {
			c.Vuelo = v
			return
		}
	}

	log.Fatalf("Vuelo con número %s no encontrado en vuelos.txt", nroVuelo)
}
