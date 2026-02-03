package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Adapter is a PostgreSQL adapter for Casbin using pgx v5.
// It implements the persist.Adapter interface.
type Adapter struct {
	pool      *pgxpool.Pool
	tableName string
}

// CasbinRule represents a single policy rule in the database.
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// NewAdapter creates a new PostgreSQL adapter for Casbin.
func NewAdapter(pool *pgxpool.Pool) *Adapter {
	return &Adapter{
		pool:      pool,
		tableName: "shared.casbin_rule",
	}
}

// LoadPolicy loads all policy rules from the database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()

	query := fmt.Sprintf("SELECT ptype, v0, v1, v2, v3, v4, v5 FROM %s", a.tableName)
	rows, err := a.pool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query casbin rules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rule CasbinRule
		var v3, v4, v5 *string
		if err := rows.Scan(&rule.PType, &rule.V0, &rule.V1, &rule.V2, &v3, &v4, &v5); err != nil {
			return fmt.Errorf("failed to scan casbin rule: %w", err)
		}

		// Handle nullable columns
		if v3 != nil {
			rule.V3 = *v3
		}
		if v4 != nil {
			rule.V4 = *v4
		}
		if v5 != nil {
			rule.V5 = *v5
		}

		loadPolicyLine(&rule, model)
	}

	return rows.Err()
}

// SavePolicy saves all policy rules to the database.
func (a *Adapter) SavePolicy(model model.Model) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback will fail if transaction is committed
	}()

	// Clear existing policies
	if _, err := tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s", a.tableName)); err != nil {
		return fmt.Errorf("failed to clear casbin rules: %w", err)
	}

	// Insert all policies
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			if err := a.savePolicyLine(ctx, tx, ptype, rule); err != nil {
				return err
			}
		}
	}

	// Insert all roles
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			if err := a.savePolicyLine(ctx, tx, ptype, rule); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

// AddPolicy adds a policy rule to the database.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	ctx := context.Background()
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := a.savePolicyLine(ctx, tx, ptype, rule); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RemovePolicy removes a policy rule from the database.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	ctx := context.Background()

	query := fmt.Sprintf("DELETE FROM %s WHERE ptype = $1", a.tableName)
	args := []interface{}{ptype}

	for i, v := range rule {
		if v != "" {
			query += fmt.Sprintf(" AND v%d = $%d", i, i+2)
			args = append(args, v)
		}
	}

	result, err := a.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("policy not found")
	}

	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the database.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	ctx := context.Background()

	query := fmt.Sprintf("DELETE FROM %s WHERE ptype = $1", a.tableName)
	args := []interface{}{ptype}

	for i, v := range fieldValues {
		if v != "" {
			query += fmt.Sprintf(" AND v%d = $%d", fieldIndex+i, i+2)
			args = append(args, v)
		}
	}

	_, err := a.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove filtered policy: %w", err)
	}

	return nil
}

// savePolicyLine saves a single policy rule to the database within a transaction.
func (a *Adapter) savePolicyLine(ctx context.Context, tx pgx.Tx, ptype string, rule []string) error {
	values := make([]string, 6)
	for i := 0; i < len(rule) && i < 6; i++ {
		values[i] = rule[i]
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		a.tableName,
	)

	_, err := tx.Exec(ctx, query, ptype, values[0], values[1], values[2], values[3], values[4], values[5])
	if err != nil {
		return fmt.Errorf("failed to insert policy: %w", err)
	}

	return nil
}

// loadPolicyLine loads a single policy rule into the model.
func loadPolicyLine(rule *CasbinRule, model model.Model) {
	var p = []string{rule.PType}

	if rule.V0 != "" {
		p = append(p, rule.V0)
	}
	if rule.V1 != "" {
		p = append(p, rule.V1)
	}
	if rule.V2 != "" {
		p = append(p, rule.V2)
	}
	if rule.V3 != "" {
		p = append(p, rule.V3)
	}
	if rule.V4 != "" {
		p = append(p, rule.V4)
	}
	if rule.V5 != "" {
		p = append(p, rule.V5)
	}

	key := p[0]
	sec := key[:1]

	// Add the rule to the model
	if _, ok := model[sec]; ok {
		if _, ok := model[sec][key]; ok {
			model[sec][key].Policy = append(model[sec][key].Policy, p[1:])
		}
	}
}

// Ensure Adapter implements persist.Adapter interface.
var _ persist.Adapter = (*Adapter)(nil)
