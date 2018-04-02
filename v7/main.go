// v7: use a processor to execute a series of functions
// https://play.golang.org/p/m-4Ah-Skb29

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
func Drive(s Shopper, fuelRequired int) (Shopper, error) {

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
func BuyEggs(s Shopper, howMany int) (Shopper, error) {

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

// Flavorize creates a single-arg function by closing over the integer argument.
func Flavorize(f func(Shopper, int) (Shopper, error), arg int) func(Shopper) (Shopper, error) {
	return func(s Shopper) (Shopper, error) {
		return f(s, arg)
	}
}

// ProcessSteps executes a series of functions. Stop processing on error.
func ProcessSteps(s Shopper, steps ...func(Shopper) (Shopper, error)) (Shopper, error) {

	for _, oneStep := range steps {
		var err error // avoid shadowing s
		s, err = oneStep(s)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

func main() {

	shopper := NewShopper(FuelInCar, DollarsInWallet)

	log.Printf("gonna try and buy eggs")
	log.Printf("shopper: %s", shopper)

	driveToStore := Flavorize(Drive, FuelNeededToGetToStore)
	buyEggs := Flavorize(BuyEggs, EggsRequired)
	driveHome := Flavorize(Drive, FuelNeededToGetHome)

	shopper, err := ProcessSteps(shopper,
		driveToStore,
		buyEggs,
		driveHome,
	)

	if err != nil {
		log.Fatalf("could not complete shopping: %s", err)
	}

	log.Printf("got the eggs!")
	log.Printf("shopper: %s", shopper)
}
