package main

import (
	"fmt"
	"os" // Importar el paquete os
	"sigoa/pkg/huffman"
)

func main() {
	// Verificar si se proporcionó un nombre de archivo como argumento
	if len(os.Args) < 2 {
		utils.PrintLog(fmt.Sprint("Uso: go run tu_programa.go <nombre_del_archivo.huff>"))
		utils.PrintLog(fmt.Sprint("Ejemplo: go run main.go data_1700000000000000000.huff"))
		return
	}

	// El nombre del archivo es el primer argumento después del nombre del programa
	filename := os.Args[1]

	// Leer y desencriptar el archivo
	err := huffman.Leer(filename) // Asegúrate de que 'err' sea una nueva variable aquí
	if err != nil {
		utils.PrintLog(fmt.Sprintf("Error al leer el archivo '%s': %v", filename, err)
		return
	}

	utils.PrintLog(fmt.Sprintf("Archivo '%s' leído y desencriptado con éxito.", filename)
}
