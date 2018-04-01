// v3: include err pre-check in each method
// https://play.golang.org/p/t60wdCWFDqB

package main

import (
	"fmt"
	"log"
)

// Configuration constants for the program. Done as constants so that
// program will run in go playground (no command line required).
const (
	FuelNeededToGetToStore = 5
	FuelNeededToGetHome    = 4
	EggsAvailableAtStore   = 24
	CostPerEgg             = 1
	FuelInCar              = 15
	DollarsInWallet        = 10
	EggsRequired           = 6
)

// Shopper is a resource to go buy eggs.
type Shopper struct {
	Fuel    int
	Dollars int
	Eggs    int
}

// String method produces a string from a Shopper.
func (s Shopper) String() string {
	return fmt.Sprintf("fuel=%d, dollars=%d, eggs=%d", s.Fuel, s.Dollars, s.Eggs)
}

// NewShopper constructs a new Shopper with specified fuel, dollars, but 0 eggs.
func NewShopper(fuelInCar, dollarsInWallet int) Shopper {
	return Shopper{
		Fuel:    fuelInCar,
		Dollars: dollarsInWallet,
		Eggs:    0,
	}
}

// Drive reduces fuel by distance. If there's not enough fuel, we'll run out of gas.
func (s Shopper) Drive(fuelRequired int, err error) (Shopper, error) {
	if err != nil {
		return s, err
	}

	s.Fuel -= fuelRequired
	if s.Fuel < 0 {
		err := fmt.Errorf("ran out of gas %d from destination", s.Fuel*-1)
		s.Fuel = 0
		return s, err
	}

	return s, nil
}

// BuyEggs checks if there are enough eggs, and if the shopper can afford them.
// Add eggs, reduce dollars if possible to buy.
func (s Shopper) BuyEggs(howMany int, err error) (Shopper, error) {
	if err != nil {
		return s, err
	}

	if howMany > EggsAvailableAtStore {
		return s, fmt.Errorf("there are only %d eggs available, but we need %d", EggsAvailableAtStore, howMany)
	}

	totalCost := howMany * CostPerEgg

	if totalCost > s.Dollars {
		return s, fmt.Errorf("only have %d dollars, but need %d to buy %d eggs", s.Dollars, totalCost, howMany)
	}

	s.Eggs += howMany
	s.Dollars -= totalCost

	return s, nil
}

func main() {

	shopper := NewShopper(FuelInCar, DollarsInWallet)

	log.Printf("gonna try and buy eggs")
	log.Printf("shopper: %s", shopper)

	shopper, err := shopper.Drive(FuelNeededToGetToStore, nil)
	shopper, err = shopper.BuyEggs(EggsRequired, err)
	shopper, err = shopper.Drive(FuelNeededToGetHome, err)

	if err != nil {
		log.Fatalf("could not complete shopping: %s", err)
	}

	log.Printf("got the eggs!")
	log.Printf("shopper: %s", shopper)
}
