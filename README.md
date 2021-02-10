# Coffee Machine

## UPDATE

coffeeMachine Interface added in coffeeMachine.go

coffeeV1 and coffeeV2 are implementation of coffeeMachine.

In this particular code example,

coffeeV2 implements outlet to also tell time spent to acquire lock

Comment and Uncomment line 5 and 6 in main() in coffee.go to switch between coffeeV1 and coffeeV2

Example outputs:

When main() is run with coffeeV1,
<pre><code>
2021/02/11 03:05:27 Outlet 1 : Requesting hot_coffee
2021/02/11 03:05:27 Outlet 1 : After reserving - map[ginger_syrup:70 hot_milk:100 hot_water:400 sugar_syrup:50 tea_leaves_syrup:70]
2021/02/11 03:05:27 Outlet 2 : Requesting green_tea
2021/02/11 03:05:27 Outlet 2 : green_tea cannot be prepared because: 
                Outlet: 2 green_mixture is not available
2021/02/11 03:05:27 Outlet 2 : Requesting black_tea
2021/02/11 03:05:27 Outlet 2 : After reserving - map[ginger_syrup:40 hot_milk:100 hot_water:100 sugar_syrup:0 tea_leaves_syrup:40]
2021/02/11 03:05:27 Outlet 3 : Requesting hot_tea
2021/02/11 03:05:27 Outlet 3 : hot_tea cannot be prepared because: 
                Outlet: 3 sugar_syrup is not available
                Outlet: 3 hot_water is not sufficient
2021/02/11 03:05:29 Outlet 2 : black_tea PREPARED
2021/02/11 03:05:29 Outlet 1 : hot_coffee PREPARED
</code></pre>

When main() is run with coffeeV2,

<pre><code>
2021/02/11 03:06:47 Outlet 2 : Requesting hot_tea
<b><font size="4">2021/02/11 03:06:47 V2 Outlet 2 waited for 229ns to acquire lock</font></b>
2021/02/11 03:06:47 Outlet 2 : After reserving - map[ginger_syrup:90 hot_milk:400 hot_water:300 sugar_syrup:90 tea_leaves_syrup:70]
2021/02/11 03:06:47 Outlet 1 : Requesting hot_coffee
<b><font size="4">2021/02/11 03:06:47 V2 Outlet 1 waited for 183ns to acquire lock</font></b>
2021/02/11 03:06:47 Outlet 1 : After reserving - map[ginger_syrup:60 hot_milk:0 hot_water:200 sugar_syrup:40 tea_leaves_syrup:40]
2021/02/11 03:06:47 Outlet 3 : Requesting green_tea
<b><font size="4">2021/02/11 03:06:47 V2 Outlet 3 waited for 70ns to acquire lock</font></b>
2021/02/11 03:06:47 Outlet 3 : green_tea cannot be prepared because: 
                Outlet: 3 sugar_syrup is not sufficient
                Outlet: 3 green_mixture is not available
2021/02/11 03:06:47 Outlet 3 : Requesting black_tea
<b><font size="4">2021/02/11 03:06:47 V2 Outlet 3 waited for 78ns to acquire lock</font></b>
2021/02/11 03:06:47 Outlet 3 : black_tea cannot be prepared because: 
                Outlet: 3 hot_water is not sufficient
                Outlet: 3 sugar_syrup is not sufficient
2021/02/11 03:06:49 Outlet 2 : hot_tea PREPARED
2021/02/11 03:06:49 Outlet 1 : hot_coffee PREPARED
</code></pre>

## Run

Pass the test case file as argument

```bash
go run coffee.go coffeeMachine.go coffeeV1.go coffeeV2.go input.txt
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


