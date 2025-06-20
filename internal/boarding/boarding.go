package main

import (
	"fmt"
	"time"

	"internal/models"
)

// SimulateBoarding simulates the boarding process for a given flight.
func SimulateBoarding(flight *models.Flight, checkedInPassengers map[string]*models.Passenger) {
	fmt.Printf("Starting boarding for flight %s...\n", flight.FlightNumber)

	// Assuming boarding starts some time after check-in closes
	time.Sleep(15 * time.Minute) // Simulate time between check-in closing and boarding starting

	boardingQueue := make(chan *models.Passenger, len(flight.AssignedPassengers))
	boardedCount := 0

	// Add checked-in passengers to the boarding queue
	for _, passengerID := range flight.AssignedPassengers {
		if passenger, ok := checkedInPassengers[passengerID]; ok && passenger.IsCheckedIn {
			boardingQueue <- passenger
		}
	}
	close(boardingQueue)

	// Simulate boarding at a single gate for simplicity
	for passenger := range boardingQueue {
		fmt.Printf("Passenger %s %s boarding flight %s...\n", passenger.FirstName, passenger.LastName, flight.FlightNumber)

		// Verify ticket - in a real system, this would involve scanning
		if passenger.FlightNumber != flight.FlightNumber {
			fmt.Printf("Error: Passenger %s %s trying to board wrong flight %s\n", passenger.FirstName, passenger.LastName, passenger.FlightNumber)
			continue // Skip boarding for this passenger
		}

		// Simulate boarding time
		time.Sleep(1 * time.Second) // Simulate scanning ticket and walking

		passenger.IsBoarded = true
		boardedCount++
		fmt.Printf("Passenger %s %s boarded.\n", passenger.FirstName, passenger.LastName)
	}

	flight.BoardedPassengers = boardedCount
	flight.Status = "Boarding Complete"
	fmt.Printf("Boarding complete for flight %s. %d passengers boarded.\n", flight.FlightNumber, boardedCount)
}

// Note: This is a simplified implementation. In a real system, you would
// likely have multiple boarding gates, more sophisticated queuing, and
// integration with a real-time flight status system.
// The function assumes `checkedInPassengers` contains all passengers who
// successfully checked in.