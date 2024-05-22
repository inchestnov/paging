package paging_cursor

import (
	"context"
	"database/sql"
	"paging/internal/pkg/models"
	"strconv"
	"strings"
)

type CursorPaging struct {
	db *sql.DB
}

func NewCursorPaging(db *sql.DB) *CursorPaging {
	return &CursorPaging{
		db: db,
	}
}

func (p *CursorPaging) Select(ctx context.Context, page PageableRequest) (Paged[models.User], error) {
	var totalElements int64
	if err := p.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM users").Scan(&totalElements); err != nil {
		return emptyUserPage(), err
	}
	totalPages := (totalElements + page.Size - 1) / page.Size

	var query strings.Builder
	query.WriteString("SELECT id, name FROM users ")

	if page.Cursor != "" {
		_, err := strconv.ParseInt(page.Cursor, 10, 64)
		if err != nil {
			return emptyUserPage(), err
		}

		query.WriteString("WHERE id > " + page.Cursor + " ")
	}

	query.WriteString("ORDER BY id ASC LIMIT " + strconv.FormatInt(page.Size+1, 10))

	rows, err := p.db.QueryContext(ctx, query.String())
	if err != nil {
		return emptyUserPage(), err
	}

	content := make([]models.User, 0, page.Size+1)
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

	hasNextPage := int64(len(content)) == page.Size+1
	var nextPage string
	if hasNextPage {
		content = content[:len(content)-1]
		nextPage = strconv.FormatInt(content[len(content)-1].ID, 10)
	}

	return Paged[models.User]{
		Content:        content,
		CursorNextPage: nextPage,
		Size:           int64(len(content)),
		TotalPages:     totalPages,
		TotalElements:  totalElements,
	}, nil
}

func (p *CursorPaging) SelectAll(ctx context.Context, pageSize int64) ([]models.User, error) {
	var content []models.User

	var cursor string
	for pageNumber := 1; true; pageNumber++ {
		page, err := p.Select(ctx, PageableRequest{
			Size:   pageSize,
			Cursor: cursor,
		})
		if err != nil {
			return nil, err
		}

		if content == nil {
			content = make([]models.User, 0, page.TotalElements)
		}
		content = append(content, page.Content...)

		if page.CursorNextPage == "" {
			break
		}

		cursor = page.CursorNextPage
	}
	return content, nil
}

func emptyUserPage() Paged[models.User] {
	return EmptyPage[models.User]()
}
