package cargomanagement

import (
	"fmt"

	"github.com/example/sigoa/internal/models"
)

// CheckCargoRestrictions verifica si un paquete de carga puede ser añadido a un vuelo
// sin exceder las restricciones de peso y volumen de la aeronave.
func CheckCargoRestrictions(flight *models.Flight, cargoPackage models.CargoPackage) error {
	if flight == nil {
		return fmt.Errorf("el vuelo proporcionado es nulo")
	}
	if flight.Aircraft == nil {
		return fmt.Errorf("la aeronave para el vuelo %s no está definida", flight.ID)
	}

	currentWeight := 0.0
	currentVolume := 0.0
	for _, cp := range flight.Cargo {
		currentWeight += cp.Weight
		currentVolume += cp.Volume
	}

	if currentWeight+cargoPackage.Weight > flight.Aircraft.MaxCargoWeight {
		return fmt.Errorf("excede el peso máximo de carga permitido para la aeronave (%s): %.2f kg", flight.Aircraft.Model, flight.Aircraft.MaxCargoWeight)
	}
	if currentVolume+cargoPackage.Volume > flight.Aircraft.MaxCargoVolume {
		return fmt.Errorf("excede el volumen máximo de carga permitido para la aeronave (%s): %.2f m³", flight.Aircraft.Model, flight.Aircraft.MaxCargoVolume)
	}

	return nil
}

func AddCargoToFlight(flightID string, cargoPackage models.CargoPackage, flights map[string]*models.Flight) error {
	flight, exists := flights[flightID]
	if !exists {
		return fmt.Errorf("el vuelo con ID %s no existe", flightID)
	}

// Devuelve un slice de CargoPackage y un error si el vuelo no existe.
func GetCargoForFlight(flightID string, flights map[string]*models.Flight) ([]models.CargoPackage, error) {
	flight, exists := flights[flightID]
	if !exists {
		return nil, fmt.Errorf("el vuelo con ID %s no existe", flightID)
	}

	return flight.Cargo, nil
}

// Devuelve un error si el vuelo o el paquete de carga no existen.
func RemoveCargoFromFlight(flightID string, cargoID string, flights map[string]*models.Flight) error {
	flight, exists := flights[flightID]
	if !exists {
		return fmt.Errorf("el vuelo con ID %s no existe", flightID)
	}

	foundIndex := -1
	for i, cargo := range flight.Cargo {
		if cargo.ID == cargoID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("el paquete de carga con ID %s no fue encontrado en el vuelo %s", cargoID, flightID)
	}

	// Eliminar el elemento de la slice sin mantener el orden
	flight.Cargo = append(flight.Cargo[:foundIndex], flight.Cargo[foundIndex+1:]...)
	fmt.Printf("Paquete de carga (ID: %s) removido del vuelo %s\n", cargoID, flightID)
	return nil
}

// Puedes añadir otras funciones relevantes para la gestión de carga aquí.
