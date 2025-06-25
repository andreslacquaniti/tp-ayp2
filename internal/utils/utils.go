package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sigoa/internal/app"
	"strings"
	"time"
)

// realStartTime almacenará la hora real en que la aplicación comenzó a ejecutarse.
var realStartTime time.Time

func init() {
	// Se ejecuta automáticamente al inicio del paquete.
	// Capturamos el tiempo real de inicio una sola vez.
	realStartTime = time.Now()
}

// PrintLog imprime el mensaje en la consola y lo guarda en un archivo de log.
func PrintLog(mensaje string) {
	// Imprimir en consola usando la HoraSistema simulada (como ya lo tienes)
	logLine := fmt.Sprintf("%s | %s", app.HoraSistema.Format("15:04"), mensaje)
	fmt.Println(logLine)

	// Crear el directorio 'out' si no existe
	logDir := "output"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755) // Permisos de lectura/escritura para el dueño, solo lectura para otros
		if err != nil {
			log.Printf("Error al crear el directorio de logs '%s': %v", logDir, err)
			return // No podemos continuar sin el directorio
		}
	}

	// Generar el nombre del archivo de log con la FECHA Y HORA REAL DE INICIO DEL PROGRAMA
	// Formato: YYYYMMDD_HHMMSS_proceso.log
	fileName := fmt.Sprintf("%s_proceso.log", realStartTime.Format("20060102_150405"))
	logFilePath := filepath.Join(logDir, fileName)

	// Abrir el archivo en modo de añadir (append). Si no existe, lo crea.
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error al abrir/crear el archivo de log '%s': %v", logFilePath, err)
		return
	}
	defer file.Close() // Asegura que el archivo se cierre al salir de la función

	// Escribir el mensaje en el archivo, seguido de un salto de línea
	_, err = file.WriteString(logLine + "\n")
	if err != nil {
		log.Printf("Error al escribir en el archivo de log '%s': %v", logFilePath, err)
	}
}

func GeneraNroTicket() string {
	var caracteresBase62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	longitudHash := 6
	var hashBuilder strings.Builder
	hashBuilder.Grow(longitudHash) // Pre-asigna espacio para eficiencia
	for range longitudHash {
		// Genera un número aleatorio seguro entre 0 (inclusive) y 62 (exclusive)
		// que corresponde al índice de los caracteres en caracteresBase62.
		numAleatorio, err := rand.Int(rand.Reader, big.NewInt(int64(len(caracteresBase62))))
		if err != nil {
			return GeneraNroTicket()
		}
		// Agrega el carácter correspondiente al builder
		hashBuilder.WriteByte(caracteresBase62[numAleatorio.Int64()])
	}
	return hashBuilder.String()
}
