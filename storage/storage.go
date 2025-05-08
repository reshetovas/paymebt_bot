package storage

import (
	"database/sql"
	"payment_bot/models"
)

type SQLDBPaymentStorage struct {
	db *sql.DB
}

func NewSQLDBPaymentStorage(db *sql.DB) *SQLDBPaymentStorage {
	return &SQLDBPaymentStorage{db: db}
}

type PaymentStorage interface {
	SavePayment(p *models.Payment) error
	GetPaymentsByPeriod(userID int64, from, to string) ([]*models.Payment, error)
}

func (r *SQLDBPaymentStorage) SavePayment(p *models.Payment) error {
	_, err := r.db.Exec(`INSERT INTO payments(user_id, amount, category, date) VALUES (?, ?, ?, ?)`,
		p.UserID, p.Amount, p.Category, p.Date)
	return err
}

func (r *SQLDBPaymentStorage) GetPaymentsByPeriod(userID int64, from, to string) ([]*models.Payment, error) {
	rows, err := r.db.Query(`SELECT id, user_id, amount, category, date FROM payments WHERE user_id = ? AND date BETWEEN ? AND ?`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		p := new(models.Payment)
		if err := rows.Scan(&p.ID, &p.UserID, &p.Amount, &p.Category, &p.Date); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
