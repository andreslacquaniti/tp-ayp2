package main

import (
	"fmt"
	"os"
	"sigoa/internal/carga"
	"sigoa/internal/checkin"
	"sigoa/internal/embarque"
	"sigoa/internal/models"
	"sigoa/pkg/csvfile"
	"sigoa/pkg/huffman"
	"time"
)

func main() {
	fmt.Println("🚀 Iniciando simulación de Check-in\n")

	// Paso 1: Inicializar app de check-in
	app := checkin.NewCkeckin()

	// Paso 2: Cargar vuelos (desde archivo)
	vuelos, err := csvfile.CargaCSV[models.VueloStruc]("vuelos.txt")
	if err != nil {
		fmt.Printf("❌ Error al cargar vuelos: %v\n", err)
		return
	}
	// Paso 3: Buscar vuelo específico
	numeroVuelo := "AA790" // Cambiar si hace falta
	var vueloSeleccionado *models.VueloStruc
	for _, v := range vuelos {
		if v.Numero == numeroVuelo {
			vueloSeleccionado = &v
			break
		}
	}
	if vueloSeleccionado == nil {
		fmt.Printf("⚠️  Vuelo %s no encontrado.\n", numeroVuelo)
		return
	}
	fmt.Printf("✈️  Vuelo seleccionado: %s - Destino: %s - Fecha: %s\n\n",
		vueloSeleccionado.Numero, vueloSeleccionado.Destino,
		vueloSeleccionado.FechaHora.Format("2006-01-02 15:04"))

	// Paso 4: Obtener pasajeros confirmados
	pasajeros := app.ObtenerPasajerosPorVuelo(*vueloSeleccionado)

	// Paso 5: Simular llegadas
	llegadas := app.SimularLlegadas(pasajeros)

	// Paso 6: Procesar check-in (capacidad simulada, ej: 5)
	capacidad := 5
	checkeados, listaEspera := app.ProcesarCheckin(llegadas, app.Reservas, capacidad)

	// Paso 7: Mostrar resumen
	fmt.Printf("\n🧾 Resumen del Check-in:\n")
	fmt.Printf("✔️  Pasajeros checkeados: %d\n", len(checkeados))
	fmt.Printf("⛔ En lista de espera: %d\n", len(listaEspera))
	app.MostrarListaEspera(listaEspera)

	fmt.Println("\n✅ Simulación finalizada.")

	// Paso 8: Asignar carga a la aeronave
	// Inicializa el módulo de carga
	appCarga := carga.GetInstance()
	appCarga.Inicializar()

	// Asigna carga a la aeronave "LV085"
	resultado := appCarga.AsignarCarga("LV085")

	// Mostrar carga seleccionada de forma legible
	fmt.Println("\n📋 Carga asignada a la aeronave LV085:")
	for _, c := range resultado {
		// Convertimos peso a toneladas si es necesario
		pesoStr := ""
		if c.Peso >= 1000 {
			pesoStr = fmt.Sprintf("%.2f toneladas", c.Peso/1000)
		} else {
			pesoStr = fmt.Sprintf("%.0f kg", c.Peso)
		}

		// Redondear volumen a 1 decimal
		volumenStr := fmt.Sprintf("%.1f m³", c.Volumen)

		fmt.Printf("✔ Destino: %-5s | Peso: %-14s | Volumen: %s\n", c.Destino, pesoStr, volumenStr)
	}

	// Paso 9: Ejecutar embarque

	// Obtener la instancia del módulo de embarque
	appEmbarque := embarque.GetInstance()
	// Inicializar el módulo (carga configuraciones, reservas y clientes)
	appEmbarque.Inicializar()
	// Ejecutar el embarque para una aeronave específica
	appEmbarque.EjecutarEmbarque("LV085") // Reemplazá por la matrícula que necesites

	// Leer archivo codificado
	data, err := os.ReadFile("output/20250622_194206.out") // Reemplazá por el archivo real
	if err != nil {
		panic(err)
	}

	// Decodificar contenido con Huffman
	decoded := huffman.HuffmanDecode(data)

	// Mostrar resultado
	fmt.Println("📋 Resultado decodificado:")
	fmt.Println(decoded)

	//

	fmt.Println("🚀 Iniciando simulación de despacho de vuelos...")

	// Obtener instancia de vuelo (Singleton)
	despachoApp := despacho.GetInstance()

	// Inicializa datos desde archivos
	despachoApp.Inicializar()

	// Calcula línea del horizonte y la guarda en archivo
	despachoApp.CalcularHorizonte("output/horizonte.out")

	// Obtener vuelos seguros listos para despachar
	vuelosListos := despachoApp.ObtenerVuelosSeguros()

	// Mostrar vuelos despachados
	fmt.Println("🛫 Vuelos despachados (ordenados por hora):")
	for _, v := range vuelosListos {
		fmt.Printf("✔ Vuelo %s a %s - %s - Aeronave: %s\n", v.Numero, v.Destino, v.FechaHora.Format("02/01/2006 15:04"), v.Aeronave)
	}

	fmt.Printf("✅ Total de vuelos despachados: %d\n", len(vuelosListos))
	time.Sleep(1 * time.Second)
}
