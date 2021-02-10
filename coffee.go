package main

func main() {
	// uncomment & comment to run with coffeeV1 or coffeeV2
	//c1 := new(coffeeV1)
	c1 := new(coffeeV2)
	c1.InitialiseMachine()

	// channel for placing orders
	bvg := make(chan string)

	//channel for refill instructions
	refill := make(chan map[string]int)

	var cm coffeeMachine
	cm = c1

	// Launch desired number of Go routines for outlet
	for i := 1; i <= c1.Outlets; i++ {
		go cm.outlet(i, bvg)
	}

	//Launch inlet Go Routine
	go cm.inlet(refill)

	c1.ExecuteCommands(bvg, refill)

	close(bvg)
	close(refill)
	c1.wg.Wait()

}
