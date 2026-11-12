package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rlawnsxo131/go-webproject-base/internal/domain"
	"github.com/rlawnsxo131/go-webproject-base/internal/repository/entity"
)

var exampleStruct = sqlbuilder.NewStruct(new(entity.Example))

type ExampleRepository struct {
	db *sql.DB
}

func NewExampleRepository(db *sql.DB) *ExampleRepository {
	return &ExampleRepository{
		db: db,
	}
}

func (r *ExampleRepository) FindByID(ctx context.Context, exampleID int64, opts ...queryOptionFunc) (*domain.Example, error) {
	qOpts := parseOptions(opts)

	sb := exampleStruct.SelectFrom("example")
	sb.Where(sb.Equal("example_id", exampleID))

	if qOpts.WithLock {
		sb.ForUpdate()
	}

	var row *sql.Row
	query, args := sb.Build()

	debug(query, args)

	if qOpts.Tx != nil {
		row = qOpts.Tx.QueryRowContext(ctx, query, args...)
	} else {
		row = r.db.QueryRowContext(ctx, query, args...)
	}

	var ent entity.Example

	if err := row.Scan(exampleStruct.Addr(&ent)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("example not found: %w", domain.ErrExampleNotFound)
		}
		return nil, fmt.Errorf("scan example record: %w", err)
	}

	dom, err := ent.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("convert example entity to domain (example_id: %d): %w", ent.ExampleID, err)
	}

	return dom, nil
}

func (r *ExampleRepository) Insert(ctx context.Context, dom *domain.Example) (int64, error) {
	ent := entity.NewExampleEntityFromDomain(dom)

	ib := exampleStruct.InsertInto("example").
		Cols("name").
		Values(ent.Name)
	query, args := ib.Build()

	debug(query, args)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("execute insert example query: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id for example: %w", err)
	}

	return id, nil
}
