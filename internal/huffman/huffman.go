package huffman

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// Nodo del árbol de Huffman
type nodoHuffman struct {
	caracter byte
	frecuencia int
	izquierda *nodoHuffman
	derecha *nodoHuffman
}

// Implementación de la interfaz heap.Interface para la cola de prioridad
type colaPrioridadNodos []*nodoHuffman

func (cp colaPrioridadNodos) Len() int { return len(cp) }
func (cp colaPrioridadNodos) Less(i, j int) bool { return cp[i].frecuencia < cp[j].frecuencia }
func (cp colaPrioridadNodos) Swap(i, j int) { cp[i], cp[j] = cp[j], cp[i] }

func (cp *colaPrioridadNodos) Push(x interface{}) {
	nodo := x.(*nodoHuffman)
	*cp = append(*cp, nodo)
}

func (cp *colaPrioridadNodos) Pop() interface{} {
	old := *cp
	n := len(old)
	nodo := old[n-1]
	*cp = old[0 : n-1]
	return nodo
}

// Construye el árbol de Huffman a partir de un mapa de frecuencias de caracteres
func construirArbolHuffman(frecuencias map[byte]int) *nodoHuffman {
	cola := &colaPrioridadNodos{}
	heap.Init(cola)

	for caracter, frecuencia := range frecuencias {
		heap.Push(cola, &nodoHuffman{caracter: caracter, frecuencia: frecuencia})
	}

	for cola.Len() > 1 {
		nodo1 := heap.Pop(cola).(*nodoHuffman)
		nodo2 := heap.Pop(cola).(*nodoHuffman)

		nodoCombinado := &nodoHuffman{
			frecuencia: nodo1.frecuencia + nodo2.frecuencia,
			izquierda:  nodo1,
			derecha:    nodo2,
		}
		heap.Push(cola, nodoCombinado)
	}

	if cola.Len() == 0 {
		return nil // Handle empty input
	}
	return heap.Pop(cola).(*nodoHuffman)
}

// Genera los códigos de Huffman a partir del árbol
func generarCodigos(arbol *nodoHuffman, prefijo string, codigos map[byte]string) {
	if arbol == nil {
		return
	}

	if arbol.izquierda == nil && arbol.derecha == nil {
		codigos[arbol.caracter] = prefijo
		return
	}

	generarCodigos(arbol.izquierda, prefijo+"0", codigos)
	generarCodigos(arbol.derecha, prefijo+"1", codigos)
}

// Comprime los datos utilizando los códigos de Huffman
func comprimirDatos(datos []byte, codigos map[byte]string) ([]byte, error) {
	var bits string
	for _, b := range datos {
		codigo, ok := codigos[b]
		if !ok {
			return nil, fmt.Errorf("caracter '%c' sin código Huffman", b)
		}
		bits += codigo
	}

	// Rellenar con ceros al final para completar el último byte
	for len(bits)%8 != 0 {
		bits += "0"
	}

	bytesComprimidos := make([]byte, len(bits)/8)
	for i := 0; i < len(bits); i += 8 {
		byteStr := bits[i : i+8]
		var byteVal byte
		for j := 0; j < 8; j++ {
			if byteStr[j] == '1' {
				byteVal |= (1 << (7 - j))
			}
		}
		bytesComprimidos[i/8] = byteVal
	}

	return bytesComprimidos, nil
}

// Comprime un archivo utilizando Huffman y guarda el resultado
func ComprimirArchivo(rutaArchivo string, rutaSalida string) error {
	datos, err := ioutil.ReadFile(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error al leer el archivo: %w", err)
	}

	// Calcular frecuencias de caracteres
	frecuencias := make(map[byte]int)
	for _, b := range datos {
		frecuencias[b]++
	}

	// Construir árbol de Huffman
	arbol := construirArbolHuffman(frecuencias)

	// Generar códigos de Huffman
	codigos := make(map[byte]string)
	generarCodigos(arbol, "", codigos)

	// Comprimir datos
	datosComprimidos, err := comprimirDatos(datos, codigos)
	if err != nil {
		return fmt.Errorf("error al comprimir los datos: %w", err)
	}

	// Guardar datos comprimidos
	err = ioutil.WriteFile(rutaSalida, datosComprimidos, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo comprimido: %w", err)
	}

	// Opcional: guardar también el árbol de Huffman o la tabla de códigos para la descompresión
	// Esto requeriría definir un formato para serializar el árbol o la tabla.
	// Por simplicidad, asumimos que el descompresor puede reconstruir el árbol si lee las frecuencias o si se incluye la tabla de códigos en el archivo comprimido.
	// Una forma común es incluir la tabla de frecuencias o la tabla de códigos al inicio del archivo comprimido.
	// Para este ejemplo, no estamos serializando el árbol o la tabla.

	fmt.Printf("Archivo comprimido con éxito en: %s\n", rutaSalida)
	return nil
}
