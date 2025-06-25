package carga

import (
	"fmt"
	"log"
	"sigoa/internal/models"
	"sigoa/internal/utils"
	"sigoa/pkg/csvfile"
)

// Estructura principal de la app de carga
type CargaApp struct {
	Aeronaves []models.AeronaveStruc
	Cargas    []models.CargaStruc
}

var instancia *CargaApp

// Singleton: Obtener la única instancia de CargaApp
func GetInstance() *CargaApp {
	if instancia == nil {
		instancia = &CargaApp{}
		instancia.inicializar()
	}
	return instancia
}

// Inicializa los datos desde archivos CSV
func (c *CargaApp) inicializar() {
	utils.PrintLog("🚛 Inicializando módulo de Carga...")

	// Cargar aeronaves
	aeronaves, err := csvfile.CargaCSV[models.AeronaveStruc]("aeronaves.txt")
	if err != nil {
		log.Fatalf("❌ Error cargando aeronaves: %v", err)
	}
	c.Aeronaves = aeronaves

	// Cargar cargas
	cargas, err := csvfile.CargaCSV[models.CargaStruc]("cargas.txt")
	if err != nil {
		log.Fatalf("❌ Error cargando cargas: %v", err)
	}
	c.Cargas = cargas

	utils.PrintLog("✅ Datos cargados correctamente.")
}

func (c *CargaApp) getAeronave(matricula string) models.AeronaveStruc {
	for _, a := range c.Aeronaves {
		if a.Matricula == matricula {
			return a
		}
	}
	utils.PrintLog(fmt.Sprintf("⚠️ Aeronave %s no encontrada.", matricula))
	return models.AeronaveStruc{}
}

func (c *CargaApp) getCargas(destino string) []models.CargaStruc {
	var cargasDestino []models.CargaStruc
	for _, carga := range c.Cargas {
		if carga.Destino == destino {
			cargasDestino = append(cargasDestino, carga)
		}
	}
	if len(cargasDestino) == 0 {
		utils.PrintLog(fmt.Sprintf("⚠️ Carga para destino %s no encontrada.", destino))
	}
	return cargasDestino
}

// Asignar carga a una aeronave específica utilizando el algoritmo de la mochila
func (c *CargaApp) ProcesarCarga(vuelo models.VueloStruc) {
	aeronave := c.getAeronave(vuelo.Matricula)
	cargas := c.getCargas(vuelo.Destino)

	// Algoritmo de la mochila extendido: peso y volumen como restricciones
	mochilaCarga(cargas, aeronave, vuelo.Numero)

}
func mochilaCarga(items []models.CargaStruc, aeronave models.AeronaveStruc, vuelo string) {

	var resultado []models.CargaStruc
	var fuera []models.CargaStruc
	pesoActual := 0.0
	volumenActual := 0.0

	for _, item := range items {
		if pesoActual+item.Peso <= aeronave.CapacidadCarga && volumenActual+item.Volumen <= aeronave.VolumenCarga {
			resultado = append(resultado, item)
			pesoActual += item.Peso
			volumenActual += item.Volumen
		} else {
			fuera = append(fuera, item)
		}
	}
	utils.PrintLog(fmt.Sprintf("✈️ Carga : Capacidad: %s %.2f kg | %.2f m3", vuelo, aeronave.CapacidadCarga, aeronave.VolumenCarga))
	utils.PrintLog(fmt.Sprintf("✅ Carga : Vuelo %s, Asignada: %d items | Total: %.2f kg, %.2f m3", vuelo, len(resultado), aeronave.CapacidadCarga, aeronave.VolumenCarga))
	utils.PrintLog(fmt.Sprintf("⚠️ Carga : Vuelo %s, Fuera de la aeronave: %d items", vuelo, len(fuera)))
}

// Algoritmo tipo "Mochila" para seleccionar la mejor combinación de cargas
func mochilaCarga2(items []models.CargaStruc, maxPeso, maxVolumen float64) {
	n := len(items)
	dp := make([][]float64, n+1)

	// Usamos un mapa auxiliar para guardar qué ítems se toman
	seleccionados := make([]bool, n)

	// Inicialización de la matriz dinámica
	for i := 0; i <= n; i++ {
		dp[i] = make([]float64, int(maxPeso)+1)
	}

	// Lógica: Maximizar volumen usado sin exceder peso ni volumen
	for i := 1; i <= n; i++ {
		for w := int(maxPeso); w >= int(items[i-1].Peso); w-- {
			if items[i-1].Volumen+dp[i-1][w-int(items[i-1].Peso)] > dp[i][w] &&
				items[i-1].Volumen <= maxVolumen {
				dp[i][w] = items[i-1].Volumen + dp[i-1][w-int(items[i-1].Peso)]
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	// Reconstrucción de la solución
	w := int(maxPeso)
	for i := n; i >= 1; i-- {
		if dp[i][w] != dp[i-1][w] {
			seleccionados[i-1] = true
			w -= int(items[i-1].Peso)
		}
	}

	// Extraer ítems seleccionados
	var resultado []models.CargaStruc
	totalPeso := 0.0
	totalVol := 0.0
	for i, sel := range seleccionados {
		if sel {
			resultado = append(resultado, items[i])
			totalPeso += items[i].Peso
			totalVol += items[i].Volumen
		}
	}

	utils.PrintLog(fmt.Sprintf("✅ Carga asignada: %d items | Total: %.2f kg, %.2f m3", len(resultado), totalPeso, totalVol))
}
