package vuelo

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sigoa/internal/utils"
	"sort"
	"time"
)

// VuelosSeguro evalúa los vuelos según los edificios para evitar colisiones.
// Complejidad del algoritmo:
// - Verificación de seguridad: O(n × m), donde n = cantidad de vuelos, m = cantidad de edificios
// - Ordenamiento por fecha: O(n log n), donde n = cantidad de vuelos seguros
func (c *VueloApp) VuelosSeguro() (bool, int) {
	altitudVuelo := c.generarAltitudRandom() // Genera una altitud aleatoria para el vuelo
	seguro := true
	for _, edificio := range Edificios {
		// Verifica si algún edificio supera la altitud permitida
		if edificio.Altura >= altitudVuelo {
			utils.PrintLog(fmt.Sprintf("⚠️ Despague %s potencialmente en conflicto con edificio entre %.1f y %.1f (altura: %.1f)",
				c.Vuelo.Numero, edificio.Xi, edificio.Xf, edificio.Altura))
			seguro = false
			break // No es necesario seguir comprobando este vuelo
		}
	}
	return seguro, int(altitudVuelo)
}

// genera un número aleatorio tipo float64 entre 150 y 800.
func (c *VueloApp) generarAltitudRandom() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 400.0
	max := 800.0
	// Retorna un número aleatorio dentro del rango [min, max].
	// r.Float64() devuelve un valor en [0.0, 1.0), entonces se escala al rango deseado.
	return min + r.Float64()*(max-min)
}

// CalcularHorizonte genera la línea del horizonte a partir de los edificios y la guarda en un archivo.
// Complejidad: O(n log n) utilizando técnica sweep line y heap implícito con mapa.
func CalcularHorizonte(nombreArchivo string) {
	type Punto struct {
		x      float64
		esIni  bool
		height float64
	}

	var puntos []Punto
	for _, e := range Edificios {
		puntos = append(puntos, Punto{e.Xi, true, e.Altura})
		puntos = append(puntos, Punto{e.Xf, false, e.Altura})
	}

	// Ordenar puntos por posición x, entradas antes que salidas, mayor altura primero
	sort.Slice(puntos, func(i, j int) bool {
		if puntos[i].x == puntos[j].x {
			if puntos[i].esIni == puntos[j].esIni {
				if puntos[i].esIni {
					return puntos[i].height > puntos[j].height
				} else {
					return puntos[i].height < puntos[j].height
				}
			}
			return puntos[i].esIni
		}
		return puntos[i].x < puntos[j].x
	})

	heights := map[float64]int{0: 1} // Alturas activas con conteo
	maxHeight := 0.0
	var resultado []string

	for _, p := range puntos {
		if p.esIni {
			heights[p.height]++
			if p.height > maxHeight {
				maxHeight = p.height
				resultado = append(resultado, fmt.Sprintf("%.0f %.0f", p.x, p.height))
			}
		} else {
			heights[p.height]--
			if heights[p.height] == 0 {
				delete(heights, p.height)
			}
			nuevoMax := 0.0
			for h := range heights {
				if h > nuevoMax {
					nuevoMax = h
				}
			}
			if nuevoMax != maxHeight {
				maxHeight = nuevoMax
				resultado = append(resultado, fmt.Sprintf("%.0f %.0f", p.x, maxHeight))
			}
		}
	}

	// Guardar en archivo
	f, err := os.Create(nombreArchivo)
	if err != nil {
		log.Fatalf("Error creando archivo de horizonte: %v", err)
	}
	defer f.Close()
	for _, linea := range resultado {
		f.WriteString(linea + "")
	}

	utils.PrintLog(fmt.Sprint("✅ Línea del horizonte guardada en:", nombreArchivo))
}
