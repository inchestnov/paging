package paging_standard

import (
	"context"
	"database/sql"
	"fmt"
	"paging/internal/pkg/models"
)

type StandardPaging struct {
	db *sql.DB
}

func NewStandardPaging(db *sql.DB) *StandardPaging {
	return &StandardPaging{
		db: db,
	}
}

func (p *StandardPaging) Select(ctx context.Context, page PageableRequest) (Paged[models.User], error) {
	var totalElements int64
	if err := p.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM users").Scan(&totalElements); err != nil {
		return emptyUserPage(), err
	}
	totalPages := (totalElements + page.Size - 1) / page.Size

	query := fmt.Sprintf("SELECT id, name FROM users ORDER BY id LIMIT %d OFFSET %d", page.Size, (page.Number-1)*page.Size)
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return emptyUserPage(), err
	}
	defer rows.Close()

	content := make([]models.User, 0, page.Size)
	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err = rows.Scan(&id, &name); err != nil {
			return emptyUserPage(), err
		}

		content = append(content, models.User{ID: id, Name: name})
	}

	return Paged[models.User]{
		Content:       content,
		Number:        page.Number,
		Size:          int64(len(content)),
		TotalPages:    totalPages,
		TotalElements: totalElements,
	}, nil
}

func (p *StandardPaging) SelectAll(ctx context.Context, pageSize int64) ([]models.User, error) {
	var content []models.User

	for pageNumber := 1; true; pageNumber++ {
		page, err := p.Select(ctx, PageableRequest{
			Number: int64(pageNumber),
			Size:   pageSize,
		})
		if err != nil {
			return nil, err
		}

		if content == nil {
			content = make([]models.User, 0, page.TotalElements)
		}
		content = append(content, page.Content...)

		if page.TotalElements == 0 || page.Number == page.TotalPages {
			break
		}
	}
	return content, nil
}

func emptyUserPage() Paged[models.User] {
	return EmptyPage[models.User]()
}
