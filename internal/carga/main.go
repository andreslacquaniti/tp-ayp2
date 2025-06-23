package carga

import (
	"fmt"
	"log"
	"sigoa/internal/models"
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
	}
	return instancia
}

// Inicializa los datos desde archivos CSV
func (c *CargaApp) Inicializar() {
	fmt.Println("🚛 Inicializando módulo de Carga...")

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

	fmt.Println("✅ Datos cargados correctamente.")
}

// Asignar carga a una aeronave específica utilizando el algoritmo de la mochila
func (c *CargaApp) AsignarCarga(matricula string) []models.CargaStruc {
	var aeronave *models.AeronaveStruc
	for _, a := range c.Aeronaves {
		if a.Matricula == matricula {
			aeronave = &a
			break
		}
	}
	if aeronave == nil {
		fmt.Printf("⚠️ Aeronave %s no encontrada.\n", matricula)
		return nil
	}

	fmt.Printf("\n📦 Asignando carga a la aeronave %s (Capacidad: %.2f kg, Volumen: %.2f m3)\n",
		aeronave.Matricula, aeronave.CapacidadCarga, aeronave.VolumenCarga)

	// Algoritmo de la mochila extendido: peso y volumen como restricciones
	return mochilaCarga(c.Cargas, aeronave.CapacidadCarga, aeronave.VolumenCarga)
}

// Algoritmo tipo "Mochila" para seleccionar la mejor combinación de cargas
func mochilaCarga(items []models.CargaStruc, maxPeso, maxVolumen float64) []models.CargaStruc {
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

	fmt.Printf("✅ Carga asignada: %d items | Total: %.2f kg, %.2f m3\n", len(resultado), totalPeso, totalVol)
	return resultado
}
