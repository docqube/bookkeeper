package category

import (
	"database/sql"
	"sync"
	"time"

	"docqube.de/bookkeeper/pkg/database"
)

var (
	categoryCache      = map[string]*Category{}
	categoryCacheMutex sync.Mutex
	categoryLastSync   time.Time
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) List() ([]Category, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
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
		)

		err = rows.Scan(&category.ID, &category.Name, &rawDescription, &category.Color)
		if err != nil {
			return nil, err
		}

		if rawDescription.Valid {
			category.Description = &rawDescription.String
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
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
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
