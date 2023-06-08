package transaction

import (
	"database/sql"
	"fmt"
	"time"

	"docqube.de/bookkeeper/pkg/database"
	"docqube.de/bookkeeper/pkg/services/category"
)

var (
	ErrTransactionExists = fmt.Errorf("transaction already exists in database")
)

type Service struct {
	categoryService *category.Service
	categories      []category.Category
}

func NewService() *Service {
	return &Service{
		categoryService: category.NewService(),
		categories:      []category.Category{},
	}
}

func (s *Service) CategorizeAndImport(transactions []Transaction) error {
	categories, err := s.categoryService.List()
	if err != nil {
		return err
	}
	s.categories = categories

	for _, t := range transactions {
		category, err := s.MatchTransactionCategory(&t)
		if err != nil {
			return err
		}
		t.Category = category

		_, err = s.Create(t)
		if err != nil {
			if err == ErrTransactionExists {
				continue
			}
			return err
		}
	}

	return nil
}

func (s *Service) MatchTransactionCategory(transaction *Transaction) (*category.Category, error) {
	for _, c := range s.categories {
		matches, err := transaction.MatchesCategory(&c)
		if err != nil {
			return nil, err
		}

		if matches {
			return &c, nil
		}
	}
	return nil, nil
}

func (s *Service) Create(transaction Transaction) (*Transaction, error) {
	exists, err := s.Exists(transaction)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTransactionExists
	}

	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	var (
		id         int64
		categoryID *int64
	)

	if transaction.Category != nil {
		categoryID = &transaction.Category.ID
	}

	hash, err := transaction.Hash()
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`
		INSERT INTO transactions (
			booking_date,
			valuta_date,
			recipient,
			booking_text,
			purpose,
			balance,
			amount,
			category_id,
			hash
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		) RETURNING id;
	`,
		transaction.BookingDate,
		transaction.ValutaDate,
		transaction.Recipient,
		transaction.BookingText,
		transaction.Purpose,
		transaction.Balance,
		transaction.Amount,
		categoryID,
		hash,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	transaction.ID = id
	return &transaction, nil
}

func (s *Service) Get(id int64) (*Transaction, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	var (
		transaction         Transaction
		recipient           sql.NullString
		purpose             sql.NullString
		categoryID          sql.NullInt64
		categoryName        sql.NullString
		categoryDescription sql.NullString
		categoryColor       sql.NullString
	)
	err = db.QueryRow(`
		SELECT
			t.id,
			t.booking_date,
			t.valuta_date,
			t.recipient,
			t.booking_text,
			t.purpose,
			t.balance,
			t.amount,
			c.id,
			c.name,
			c.description,
			c.color
		FROM transactions AS t
			LEFT JOIN categories AS c ON t.category_id = c.id
		WHERE
			t.id = $1;
	`, id).Scan(
		&transaction.ID,
		&transaction.BookingDate,
		&transaction.ValutaDate,
		&recipient,
		&transaction.BookingText,
		&purpose,
		&transaction.Balance,
		&transaction.Amount,
		&categoryID,
		&categoryName,
		&categoryDescription,
		&categoryColor,
	)
	if err != nil {
		return nil, err
	}

	if recipient.Valid {
		transaction.Recipient = &recipient.String
	}
	if purpose.Valid {
		transaction.Purpose = &purpose.String
	}

	if categoryID.Valid {
		category := category.Category{
			ID:    categoryID.Int64,
			Name:  categoryName.String,
			Color: categoryColor.String,
		}
		if categoryDescription.Valid {
			category.Description = &categoryDescription.String
		}
		transaction.Category = &category
	}

	return &transaction, nil
}

