package main

import (
	"fmt"
	"os" // Importar el paquete os
	"sigoa/pkg/huffman"
)

func main() {
	// Verificar si se proporcionó un nombre de archivo como argumento
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run tu_programa.go <nombre_del_archivo.huff>")
		fmt.Println("Ejemplo: go run main.go data_1700000000000000000.huff")
		return
	}

	// El nombre del archivo es el primer argumento después del nombre del programa
	filename := os.Args[1]

	// Leer y desencriptar el archivo
	err := huffman.Leer(filename) // Asegúrate de que 'err' sea una nueva variable aquí
	if err != nil {
		fmt.Printf("Error al leer el archivo '%s': %v\n", filename, err)
		return
	}

	fmt.Printf("Archivo '%s' leído y desencriptado con éxito.\n", filename)
}
