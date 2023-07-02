package category

import (
	"database/sql"
	"sync"
	"time"
)

var (
	categoryCache      = map[string]*Category{}
	categoryCacheMutex sync.Mutex
	categoryLastSync   time.Time
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) List() ([]Category, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, color
		FROM categories;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var (
			category       Category
			rawDescription sql.NullString
			rawColor       sql.NullString
		)

		err = rows.Scan(&category.ID, &category.Name, &rawDescription, &rawColor)
		if err != nil {
			return nil, err
		}

		if rawDescription.Valid {
			category.Description = &rawDescription.String
		}
		if rawColor.Valid {
			category.Color = &rawColor.String
		}

		rules, err := s.GetRules(category.ID)
		if err != nil {
			return nil, err
		}
		category.Rules = rules

		categories = append(categories, category)
	}

	return categories, nil
}

func (s *Service) GetRules(categoryID int64) ([]CategoryRule, error) {
	rows, err := s.db.Query(`
		SELECT id, category_id, regex, mapping_field, description
		FROM category_rules
		WHERE category_id = $1;
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]CategoryRule, 0)
	for rows.Next() {
		var (
			rule           CategoryRule
			rawDescription sql.NullString
		)

		err = rows.Scan(&rule.ID, &rule.CategoryID, &rule.Regex, &rule.MappingField, &rawDescription)
		if err != nil {
			return nil, err
		}

		if rawDescription.Valid {
			rule.Description = &rawDescription.String
		}

		rules = append(rules, rule)
	}

	return rules, nil
}
