package guardar

import (
	"encoding/json"
	"fmt"
	"os"
	"sigoa/internal/app"
	"sigoa/internal/models"
	"sigoa/internal/utils"
	"sigoa/internal/vuelo"
	"sigoa/pkg/huffman"
	"time"
)

var RegistroFinalStruc map[string]*RegistroVueloStruc

// RegistroVuelo agrupa toda la información de un vuelo tras el embarque.
type RegistroVueloStruc struct {
	NumeroVuelo          string                   `json:"numeroVuelo"`
	SalidaProgramada     time.Time                `json:"salidaProgramada"`
	PartidaReal          time.Time                `json:"partidaReal"`
	PasajerosEmbarcados  []models.LlegadaPasajero `json:"pasajerosEmbarcados"`
	PasajerosNoPresentes []models.ClienteStruc    `json:"pasajerosNoPresentes"`
	ListaEspera          []models.ClienteStruc    `json:"listaEspera"`
}

// en el directorio "out".
func Init() {
	RegistroFinalStruc = make(map[string]*RegistroVueloStruc)

	vuelos := vuelo.GetVuelos()
	for _, v := range vuelos {
		RegistroFinalStruc[v.Numero] = &RegistroVueloStruc{
			NumeroVuelo:          v.Numero,
			SalidaProgramada:     v.FechaHora,
			PasajerosEmbarcados:  []models.LlegadaPasajero{},
			PasajerosNoPresentes: []models.ClienteStruc{},
			ListaEspera:          []models.ClienteStruc{},
		}
	}

}

// GuardarRegistroVueloEnJson toma una estructura RegistroVueloStruc y la guarda como un archivo JSON
// en el directorio "out".
func GuardarRegistroVueloEnJson(vuelo string) error {
	registro := RegistroFinalStruc[vuelo]
	tiempoAhora := time.Now().Format("2006-01-02_15-04-05")

	registro.PartidaReal = app.HoraSistema
	// Convierte la estructura a JSON con indentación para que sea legible
	datosJson, err := json.MarshalIndent(registro, "", "  ")
	if err != nil {
		return fmt.Errorf("error al convertir los datos a JSON: %w", err)
	}

	// Define la ruta del archivo
	rutaArchivo := fmt.Sprintf("output/registro_vuelo_%s_%s", registro.NumeroVuelo, tiempoAhora)

	// Escribe los datos JSON en el archivo
	err = os.WriteFile(rutaArchivo+".json", datosJson, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el JSON en el archivo: %w", err)
	}

	utils.PrintLog(fmt.Sprintf("Registro de vuelo para el vuelo %s guardado en %s", registro.NumeroVuelo, rutaArchivo))

	huffman.Guardar(string(datosJson), rutaArchivo+".huff")
	return nil
}
