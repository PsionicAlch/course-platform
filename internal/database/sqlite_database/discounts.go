package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
)

func (db *SQLiteDatabase) GetDiscountsPaginated(term string, active *bool, page, elements uint) ([]*models.DiscountModel, error) {
	query := `SELECT id, title, description, code, discount, uses, active, created_at, updated_at FROM discounts WHERE (LOWER(id) LIKE '%' || ? ||'%' OR LOWER(title) LIKE '%' || ? || '%' OR LOWER(description) LIKE '%' || ? || '%' OR LOWER(code) LIKE '%' || ? || '%')`

	args := []any{term, term, term, term}

	if active != nil {
		query += " AND active = ?"

		if *active {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}

	offset := (page - 1) * elements
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var discounts []*models.DiscountModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all discounts from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var discount models.DiscountModel
		var active uint

		if err := rows.Scan(&discount.ID, &discount.Title, &discount.Description, &discount.Code, &discount.Discount, &discount.Uses, &active, &discount.CreatedAt, &discount.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from discounts table: %s\n", err)
			return nil, err
		}

		if active == 1 {
			discount.Active = true
		}

		discounts = append(discounts, &discount)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all discounts from the database: %s\n", err)
		return nil, err
	}

	return discounts, nil
}

func (db *SQLiteDatabase) GetAllDiscounts() ([]*models.DiscountModel, error) {
	query := `SELECT id, title, description, code, discount, uses, active, created_at, updated_at FROM discounts;`

	var discounts []*models.DiscountModel

	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all discounts from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var discount models.DiscountModel
		var active uint

		if err := rows.Scan(&discount.ID, &discount.Title, &discount.Description, &discount.Code, &discount.Discount, &discount.Uses, &active, &discount.CreatedAt, &discount.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from discounts table: %s\n", err)
			return nil, err
		}

		if active == 1 {
			discount.Active = true
		}

		discounts = append(discounts, &discount)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all discounts from the database: %s\n", err)
		return nil, err
	}

	return discounts, nil
}

func (db *SQLiteDatabase) CountDiscounts() (uint, error) {
	query := `SELECT COUNT(id) FROM discounts;`

	var count uint

	row := db.connection.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number of discounts in the database: %s\n", err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) AddDiscount(title, description string, discount, uses uint64) (string, error) {
	query := `INSERT INTO discounts (id, title, description, code, discount, uses) VALUES (?, ?, ?, ?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new discount: %s\n", err)
		return "", err
	}

	code, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate code for new discount: %s\n", err)
		return "", err
	}

	result, err := db.connection.Exec(query, id, title, description, code, discount, uses)
	if err != nil {
		db.ErrorLog.Printf("Failed to add new discount \"%s\" to the database: %s\n", title, err)
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query the database for the number of rows affected after adding new discount: %s\n", err)
		return "", err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after adding new discount to the database")
		return "", database.ErrNoRowsAffected
	}

	return id, nil
}

func (db *SQLiteDatabase) GetDiscountByID(discountId string) (*models.DiscountModel, error) {
	query := `SELECT id, title, description, code, discount, uses, active, created_at, updated_at FROM discounts WHERE id = ?;`

	discount := new(models.DiscountModel)

	var active int

	row := db.connection.QueryRow(query, discountId)
	if err := row.Scan(&discount.ID, &discount.Title, &discount.Description, &discount.Code, &discount.Discount, &discount.Uses, &active, &discount.CreatedAt, &discount.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get discount \"%s\" from the database: %s\n", discountId, err)
		return nil, err
	}

	discount.Active = active == 1

	return discount, nil
}

func (db *SQLiteDatabase) GetDiscountByCode(discountCode string) (*models.DiscountModel, error) {
	query := `SELECT id, title, description, code, discount, uses, active, created_at, updated_at FROM discounts WHERE code = ?;`

	discount := new(models.DiscountModel)

	var active int

	row := db.connection.QueryRow(query, discountCode)
	if err := row.Scan(&discount.ID, &discount.Title, &discount.Description, &discount.Code, &discount.Discount, &discount.Uses, &active, &discount.CreatedAt, &discount.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get discount by code \"%s\" from the database: %s\n", discountCode, err)
		return nil, err
	}

	discount.Active = active == 1

	return discount, nil
}

func (db *SQLiteDatabase) ActivateDiscount(discountId string) error {
	query := `UPDATE discounts SET active = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	if _, err := db.connection.Exec(query, discountId); err != nil {
		db.ErrorLog.Printf("Failed to update discount \"%s\" active status: %s\n", discountId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) DeactivateDiscount(discountId string) error {
	query := `UPDATE discounts SET active = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	if _, err := db.connection.Exec(query, discountId); err != nil {
		db.ErrorLog.Printf("Failed to update discount \"%s\" active status: %s\n", discountId, err)
		return err
	}

	return nil
}
