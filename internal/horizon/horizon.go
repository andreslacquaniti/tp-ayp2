package horizon

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Estructura que representa un edificio con su posición inicial (x1),
// altura y posición final (x2)
type Edificio struct {
	x1, altura, x2 int
}

// Punto del perfil del horizonte con una coordenada (x, y)
// donde y representa la altura visible en x
type Punto struct {
	x, y int
}

// ------------------------------------------------------------
// LECTURA DE DATOS DESDE ARCHIVO
// ------------------------------------------------------------

// Lee edificios desde un archivo de texto con formato: x1;altura;x2 por línea
func leerEdificios(ruta string) ([]Edificio, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	var edificios []Edificio

	for scanner.Scan() {
		partes := strings.Split(scanner.Text(), ";")
		if len(partes) != 3 {
			continue // si la línea no tiene el formato correcto, se salta
		}

		// Convertimos los valores a enteros
		x1, _ := strconv.Atoi(partes[0])
		h, _ := strconv.Atoi(partes[1])
		x2, _ := strconv.Atoi(partes[2])

		edificios = append(edificios, Edificio{x1, h, x2})
	}

	return edificios, nil
}

// ------------------------------------------------------------
// FUNCIÓN DE COMBINACIÓN (MERGE) DE DOS PERFILES DE HORIZONTE
// ------------------------------------------------------------

// Combina dos perfiles del horizonte en uno solo manteniendo la forma correcta
func merge(left, right []Punto) []Punto {
	var resultado []Punto
	h1, h2 := 0, 0 // alturas actuales de cada mitad
	i, j := 0, 0   // índices de recorrido para ambos perfiles

	for i < len(left) && j < len(right) {
		var x int

		// Comparar coordenadas x de ambos perfiles
		if left[i].x < right[j].x {
			x = left[i].x
			h1 = left[i].y
			i++
		} else if left[i].x > right[j].x {
			x = right[j].x
			h2 = right[j].y
			j++
		} else {
			// Si ambos tienen el mismo x, tomamos ambos y avanzamos
			x = left[i].x
			h1 = left[i].y
			h2 = right[j].y
			i++
			j++
		}

		// Se toma la mayor altura visible entre los dos perfiles
		maxH := max(h1, h2)

		// Agregamos el punto solo si cambia la altura respecto al último punto agregado
		if len(resultado) == 0 || resultado[len(resultado)-1].y != maxH {
			resultado = append(resultado, Punto{x, maxH})
		}
	}

	// Añadir puntos restantes que no se recorrieron del todo
	for ; i < len(left); i++ {
		if len(resultado) == 0 || resultado[len(resultado)-1].y != left[i].y {
			resultado = append(resultado, left[i])
		}
	}

	for ; j < len(right); j++ {
		if len(resultado) == 0 || resultado[len(resultado)-1].y != right[j].y {
			resultado = append(resultado, right[j])
		}
	}

	return resultado
}

// ------------------------------------------------------------
// FUNCIÓN PRINCIPAL DE DIVISIÓN Y CONQUISTA
// ------------------------------------------------------------

// Aplica recursivamente el algoritmo para generar la línea del horizonte
func lineaHorizonte(edificios []Edificio) []Punto {
	// Caso base: sin edificios
	if len(edificios) == 0 {
		return []Punto{}
	}

	// Caso base: un solo edificio
	if len(edificios) == 1 {
		e := edificios[0]
		return []Punto{
			{e.x1, e.altura}, // Inicio del edificio
			{e.x2, 0},        // Fin del edificio (vuelve a nivel 0)
		}
	}

	// Dividir en dos mitades
	mid := len(edificios) / 2
	izquierda := lineaHorizonte(edificios[:mid])
	derecha := lineaHorizonte(edificios[mid:])

	// Combinar los resultados
	return merge(izquierda, derecha)
}

// ------------------------------------------------------------
// GUARDAR RESULTADO EN ARCHIVO DE SALIDA
// ------------------------------------------------------------

// Guarda los puntos resultantes del horizonte en un archivo de texto
func guardarLineaHorizonte(puntos []Punto, rutaSalida string) error {
	os.MkdirAll("output", os.ModePerm) // crea carpeta si no existe
	archivo, err := os.Create(rutaSalida)
	if err != nil {
		return err
	}
	defer archivo.Close()

	writer := bufio.NewWriter(archivo)
	for _, punto := range puntos {
		linea := fmt.Sprintf("%d;%d\n", punto.x, punto.y)
		writer.WriteString(linea)
	}

	return writer.Flush()
}

// Función auxiliar para máximo entre dos enteros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ------------------------------------------------------------
// FUNCIÓN MAIN: entrada del programa
// ------------------------------------------------------------

func main() {
	// Leer edificios desde archivo de entrada
	edificios, err := leerEdificios("data/edificios.txt")
	if err != nil {
		fmt.Println("Error leyendo edificios:", err)
		return
	}

	// Ordenar los edificios por su punto inicial x1 (importante para D&C)
	sort.Slice(edificios, func(i, j int) bool {
		return edificios[i].x1 < edificios[j].x1
	})

	// Obtener la línea del horizonte con D&C
	skyline := lineaHorizonte(edificios)

	// Guardar resultados en archivo de salida
	err = guardarLineaHorizonte(skyline, "output/linea_horizonte.txt")
	if err != nil {
		fmt.Println("Error guardando línea del horizonte:", err)
		return
	}

	fmt.Println("Línea del horizonte generada correctamente.")
}
