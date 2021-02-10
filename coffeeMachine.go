package main

type coffeeMachine interface {
	inlet(ing <-chan map[string]int)
	outlet(id int, bvg <-chan string)
}
