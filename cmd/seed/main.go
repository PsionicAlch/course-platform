package main

import (
	"fmt"
	"sync"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database"
	"github.com/PsionicAlch/course-platform/internal/utils"
)

type DatabaseSeeder struct {
	utils.Loggers
	Database       database.Database
	Authentication *authentication.Authentication
}

func (ds *DatabaseSeeder) SeedUsers(wg *sync.WaitGroup) {
	defer wg.Done()

	ds.InfoLog.Println("Seeding users...")

	for _, user := range Users {
		name := user.Name
		surname := user.Surname
		email := fmt.Sprintf("%s.%s@gmail.com", name, surname)
		password := "SuperSecurePassword123"

		if err := ds.Authentication.NewUser(name, surname, email, password); err != nil {
			ds.ErrorLog.Fatalf("Failed to add new user to the database: %s\n", err)
		}
	}

	ds.InfoLog.Println("Finished seeding users!")
}

func (ds *DatabaseSeeder) SeedAdmins(wg *sync.WaitGroup) {
	defer wg.Done()

	ds.InfoLog.Println("Seeding admins...")

	for _, admin := range Admins {
		name := admin.Name
		surname := admin.Surname
		email := fmt.Sprintf("%s.%s@gmail.com", name, surname)
		password := "SuperSecurePassword123"

		if err := ds.Authentication.NewAdminUser(name, surname, email, password); err != nil {
			ds.ErrorLog.Fatalf("Failed to add new admin user to the database: %s\n", err)
		}
	}

	ds.InfoLog.Println("Finished seeding admins!")
}

func (ds *DatabaseSeeder) SeedDiscounts(wg *sync.WaitGroup) {
	defer wg.Done()

	ds.InfoLog.Printf("Seeding discounts...")

	for _, discount := range Discounts {
		if _, err := ds.Database.AddDiscount(discount.Title, discount.Description, uint64(discount.Amount), uint64(discount.Uses)); err != nil {
			ds.ErrorLog.Fatalf("Failed to add new discount to the database: %s\n", err)
		}
	}

	ds.InfoLog.Println("Finished seeding discounts!")
}

func main() {
	loggers := utils.CreateLoggers("DATABASE SEEDER")

	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to open database connection: ", err)
	}
	defer db.Close()

	auth, err := authentication.SetupAuthentication(db, nil, 0, 0, "", "", "", "")
	if err != nil {
		loggers.ErrorLog.Fatalf("Failed to set up authentication: %s\n", err)
	}

	ds := &DatabaseSeeder{
		Loggers:        loggers,
		Database:       db,
		Authentication: auth,
	}

	wg := new(sync.WaitGroup)

	wg.Add(3)

	go ds.SeedUsers(wg)
	go ds.SeedAdmins(wg)
	go ds.SeedDiscounts(wg)

	wg.Wait()
}
