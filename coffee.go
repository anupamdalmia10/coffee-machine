package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Mutex for synchronisation for accessing the ingredients
var mutex = &sync.Mutex{}

// Wait group to wait till all instructions are processed
var wg sync.WaitGroup

// Initial configuration of coffee machine
var Outlets int
var Ingredients = make(map[string]int)
var Capacity = make(map[string]int)
var Beverages = make(map[string]map[string]int)

// Variable to hold instructions passed by user
var Instruction []string

// Function to check if Ingredients for a particular beverage are present in terms of quantity
// Returns: 1) Approriate error message if Ingredient is less or not present
//			2) Empty string if ingredient is sufficient
func checkIngredients(beverage string, id int) (bool, []string) {

	result := []string{}

	// Loop over the recipe of given beverage and check for availability
	for v := range Beverages[beverage] {
		if Ingredients[v] < Beverages[beverage][v] {
			errorMessage := "is not available"

			if Ingredients[v] > 0 {
				errorMessage = "is not sufficient"
			}
			s := v + " " + errorMessage

			result = append(result, s)

		}
	}

	// Case for less quantity or unavailability of an Ingredient
	// result will contain at least 1 error in case of  shortage or non availability
	if len(result) > 0 {
		return false, result
	}

	//Case for suffiecient quantity, result will be empty in this case
	return true, result

}

// Function to reserve the ingredients for a beverage by reducing their quantity
func reserveIngredients(beverage string, id int) {
	// Loop over the recipe of given beverage to reduce its quantity
	for v := range Beverages[beverage] {
		Ingredients[v] -= Beverages[beverage][v]
	}
}

// Go routine to simulate Inlet(refill area)
func inlet(ing <-chan map[string]int) {
	defer wg.Done() // Tell the WorkGroup that lifespan of this inlet is over, right after exiting from this block

	// Keep waiting for data on ing channel
	// Terminates when channel is closed in main go Routine
	for i := range ing {
		//simulate processing delay for inlet
		time.Sleep(time.Second)
		for v := range i {
			// Check slot for given ingredient exists
			if _, present := Ingredients[v]; present {
				if i[v] < 1 {
					log.Println("Inlet:", "Taking ingredient out of me not allowed :P")
					break
				}
				// Acquire exclusive access on Ingredients to avoid updating while an outlet is reading Ingredients
				mutex.Lock()
				result := ""
				// Check if refilling with given quantity is within capacity
				// If it exceeds the capacity, refill upto max capacity
				if i[v]+Ingredients[v] <= Capacity[v] {
					Ingredients[v] = i[v]
					result = "refilled"
				} else if Ingredients[v] < Capacity[v] {
					Ingredients[v] = Capacity[v]
					result = "refilled full to capacity discarded the leftover"
				} else {
					result = "already full"
				}
				mutex.Unlock()

				log.Println("Inlet:", v, result)

			} else {
				log.Println("Inlet: Slot for", v, "does not exist")
			}
		}
	}
}

// Go Routing to simulate an outlet
// Receives an outlet id/number and a bvg channel, which has orders for various beverages
func outlet(id int, bvg <-chan string) {
	defer wg.Done() // Tell the WorkGroup that lifespan of this inlet is over, right after exiting from this block

	// Keep waiting for data on bvg channel
	// Terminates when channel is closed in main go Routine
	for beverage := range bvg {

		log.Println("Outlet", id, ":", "Requesting", beverage)

		if _, pres := Beverages[beverage]; pres {

			// Acquire exclusive access over Ingredients to avoid Reading simultaneously with other outlets
			// This will ensure that only 1 outlet is able to reserve the ingredients if they are sufficiently
			// available for beverage ordered on it
			mutex.Lock()

			ingredientsAvailable, result := checkIngredients(beverage, id)

			// Reserve the ingredients only if they are sufficiently available for the given beverage
			// Here reserve means that outlet has taken out the ingredient from inventory
			// and is mixing them in its mixing unit
			if ingredientsAvailable {
				reserveIngredients(beverage, id)
				log.Println("Outlet", id, ": After reserving -", Ingredients)
			} else {
				log.Println("Outlet", id, ":", beverage, "cannot be prepared because: ")
				// Display various insufficient or unavailable ingredients for the beverage
				for _, val := range result {
					fmt.Println("\t\tOutlet:", id, val)
				}
			}
			//Display the state of Ingredients

			//Release the exclusive access to allow other outlets to access Ingredients
			mutex.Unlock()

			if ingredientsAvailable {
				time.Sleep(time.Second * 2)
				log.Println("Outlet", id, ":", beverage, "PREPARED")
			}

		} else {
			log.Println("Outlet", id, ":", beverage, "DOES NOT EXIST")
		}

	}

}

func main() {
	///////////////////// Machine Initialisation Section /////////////////////

	//Reading from input file and parsing in to various data structures
	if len(os.Args) <= 1 {
		log.Fatal("No file specified")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())

		if len(words) >= 2 {
			switch words[0] {
			case "ingredients":
				Ingredients[words[1]], _ = strconv.Atoi(words[2])
				Capacity[words[1]], _ = strconv.Atoi(words[2])
			case "beverages":
				if _, pres := Beverages[words[1]]; pres {
					Beverages[words[1]][words[2]], _ = strconv.Atoi(words[3])
				} else {
					Beverages[words[1]] = make(map[string]int)
					Beverages[words[1]][words[2]], _ = strconv.Atoi(words[3])

				}
			case "outlets":
				Outlets, _ = strconv.Atoi(words[1])
			case "instruction":
				Instruction = append(Instruction, scanner.Text())
			case "//":
				continue
			}
		}

	}
	///////////////////// Machine Initialisation Section End /////////////////////

	if Outlets <= 0 {
		log.Fatal("No Outlet!!!, exiting")
	}

	//Add number of members to wait group
	// +1 for the inlet
	wg.Add(Outlets + 1)
	// channel for placing orders
	bvg := make(chan string)

	//channel for refill instructions
	refill := make(chan map[string]int)

	// Launch desired number of Go routines for outlet
	for i := 1; i <= Outlets; i++ {
		go outlet(i, bvg)
	}

	//Launch inlet Go Routine
	go inlet(refill)

	///////////////////// Executing given commands /////////////////////
	for _, v := range Instruction {
		words := strings.Fields(v)

		switch words[1] {
		case "order":
			bvg <- words[2]
		case "refill":
			amount, _ := strconv.Atoi(words[3])
			refill <- map[string]int{words[2]: amount}
		case "delay":
			amount, _ := strconv.Atoi(words[2])
			time.Sleep(time.Second * time.Duration(amount))
		}
	}
	/////////////////////  //////////////////////  /////////////////////

	close(bvg)
	close(refill)
	wg.Wait()

}
