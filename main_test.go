package main_test

import (
	"context"
	"fmt"
	"time"

	"github.com/itsubaki/ent-training/ent"
	"github.com/itsubaki/ent-training/ent/car"
	"github.com/itsubaki/ent-training/ent/user"

	_ "github.com/mattn/go-sqlite3"
)

func Example() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		fmt.Printf("failed opening connection to sqlite: %v", err)
		return
	}
	defer client.Close()

	// Run the auto migration tool.
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		fmt.Printf("failed creating schema resources: %v", err)
		return
	}

	// create
	{
		u, err := client.User.
			Create().
			SetAge(30).
			SetName("a8m").
			Save(ctx)
		if err != nil {
			fmt.Printf("failed creating user: %v", err)
			return
		}

		fmt.Println("user was created: ", u)
	}

	// query
	{
		u, err := client.User.
			Query().
			Where(user.Name("a8m")).
			// `Only` fails if no user found,
			// or more than 1 user returned.
			Only(ctx)
		if err != nil {
			fmt.Printf("failed querying user: %v", err)
			return
		}

		fmt.Println("user returned: ", u)
	}

	// Output:
	// user was created:  User(id=1, age=30, name=a8m)
	// user returned:  User(id=1, age=30, name=a8m)
}

func Example_edge() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		fmt.Printf("failed opening connection to sqlite: %v", err)
		return
	}
	defer client.Close()

	// Run the auto migration tool.
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		fmt.Printf("failed creating schema resources: %v", err)
		return
	}

	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC)).
		Save(ctx)
	if err != nil {
		fmt.Printf("failed creating car: %v", err)
		return
	}
	fmt.Println("car was created: ", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC)).
		Save(ctx)
	if err != nil {
		fmt.Printf("failed creating car: %v", err)
		return
	}
	fmt.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		fmt.Printf("failed creating user: %v", err)
		return
	}
	fmt.Println("user was created: ", a8m)

	// query
	{
		cars, err := a8m.QueryCars().All(ctx)
		if err != nil {
			fmt.Printf("failed querying user cars: %v", err)
			return
		}
		fmt.Println("returned cars:", cars)

		// What about filtering specific cars.
		ford, err := a8m.QueryCars().
			Where(car.Model("Ford")).
			Only(ctx)
		if err != nil {
			fmt.Printf("failed querying user cars: %v", err)
			return
		}
		fmt.Println(ford)
	}

	// Output:
	// car was created:  Car(id=1, model=Tesla, registered_at=Sun Mar 12 00:00:00 2023)
	// car was created:  Car(id=2, model=Ford, registered_at=Sun Mar 12 00:00:00 2023)
	// user was created:  User(id=1, age=30, name=a8m)
	// returned cars: [Car(id=1, model=Tesla, registered_at=Sun Mar 12 00:00:00 2023) Car(id=2, model=Ford, registered_at=Sun Mar 12 00:00:00 2023)]
	// Car(id=2, model=Ford, registered_at=Sun Mar 12 00:00:00 2023)
}
