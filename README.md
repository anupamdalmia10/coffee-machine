# Coffee Machine



## Run

Pass the test case file as argument

```bash
go run coffee.go input.txt
```

- To introduce a delay between instructions in test case files, add this to the point you want to introduce delay ( in seconds )

```bash
instruction delay 2
```

## Information

This uses Go routine to somewhat simulate the desired parallel behavior of the machine.

Go routine can execute in random order, so there can be different output for different runs.

Mutex is used to synchronize access to Ingredients.

Channels are used to communicate to Go routine about orders and refills.

## Assumptions

* As soon as the outlet has reserved the ingredients, it means it has taken out the required ingredients from inventory.
* Once the outlet has exclusive access to ingredients, checking and reserving the Ingredients happens instantly.
* The coffee making process delay kicks in after taking out the ingredients. Which means other go routines are free to access Ingredients once a Go routine is done with its critical section of checking and reserving ingredients.
* Time taken to prepare a beverage is 2 seconds.
* Time taken for inlet operation is 1 second.