func (s *Service) List(from, to time.Time) ([]Transaction, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			t.id,
			t.booking_date,
			t.valuta_date,
			t.recipient,
			t.booking_text,
			t.purpose,
			t.balance,
			t.amount,
			c.id,
			c.name,
			c.description,
			c.color
		FROM transactions AS t
			LEFT JOIN categories AS c ON t.category_id = c.id
		WHERE
			booking_date BETWEEN $1 AND $2
		ORDER BY booking_date DESC;
	`, database.NormalizeTime(from), database.NormalizeTime(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var (
			transaction         Transaction
			recipient           sql.NullString
			purpose             sql.NullString
			categoryID          sql.NullInt64
			categoryName        sql.NullString
			categoryDescription sql.NullString
			categoryColor       sql.NullString
		)
		err = rows.Scan(
			&transaction.ID,
			&transaction.BookingDate,
			&transaction.ValutaDate,
			&recipient,
			&transaction.BookingText,
			&purpose,
			&transaction.Balance,
			&transaction.Amount,
			&categoryID,
			&categoryName,
			&categoryDescription,
			&categoryColor,
		)
		if err != nil {
			return nil, err
		}

		if recipient.Valid {
			transaction.Recipient = &recipient.String
		}
		if purpose.Valid {
			transaction.Purpose = &purpose.String
		}

		if categoryID.Valid {
			category := category.Category{
				ID:    categoryID.Int64,
				Name:  categoryName.String,
				Color: categoryColor.String,
			}
			if categoryDescription.Valid {
				category.Description = &categoryDescription.String
			}
			transaction.Category = &category
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Service) ListByCategoryID(from time.Time, to time.Time, categoryID int64) ([]Transaction, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			t.id,
			t.booking_date,
			t.valuta_date,
			t.recipient,
			t.booking_text,
			t.purpose,
			t.balance,
			t.amount,
			c.id,
			c.name,
			c.description,
			c.color
		FROM transactions AS t
			INNER JOIN categories AS c ON t.category_id = c.id
		WHERE
			booking_date BETWEEN $1 AND $2
			AND category_id = $3
		ORDER BY booking_date DESC;
	`, database.NormalizeTime(from), database.NormalizeTime(to), categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var (
			transaction         Transaction
			recipient           sql.NullString
			purpose             sql.NullString
			category            category.Category
			categoryDescription sql.NullString
		)
		err = rows.Scan(
			&transaction.ID,
			&transaction.BookingDate,
			&transaction.ValutaDate,
			&recipient,
			&transaction.BookingText,
			&purpose,
			&transaction.Balance,
			&transaction.Amount,
			&category.ID,
			&category.Name,
			&categoryDescription,
			&category.Color,
		)
		if err != nil {
			return nil, err
		}

		if recipient.Valid {
			transaction.Recipient = &recipient.String
		}
		if purpose.Valid {
			transaction.Purpose = &purpose.String
		}

		if categoryDescription.Valid {
			category.Description = &categoryDescription.String
		}

		transaction.Category = &category
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Service) Categorize(id, categoryID int64) error {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE transactions
		SET category_id = $1
		WHERE id = $2;
	`, categoryID, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ListUnclassified(from time.Time, to time.Time) ([]Transaction, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			id,
			booking_date,
			valuta_date,
			recipient,
			booking_text,
			purpose,
			balance,
			amount
		FROM transactions
		WHERE
			category_id IS NULL
			AND
			booking_date BETWEEN $1 AND $2
		ORDER BY booking_date DESC;
	`, database.NormalizeTime(from), database.NormalizeTime(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var (
			transaction Transaction
			recipient   sql.NullString
			purpose     sql.NullString
		)
		err = rows.Scan(
			&transaction.ID,
			&transaction.BookingDate,
			&transaction.ValutaDate,
			&recipient,
			&transaction.BookingText,
			&purpose,
			&transaction.Balance,
			&transaction.Amount,
		)
		if err != nil {
			return nil, err
		}

		if recipient.Valid {
			transaction.Recipient = &recipient.String
		}
		if purpose.Valid {
			transaction.Purpose = &purpose.String
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Service) Exists(transaction Transaction) (bool, error) {
	hash, err := transaction.Hash()
	if err != nil {
		return false, err
	}

	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}

	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM transactions
			WHERE hash = $1
		);
	`, hash).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
