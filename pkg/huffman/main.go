package huffman

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
	"sigoa/internal/utils"
	"strings"
)

// --- Estructuras para la Codificación Huffman ---

// HuffmanNode representa un nodo en el árbol de Huffman
type HuffmanNode struct {
	char  rune
	freq  int
	left  *HuffmanNode
	right *HuffmanNode
}

// IsLeaf verifica si el nodo es una hoja
func (n *HuffmanNode) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

// PriorityQueue implementa heap.Interface para los nodos de Huffman
type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].freq < pq[j].freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*HuffmanNode)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// --- Funciones de Codificación Huffman ---

// buildFrequencyMap calcula la frecuencia de cada carácter en el string de entrada
func buildFrequencyMap(data string) map[rune]int {
	freqMap := make(map[rune]int)
	for _, char := range data {
		freqMap[char]++
	}
	return freqMap
}

// buildHuffmanTree construye el árbol de Huffman a partir del mapa de frecuencias
func buildHuffmanTree(freqMap map[rune]int) *HuffmanNode {
	pq := make(PriorityQueue, 0, len(freqMap))
	for char, freq := range freqMap {
		pq.Push(&HuffmanNode{char: char, freq: freq})
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		node1 := heap.Pop(&pq).(*HuffmanNode)
		node2 := heap.Pop(&pq).(*HuffmanNode)
		newNode := &HuffmanNode{
			char:  0, // No es una hoja, no tiene un carácter asociado
			freq:  node1.freq + node2.freq,
			left:  node1,
			right: node2,
		}
		heap.Push(&pq, newNode)
	}
	return heap.Pop(&pq).(*HuffmanNode)
}

// buildHuffmanCodes genera los códigos binarios para cada carácter
func buildHuffmanCodes(root *HuffmanNode, currentCode string, codes map[rune]string) {
	if root == nil {
		return
	}

	if root.IsLeaf() {
		codes[root.char] = currentCode
		return
	}

	buildHuffmanCodes(root.left, currentCode+"0", codes)
	buildHuffmanCodes(root.right, currentCode+"1", codes)
}

// encodeString codifica el string usando los códigos de Huffman
func encodeString(data string, codes map[rune]string) string {
	var encodedBuilder strings.Builder
	for _, char := range data {
		encodedBuilder.WriteString(codes[char])
	}
	return encodedBuilder.String()
}

// serializeTree serializa el árbol de Huffman para guardarlo junto con los datos
func serializeTree(root *HuffmanNode) string {
	if root == nil {
		return ""
	}
	if root.IsLeaf() {
		// Formato: L<char_as_int>
		return fmt.Sprintf("L%d", root.char)
	}
	// Formato: I(<left_subtree>)(<right_subtree>)
	return fmt.Sprintf("I(%s)(%s)", serializeTree(root.left), serializeTree(root.right))
}

// deserializeTree deserializa el árbol de Huffman desde una cadena
func deserializeTree(reader *strings.Reader) (*HuffmanNode, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		return nil, err
	}

	switch char {
	case 'L': // Leaf node
		var charVal int
		_, err := fmt.Fscanf(reader, "%d", &charVal)
		if err != nil {
			return nil, err
		}
		return &HuffmanNode{char: rune(charVal), freq: 0}, nil
	case 'I': // Internal node
		_, _, err := reader.ReadRune() // Consume '('
		if err != nil {
			return nil, err
		}
		left, err := deserializeTree(reader)
		if err != nil {
			return nil, err
		}
		_, _, err = reader.ReadRune() // Consume ')'
		if err != nil {
			return nil, err
		}

		_, _, err = reader.ReadRune() // Consume '('
		if err != nil {
			return nil, err
		}
		right, err := deserializeTree(reader)
		if err != nil {
			return nil, err
		}
		_, _, err = reader.ReadRune() // Consume ')'
		if err != nil {
			return nil, err
		}
		return &HuffmanNode{left: left, right: right}, nil
	default:
		return nil, fmt.Errorf("formato de árbol de Huffman inválido: %c", char)
	}
}

// decodeString decodifica el string binario usando el árbol de Huffman
func decodeString(encodedData string, root *HuffmanNode) string {
	var decodedBuilder strings.Builder
	current := root
	for _, bit := range encodedData {
		if bit == '0' {
			current = current.left
		} else if bit == '1' {
			current = current.right
		}

		if current.IsLeaf() {
			decodedBuilder.WriteRune(current.char)
			current = root // Reiniciar desde la raíz para el siguiente carácter
		}
	}
	return decodedBuilder.String()
}

// --- Funciones de Guardado y Lectura de Archivos ---

// Guardar string en un archivo .huff
func Guardar(data string, fileName string) (string, error) {
	// Generar un nombre de archivo único
	//fileName := fmt.Sprintf("data_%d.huff", time.Now().UnixNano())
	//filePath := filepath.Join("./output/", fileName) // Guardar en el directorio actual
	//fmt.Println("Guardando archivo:", fileName)
	// 1. Construir el árbol y los códigos de Huffman
	freqMap := buildFrequencyMap(data)
	if len(freqMap) == 0 {
		return "", fmt.Errorf("no hay datos para codificar")
	}
	huffmanTree := buildHuffmanTree(freqMap)
	huffmanCodes := make(map[rune]string)
	buildHuffmanCodes(huffmanTree, "", huffmanCodes)

	// 2. Codificar el string
	encodedData := encodeString(data, huffmanCodes)

	// 3. Serializar el árbol de Huffman
	serializedTreeStr := serializeTree(huffmanTree)

	// Abrir el archivo para escritura
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Guardar el árbol serializado y los datos codificados
	// Separamos el árbol serializado y los datos codificados con un salto de línea
	// para facilitar la lectura.
	_, err = writer.WriteString(serializedTreeStr + "")
	if err != nil {
		return "", fmt.Errorf("error al escribir el árbol serializado: %w", err)
	}
	_, err = writer.WriteString(encodedData)
	if err != nil {
		return "", fmt.Errorf("error al escribir los datos codificados: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("error al vaciar el buffer del escritor: %w", err)
	}

	return fileName, nil
}

// Leer lee un archivo .huff, lo desencripta e imprime por pantalla
func Leer(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Leer el árbol serializado (hasta el primer salto de línea)
	serializedTreeStr, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error al leer el árbol serializado: %w", err)
	}
	serializedTreeStr = strings.TrimSpace(serializedTreeStr) // Eliminar el salto de línea

	// Leer los datos codificados restantes
	encodedDataBuilder := new(strings.Builder)
	_, err = io.Copy(encodedDataBuilder, reader)
	if err != nil {
		return fmt.Errorf("error al leer los datos codificados: %w", err)
	}
	encodedData := encodedDataBuilder.String()

	// Deserializar el árbol de Huffman
	treeReader := strings.NewReader(serializedTreeStr)
	huffmanTree, err := deserializeTree(treeReader)
	if err != nil {
		return fmt.Errorf("error al deserializar el árbol de Huffman: %w", err)
	}

	// Decodificar los datos
	decodedString := decodeString(encodedData, huffmanTree)

	utils.PrintLog(fmt.Sprint("--- Contenido Desencriptado ---"))
	utils.PrintLog(fmt.Sprint(decodedString))
	utils.PrintLog(fmt.Sprint("------------------------------"))

	return nil
}
