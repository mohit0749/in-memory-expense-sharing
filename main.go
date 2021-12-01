package main

import (
	"fmt"
	"testsplitwise/splitwise"
)

func main() {
	app := splitwise.NewSplitwise()
	u1 := app.CreateUser("mohit0", "mohit0@test.com", "9999929191")
	u2 := app.CreateUser("mohit1", "mohit1@test.com", "9999929192")
	u3 := app.CreateUser("mohit2", "mohit2@test.com", "9999929193")

	err := app.CreateExpense(u1, []int64{u1, u2, u3}, []float64{90.0}, "equal")
	if err != nil {
		fmt.Println(err)
	}
	app.ShowExpense(u1)
}
