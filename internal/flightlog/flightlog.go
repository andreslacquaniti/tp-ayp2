package app

import (
	"fmt"
	"os"
	"time"

	"github.com/example/sigoa/internal/models"
)

// GenerateFlightLog generates a text file containing the flight log information.
func GenerateFlightLog(flight *models.Flight) error {
	logFileName := fmt.Sprintf("flight_log_%s_%s.txt", flight.Number, flight.Schedule.Format("200601021504"))
	file, err := os.Create(logFileName)
	if err != nil {
		return fmt.Errorf("failed to create flight log file: %w", err)
	}
	defer file.Close()

	// Write Flight Details
	fmt.Fprintf(file, "--- Flight Details ---\n")
	fmt.Fprintf(file, "Flight Number: %s\n", flight.Number)
	fmt.Fprintf(file, "Origin: %s\n", flight.Origin)
	fmt.Fprintf(file, "Destination: %s\n", flight.Destination)
	fmt.Fprintf(file, "Schedule: %s\n", flight.Schedule.Format(time.RFC3339))
	fmt.Fprintf(file, "Aircraft: %s\n", flight.Aircraft.Model)
	fmt.Fprintf(file, "Status: %s\n", flight.Status)
	fmt.Fprintf(file, "\n")

	// Write Checked-in Passengers
	fmt.Fprintf(file, "--- Checked-in Passengers ---\n")
	if len(flight.Passengers) == 0 {
		fmt.Fprintf(file, "No passengers checked in.\n")
	} else {
		for _, passenger := range flight.Passengers {
			if passenger.CheckinInfo != nil {
				fmt.Fprintf(file, "  - Document: %s, Name: %s %s, Category: %s, Boarded: %t\n",
					passenger.DocumentID, passenger.FirstName, passenger.LastName, passenger.Category, passenger.Boarded)
			}
		}
	}
	fmt.Fprintf(file, "\n")

	// Write Boarded Passengers
	fmt.Fprintf(file, "--- Boarded Passengers ---\n")
	boardedCount := 0
	for _, passenger := range flight.Passengers {
		if passenger.Boarded {
			fmt.Fprintf(file, "  - Document: %s, Name: %s %s, Category: %s\n",
				passenger.DocumentID, passenger.FirstName, passenger.LastName, passenger.Category)
			boardedCount++
		}
	}
	if boardedCount == 0 {
		fmt.Fprintf(file, "No passengers boarded.\n")
	}
	fmt.Fprintf(file, "\n")

	// Write Baggage Information
	fmt.Fprintf(file, "--- Baggage Information ---\n")
	totalBaggageWeight := 0.0
	totalBaggagePieces := 0
	hasBaggage := false
	for _, passenger := range flight.Passengers {
		if passenger.CheckinInfo != nil && passenger.CheckinInfo.Baggage != nil {
			hasBaggage = true
			fmt.Fprintf(file, "  - Passenger Document: %s, Pieces: %d, Total Weight: %.2f kg\n",
				passenger.DocumentID, passenger.CheckinInfo.Baggage.Pieces, passenger.CheckinInfo.Baggage.TotalWeight)
			totalBaggageWeight += passenger.CheckinInfo.Baggage.TotalWeight
			totalBaggagePieces += passenger.CheckinInfo.Baggage.Pieces
		}
	}
	if !hasBaggage {
		fmt.Fprintf(file, "No baggage checked in.\n")
	}
	fmt.Fprintf(file, "Total Baggage Pieces: %d\n", totalBaggagePieces)
	fmt.Fprintf(file, "Total Baggage Weight: %.2f kg\n", totalBaggageWeight)
	fmt.Fprintf(file, "\n")

	// Write Cargo Information (assuming cargo is associated with the flight)
	// You would need a structure for Cargo in models and potentially a slice of cargo in the Flight struct
	fmt.Fprintf(file, "--- Cargo Information ---\n")
	if len(flight.Cargo) == 0 { // Assuming Flight struct has a []models.Cargo field
		fmt.Fprintf(file, "No cargo loaded.\n")
	} else {
		for _, cargo := range flight.Cargo { // Assuming cargo has fields like ID, Weight, Volume, Destination
			fmt.Fprintf(file, "  - Cargo ID: %s, Weight: %.2f kg, Volume: %.2f m³, Destination: %s\n",
				cargo.ID, cargo.Weight, cargo.Volume, cargo.Destination)
		}
	}
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "--- End of Log ---\n")

	return nil
}

// Note: This code assumes that the `models.Flight` struct has fields
// for `Passengers` ([]models.Passenger) and potentially `Cargo` ([]models.Cargo).
// The `models.Passenger` struct is assumed to have fields like `DocumentID`,
// `FirstName`, `LastName`, `Category`, `Boarded` (bool), and `CheckinInfo` (*models.CheckinDetails).
// The `models.CheckinDetails` struct is assumed to have a `Baggage` field (*models.Baggage).
// The `models.Baggage` struct is assumed to have `Pieces` (int) and `TotalWeight` (float64).
// The `models.Cargo` struct is assumed to have fields like `ID`, `Weight`, `Volume`, and `Destination`.