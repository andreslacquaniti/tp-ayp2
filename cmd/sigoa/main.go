package main

import (
	"fmt"
	"time"

	"github.com/Santy-J/sigoa/internal/app"
	"github.com/Santy-J/sigoa/internal/checkin"
	"github.com/Santy-J/sigoa/internal/models"
)

func main() {
	// Load data from CSV files
	aeropuertos, err := app.LoadAeropuertos("data/aeropuertos.txt")
	if err != nil {
		fmt.Println("Error loading airports:", err)
		return
	}

	vuelos, err := app.LoadVuelos("data/vuelos.txt")
	if err != nil {
		fmt.Println("Error loading flights:", err)
		return
	}

	aeronaves, err := app.LoadAeronaves("data/aeronaves.txt")
	if err != nil {
		fmt.Println("Error loading aircrafts:", err)
		return
	}

	clientes, err := app.LoadClientes("data/clientes.txt")
	if err != nil {
		fmt.Println("Error loading clients:", err)
		return
	}

	reservas, err := app.LoadReservas("data/reservas.txt")
	if err != nil {
		fmt.Println("Error loading reservations:", err)
		return
	}

	configuracionAsientos, err := app.LoadConfiguracionAsientos("data/configuracion_asientos.txt")
	if err != nil {
		fmt.Println("Error loading seat configuration:", err)
		return
	}

	edificios, err := app.LoadEdificios("data/edificios.txt")
	if err != nil {
		fmt.Println("Error loading buildings:", err)
		return
	}

	cargas, err := app.LoadCargas("data/cargas.txt")
	if err != nil {
		fmt.Println("Error loading cargo:", err)
		return
	}

	// Initialize the check-in system
	checkinSystem := checkin.NewCheckinSystem(vuelos, clientes, reservas, configuracionAsientos, buildings)

	// Simulate passenger check-in for a specific flight (example)
	flightNumber := "AR123"
	flight, found := checkinSystem.FindFlight(flightNumber)
	if !found {
		fmt.Printf("Flight %s not found.\n", flightNumber)
		return
	}

	fmt.Printf("Starting check-in simulation for flight %s...\n", flightNumber)

	// Simulate passengers arriving and queuing
	passengerQueue := make(chan models.Reservation)
	go func() {
		for _, reservation := range reservations {
			if reservation.FlightCode == flightNumber {
				passengerQueue <- reservation
				time.Sleep(time.Duration(reservation.ArrivalTime) * time.Millisecond) // Simulate arrival time
			}
		}
		close(passengerQueue)
	}()

	// Simulate check-in counters
	numCounters := 3 // Example number of counters
	counterChannel := make(chan int, numCounters)
	for i := 1; i <= numCounters; i++ {
		counterChannel <- i // Initialize counters as available
	}

	// Process passengers from the queue
	for reservation := range passengerQueue {
		counterID := <-counterChannel // Wait for an available counter
		fmt.Printf("Passenger %s (Reservation ID: %s) arriving at Counter %d.\n", reservation.ClientID, reservation.ReservationID, counterID)

		// Simulate check-in process
		err := checkinSystem.ProcessCheckin(reservation.ReservationID, counterID)
		if err != nil {
			fmt.Printf("Error processing check-in for reservation %s: %v\n", reservation.ReservationID, err)
		} else {
			fmt.Printf("Passenger %s (Reservation ID: %s) checked in successfully at Counter %d.\n", reservation.ClientID, reservation.ReservationID, counterID)
		}

		counterChannel <- counterID // Make counter available again
	}

	fmt.Printf("Check-in simulation for flight %s finished.\n", flightNumber)

	// You can add more logic here for waitlisted passengers, etc.
	// Example: Process waitlisted passengers after regular check-in
	waitlistedPassengers := checkinSystem.GetWaitlistedPassengers(flightNumber)
	if len(waitlistedPassengers) > 0 {
		fmt.Printf("\nProcessing waitlisted passengers for flight %s:\n", flightNumber)
		for _, reservation := range waitlistedPassengers {
			counterID := <-counterChannel // Wait for an available counter
			fmt.Printf("Processing waitlisted passenger %s (Reservation ID: %s) at Counter %d.\n", reservation.ClientID, reservation.ReservationID, counterID)

			err := checkinSystem.ProcessWaitlistCheckin(reservation.ReservationID, counterID)
			if err != nil {
				fmt.Printf("Error processing waitlist check-in for reservation %s: %v\n", reservation.ReservationID, err)
			} else {
				fmt.Printf("Waitlisted passenger %s (Reservation ID: %s) checked in successfully at Counter %d.\n", reservation.ClientID, reservation.ReservationID, counterID)
			}
			counterChannel <- counterID // Make counter available again
		}
	}
}