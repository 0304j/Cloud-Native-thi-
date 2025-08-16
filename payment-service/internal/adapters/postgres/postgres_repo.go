package postgres

import (
	"context"
	"payment-service/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PaymentRepositoryConnection struct {
	Conn *pgx.Conn
}

func (r *PaymentRepositoryConnection) GetPayment(ctx context.Context, id string) (domain.Payment, error) {
	var p domain.Payment
	err := r.Conn.QueryRow(ctx, `
        SELECT id, provider, amount, currency, status, created_at, updated_at
        FROM payments WHERE id = $1
    `, id).Scan(
		&p.ID, &p.Provider, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return domain.Payment{}, err
	}
	return p, nil
}

func (r *PaymentRepositoryConnection) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	rows, err := r.Conn.Query(ctx, `
		SELECT id, provider, amount, currency, status, created_at, updated_at
		FROM payments
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.Provider, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (s *PaymentRepositoryConnection) CreatePayment(ctx context.Context, payment domain.Payment) error {
	_, err := s.Conn.Exec(ctx, `
        INSERT INTO payments (provider, amount, currency, status)
        VALUES ($1, $2, $3, $4)
    `, payment.Provider, payment.Amount, payment.Currency, payment.Status)
	return err
}
