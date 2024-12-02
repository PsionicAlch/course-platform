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

func (ds *DatabaseSeeder) SeedDiscounts(wg *sync.WaitGroup) {
	defer wg.Done()

	ds.InfoLog.Printf("Seeding discounts...")

	discounts := []struct {
		Title       string
		Description string
		Amount      uint
		Uses        uint
	}{
		{"Winter Clearance Sale", "Get ready for the winter season with massive discounts on winter apparel.", 25, 500},
		{"Holiday Shopping Discount", "Enjoy 30% off all holiday gifts, from decorations to toys.", 30, 200},
		{"New Year Special", "Ring in the new year with a 20% discount on all items in our store.", 20, 300},
		{"Black Friday Deal", "Save up to 40% off on select items during our Black Friday Sale.", 40, 150},
		{"Spring Fling Promotion", "Get 15% off all floral arrangements this spring.", 15, 250},
		{"Summer Essentials", "Don't miss out on our summer essentials sale with up to 35% off.", 35, 400},
		{"Back to School Special", "Stock up on school supplies with a 10% discount for students.", 10, 600},
		{"VIP Member Discount", "A special 50% discount for our loyal members.", 50, 100},
		{"Flash Sale", "48-hour flash sale! Get 60% off select electronics.", 60, 50},
		{"Early Bird Discount", "Get 18% off when you shop before 10 AM.", 18, 350},
		{"Cyber Monday Discount", "Enjoy 30% off on all online orders this Cyber Monday.", 30, 250},
		{"Friends and Family Offer", "Share this code for 15% off for you and your loved ones.", 15, 500},
		{"Summer Blowout Sale", "Save 45% off store-wide during our Summer Blowout Sale.", 45, 350},
		{"End of Season Clearance", "Huge savings of up to 70% on last season's collection.", 70, 100},
		{"Free Shipping Day", "Enjoy free shipping on all orders over $50.", 1, 1},
		{"Mystery Discount", "Find out your mystery discount when you check out! It could be up to 50%.", 50, 75},
		{"Thanksgiving Special", "Celebrate Thanksgiving with 25% off all kitchen items.", 25, 200},
		{"Fall Fashion Discount", "Get 20% off new fall fashion styles this season.", 20, 400},
		{"Big Savings Weekend", "Take 40% off your entire order this weekend only.", 40, 300},
		{"Weekend Flash Discount", "Flash sale! 50% off all weekend long on select items.", 50, 150},
		{"Buy One Get One Free", "Buy one item and get another one free on select products.", 1, 200},
		{"2-Day Sale", "2-day sale with discounts of 30% on all home goods.", 30, 250},
		{"New Collection Discount", "10% off all items in the new collection.", 10, 450},
		{"Holiday Gift Guide", "Get 35% off select items in our holiday gift guide.", 35, 300},
		{"Flash 24-Hour Sale", "Don't miss our 24-hour sale offering 20% off select categories.", 20, 500},
		{"Buy More, Save More", "Save an extra 10% for every $100 spent.", 10, 450},
		{"Weekend Flash Savings", "Flash savings event with up to 60% off select styles.", 60, 250},
		{"Spring Renewal", "Refresh your wardrobe with 30% off spring fashion essentials.", 30, 400},
		{"New Year's Eve Discount", "Countdown to savings with 40% off everything.", 40, 100},
		{"Valentine's Day Discount", "Surprise your loved one with 15% off all gift sets.", 15, 350},
		{"Exclusive Loyalty Offer", "50% off your next purchase as a thank you for being a loyal customer.", 50, 200},
		{"4th of July Celebration", "Celebrate Independence Day with 25% off all patriotic items.", 25, 300},
		{"Flash 50% Off Deal", "Limited-time 50% off select products, act fast!", 50, 150},
		{"Early Access Discount", "Get early access to sales with a 20% discount.", 20, 500},
		{"End of Year Savings", "Enjoy huge savings of up to 60% as we close the year.", 60, 200},
		{"Buy More, Save More (Tiered)", "Get 10% off your order when you spend $50, 20% off when you spend $100, and 30% off when you spend $150.", 30, 350},
		{"Black Friday Early Access", "Get 20% off your Black Friday purchase before the crowd!", 20, 300},
		{"Free Gift with Purchase", "Get a free gift with every purchase of $100 or more.", 1, 250},
		{"Summer Special Discount", "Beat the heat with a 35% discount on all summer essentials.", 35, 400},
		{"Gift Card Bonus", "Get an extra 15% when you purchase a $50 gift card.", 15, 150},
		{"Winter Warmth Promotion", "Save 25% on all outerwear to keep warm this winter.", 25, 350},
		{"Monthly Flash Discount", "Save 40% on select products every month.", 40, 500},
		{"Year-End Clearance", "Up to 75% off select items to end the year with a bang!", 75, 100},
		{"Winter Savings Event", "50% off winter gear and accessories!", 50, 250},
		{"Sneak Peek Offer", "Get 15% off our upcoming launch by using this code.", 15, 200},
		{"Black Friday Sneak Peek", "Enjoy exclusive discounts of 25% before Black Friday starts!", 25, 300},
		{"Weekend Bonanza", "Shop our weekend sale and get up to 45% off store-wide.", 45, 200},
		{"Labor Day Sale", "Celebrate Labor Day with 30% off select work wear.", 30, 400},
		{"Fall Clearance Sale", "Enjoy 40% off all fall items during our clearance event.", 40, 250},
		{"Flash 60% Off Deal", "Grab the best deals with up to 60% off select categories.", 60, 150},
	}

	for _, discount := range discounts {
		if err := ds.Database.AddDiscount(discount.Title, discount.Description, uint64(discount.Amount), uint64(discount.Uses)); err != nil {
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

	fake := faker.New()

	ds := &DatabaseSeeder{
		Loggers:        loggers,
		Database:       db,
		Authentication: auth,
		Faker:          fake,
	}

	wg := new(sync.WaitGroup)

	wg.Add(3)

	go ds.SeedUsers(wg)
	go ds.SeedAdmins(wg)
	go ds.SeedDiscounts(wg)

	wg.Wait()
}
