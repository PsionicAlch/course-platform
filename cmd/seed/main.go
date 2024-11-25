package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/jaswdr/faker/v2"
)

const Users = 100
const Admins = 100

type DatabaseSeeder struct {
	utils.Loggers
	Database       database.Database
	Authentication *authentication.Authentication
	Faker          faker.Faker
}

func (ds *DatabaseSeeder) SeedUsers(wg *sync.WaitGroup) {
	defer wg.Done()

	ds.InfoLog.Println("Seeding users...")

	for i := 0; i < Users; i++ {
		name := ds.Faker.Person().FirstName()
		surname := ds.Faker.Person().LastName()
		email := fmt.Sprintf("%s.%s@gmail.com", strings.ToLower(name), strings.ToLower(surname))
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

	for i := 0; i < Users; i++ {
		name := ds.Faker.Person().FirstName()
		surname := ds.Faker.Person().LastName()
		email := fmt.Sprintf("%s.%s@gmail.com", strings.ToLower(name), strings.ToLower(surname))
		password := "SuperSecurePassword123"

		if err := ds.Authentication.NewAdminUser(name, surname, email, password); err != nil {
			ds.ErrorLog.Fatalf("Failed to add new admin user to the database: %s\n", err)
		}
	}

	ds.InfoLog.Println("Finished seeding admins!")
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

	fake := faker.New()

	ds := &DatabaseSeeder{
		Loggers:        loggers,
		Database:       db,
		Authentication: auth,
		Faker:          fake,
	}

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go ds.SeedUsers(wg)
	go ds.SeedAdmins(wg)

	wg.Wait()
}
