package generator

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"sigoa/internal/models"
)

// GenerateFlightData crea un archivo CSV con datos de vuelos de ejemplo.
func GenerateFlightData(filename string, numFlights int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo de vuelos: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezado
	writer.Write([]string{"numero_vuelo", "origen", "destino", "fecha_salida", "hora_salida", "capacidad"})

	// Generar datos de vuelos
	for i := 1; i <= numFlights; i++ {
		flightNumber := fmt.Sprintf("AR%03d", 100+i)
		origin := fmt.Sprintf("APT%d", 1+i%5)
		destination := fmt.Sprintf("APT%d", 6+i%5)
		departureDate := time.Now().AddDate(0, 0, i-1).Format("2006-01-02")
		departureTime := fmt.Sprintf("%02d:%02d", 8+(i%12), 0)
		capacity := 100 + (i%5)*50

		writer.Write([]string{flightNumber, origin, destination, departureDate, departureTime, strconv.Itoa(capacity)})
	}

	return nil
}

// GeneratePassengerData crea un archivo CSV con datos de pasajeros de ejemplo.
func GeneratePassengerData(filename string, numPassengers int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo de pasajeros: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezado
	writer.Write([]string{"documento", "nombre", "apellido", "categoria"})

	categories := []string{"platino", "oro", "plata", "no frecuente"}

	// Generar datos de pasajeros
	for i := 1; i <= numPassengers; i++ {
		document := fmt.Sprintf("DNI%07d", 1000000+i)
		firstName := fmt.Sprintf("Nombre%d", i)
		lastName := fmt.Sprintf("Apellido%d", i)
		category := categories[i%len(categories)]

		writer.Write([]string{document, firstName, lastName, category})
	}

	return nil
}

// GenerateBookingData crea un archivo CSV con datos de reservas de ejemplo.
// Asume que ya existen datos de vuelos y pasajeros.
func GenerateBookingData(filename string, flights []models.Flight, passengers []models.Passenger, numBookings int) error {
	if len(flights) == 0 || len(passengers) == 0 {
		return fmt.Errorf("no hay datos de vuelos o pasajeros para generar reservas")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo de reservas: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezado
	writer.Write([]string{"numero_reserva", "numero_vuelo", "documento_pasajero_1", "documento_pasajero_2", "documento_pasajero_3", "documento_pasajero_4", "documento_pasajero_5"})

	// Generar datos de reservas
	for i := 1; i <= numBookings; i++ {
		bookingNumber := fmt.Sprintf("RES%05d", 10000+i)
		flight := flights[i%len(flights)]
		row := []string{bookingNumber, flight.FlightNumber}

		// Añadir hasta 5 pasajeros por reserva
		for j := 0; j < (i%5)+1; j++ {
			passenger := passengers[(i*5+j)%len(passengers)]
			row = append(row, passenger.Document)
		}
		// Rellenar con campos vacíos si hay menos de 5 pasajeros
		for len(row) < 7 {
			row = append(row, "")
		}

		writer.Write(row)
	}

	return nil
}

// GenerateCargoData crea un archivo CSV con datos de carga de ejemplo.
func GenerateCargoData(filename string, flights []models.Flight, numCargoPackages int) error {
	if len(flights) == 0 {
		return fmt.Errorf("no hay datos de vuelos para generar carga")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo de carga: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezado
	writer.Write([]string{"id_paquete", "numero_vuelo", "peso_kg", "volumen_m3", "destino"})

	// Generar datos de carga
	for i := 1; i <= numCargoPackages; i++ {
		packageID := fmt.Sprintf("CARGO%06d", 100000+i)
		flight := flights[i%len(flights)]
		weight := 10 + (i%10)*5
		volume := 0.1 + float64(i%5)*0.05
		destination := flight.Destination // La carga va al destino del vuelo

		writer.Write([]string{packageID, flight.FlightNumber, strconv.Itoa(weight), fmt.Sprintf("%.2f", volume), destination})
	}

	return nil
}

// Puedes añadir funciones similares para generar datos de aeropuertos, aeronaves, etc.
// Asegúrate de ajustar los campos y el formato según las necesidades de tu sistema.