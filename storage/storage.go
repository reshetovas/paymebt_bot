package storage

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"payment_bot/models"
	"time"
)

type SQLDBPaymentStorage struct {
	db *sql.DB
}

func NewSQLDBPaymentStorage(db *sql.DB) *SQLDBPaymentStorage {
	return &SQLDBPaymentStorage{db: db}
}

type PaymentStorage interface {
	SavePayment(p *models.Payment) error
	GetPaymentsByPeriod(userID int64, from, to time.Time) ([]*models.Payment, error)
	GetCountByCategory(userID int64, from, to time.Time) ([]models.CategoryReport, error)
}

func (r *SQLDBPaymentStorage) SavePayment(p *models.Payment) error {
	_, err := r.db.Exec(`INSERT INTO payments(user_id, amount, category, date) VALUES (?, ?, ?, ?)`,
		p.UserID, p.Amount, p.Category, p.Date)
	return err
}

func (r *SQLDBPaymentStorage) GetPaymentsByPeriod(userID int64, from, to time.Time) ([]*models.Payment, error) {
	rows, err := r.db.Query(`SELECT id, user_id, amount, category, date FROM payments WHERE user_id = ? AND date >= ? AND date <= ? ORDER BY date ASC`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		p := models.Payment{}
		var dateStr string
		if err := rows.Scan(&p.ID, &p.UserID, &p.Amount, &p.Category, &dateStr); err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse("2006-01-02 15:04:05-07:00", dateStr)
		if err != nil {
			return nil, err
		}
		p.Date = parsedDate

		payments = append(payments, &p)
	}
	return payments, nil
}

func (r *SQLDBPaymentStorage) GetCountByCategory(userID int64, from, to time.Time) ([]models.CategoryReport, error) {
	query := `SELECT CAST(category AS TEXT) as category, IFNULL(SUM(amount), 0) as amount
	FROM payments
	WHERE user_id = ? AND date >= ? AND date <= ?
	GROUP BY category
	
	UNION ALL
	
	SELECT 'Всего' as category, IFNULL(SUM(amount), 0) as amount
	FROM payments
	WHERE user_id = ? AND date >= ? AND date <= ?`

	rows, err := r.db.Query(query, userID, from, to, userID, from, to)
	if err != nil {
		log.Error().Err(err).Msg("Error executing query")
		return nil, err
	}
	defer rows.Close()

	result := []models.CategoryReport{}
	for rows.Next() {
		var r models.CategoryReport
		if err := rows.Scan(&r.Category, &r.Amount); err != nil {
			log.Error().Err(err).Msg("Error scanning row")
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}
