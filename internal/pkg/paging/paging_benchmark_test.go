package paging

import (
	"context"
	"fmt"
	"paging/internal/pkg/models"
	paging_cursor "paging/internal/pkg/paging/cursor"
	paging_standard "paging/internal/pkg/paging/standard"
	"paging/internal/pkg/repo"
	"paging/internal/pkg/tests"
	"testing"
)

func TestMain(m *testing.M) {
	tests.TestWithPostgres(m)
}

func BenchmarkPaging_SelectAll(b *testing.B) {
	db, deferFunc := tests.ConnectToPostgres(b)
	defer deferFunc()

	usersRepo := repo.NewUsersRepo(db)
	standardPaging := paging_standard.NewStandardPaging(db)
	cursorPaging := paging_cursor.NewCursorPaging(db)

	benchmarks := []struct {
		pageSize      int64
		totalElements int64
	}{
		{
			pageSize:      100,
			totalElements: 100,
		},
		{
			pageSize:      100,
			totalElements: 1000,
		},
		{
			pageSize:      100,
			totalElements: 10_000,
		},
		{
			pageSize:      500,
			totalElements: 100_000,
		},
		{
			pageSize:      1_000,
			totalElements: 500_000,
		},
		{
			pageSize:      1_000,
			totalElements: 1_000_000,
		},
	}
	for _, bb := range benchmarks {
		benchmarkPaging(b, "StandardPaging", usersRepo, bb.pageSize, bb.totalElements, standardPaging.SelectAll)
		benchmarkPaging(b, "CursorPaging", usersRepo, bb.pageSize, bb.totalElements, cursorPaging.SelectAll)
	}
}

func benchmarkPaging(
	b *testing.B,
	name string,
	usersRepo *repo.UsersRepo,
	pageSize int64,
	totalElements int64,
	loader func(ctx context.Context, pageSize int64) ([]models.User, error),
) {
	ctx := context.Background()
	if err := usersRepo.DeleteAll(ctx); err != nil {
		b.Fatal(err)
	}
	if err := usersRepo.CreateRandom(ctx, totalElements); err != nil {
		b.Fatal(err)
	}

	b.Run(fmt.Sprintf("[%v] pageSize = %d, totalElements = %d", name, pageSize, totalElements), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := loader(ctx, pageSize); err != nil {
				b.Error(err)
			}
		}
	})
}
