package paging_standard

import (
	"context"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"paging/internal/pkg/models"
	"paging/internal/pkg/repo"
	"paging/internal/pkg/tests"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	tests.TestWithPostgres(m)
}

func TestStandardPaging_Select(t *testing.T) {
	db, deferFunc := tests.ConnectToPostgres(t)
	defer deferFunc()

	usersRepo := repo.NewUsersRepo(db)
	standardPaging := NewStandardPaging(db)

	type args struct {
		ctx  context.Context
		page PageableRequest
	}
	tests := []struct {
		name        string
		prepareFunc func(t *testing.T, ctx context.Context)
		args        args
		want        Paged[models.User]
		wantErr     bool
	}{
		{
			name: "page(1, 2) on total 10",
			prepareFunc: func(t *testing.T, ctx context.Context) {
				for i := 1; i <= 10; i++ {
					_, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
					assert.NoError(t, err)
				}
			},
			args: args{
				page: PageableRequest{Number: 1, Size: 2},
			},
			want: Paged[models.User]{
				Content: []models.User{
					{Name: "User #" + strconv.Itoa(1)},
					{Name: "User #" + strconv.Itoa(2)},
				},
				Number:        1,
				Size:          2,
				TotalElements: 10,
				TotalPages:    5,
			},
		},
		{
			name: "page(2, 1) on total 3",
			prepareFunc: func(t *testing.T, ctx context.Context) {
				for i := 1; i <= 3; i++ {
					_, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
					assert.NoError(t, err)
				}
			},
			args: args{
				page: PageableRequest{Number: 2, Size: 1},
			},
			want: Paged[models.User]{
				Content: []models.User{
					{Name: "User #" + strconv.Itoa(2)},
				},
				Number:        2,
				Size:          1,
				TotalElements: 3,
				TotalPages:    3,
			},
		},
		{
			name: "page(1, 5) on total 2",
			prepareFunc: func(t *testing.T, ctx context.Context) {
				for i := 1; i <= 2; i++ {
					_, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
					assert.NoError(t, err)
				}
			},
			args: args{
				page: PageableRequest{Number: 1, Size: 5},
			},
			want: Paged[models.User]{
				Content: []models.User{
					{Name: "User #" + strconv.Itoa(1)},
					{Name: "User #" + strconv.Itoa(2)},
				},
				Number:        1,
				Size:          2,
				TotalElements: 2,
				TotalPages:    1,
			},
		},
		{
			name: "page(1, 1) on total 1",
			prepareFunc: func(t *testing.T, ctx context.Context) {
				for i := 1; i <= 1; i++ {
					_, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
					assert.NoError(t, err)
				}
			},
			args: args{
				page: PageableRequest{Number: 1, Size: 1},
			},
			want: Paged[models.User]{
				Content: []models.User{
					{Name: "User #" + strconv.Itoa(1)},
				},
				Number:        1,
				Size:          1,
				TotalElements: 1,
				TotalPages:    1,
			},
		},
		{
			name: "page(3, 2) on total 7",
			prepareFunc: func(t *testing.T, ctx context.Context) {
				for i := 1; i <= 7; i++ {
					_, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
					assert.NoError(t, err)
				}
			},
			args: args{
				page: PageableRequest{Number: 3, Size: 2},
			},
			want: Paged[models.User]{
				Content: []models.User{
					{Name: "User #" + strconv.Itoa(5)},
					{Name: "User #" + strconv.Itoa(6)},
				},
				Number:        3,
				Size:          2,
				TotalElements: 7,
				TotalPages:    4,
			},
		},
		{
			name: "page(4, 2) on total 7",
			prepareFunc: func(t *testing.T, ctx context.Context) {
				for i := 1; i <= 7; i++ {
					_, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
					assert.NoError(t, err)
				}
			},
			args: args{
				page: PageableRequest{Number: 4, Size: 2},
			},
			want: Paged[models.User]{
				Content: []models.User{
					{Name: "User #" + strconv.Itoa(7)},
				},
				Number:        4,
				Size:          1,
				TotalElements: 7,
				TotalPages:    4,
			},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.args.ctx = ctx

			assert.NoError(t, usersRepo.DeleteAll(ctx))

			if tt.prepareFunc != nil {
				tt.prepareFunc(t, ctx)
			}

			got, err := standardPaging.Select(tt.args.ctx, tt.args.page)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			lo.ForEach(tt.want.Content, func(item models.User, index int) {
				(&item).ID = -1 // does not matter
			})
			lo.ForEach(got.Content, func(item models.User, index int) {
				(&item).ID = -1 // does not matter
			})

			assert.Equal(t, tt.want.Number, got.Number)
			assert.Equal(t, tt.want.Size, got.Size)
			assert.Equal(t, tt.want.TotalElements, got.TotalElements)
			assert.Equal(t, tt.want.TotalPages, got.TotalPages)

			var wantContent, gotContent []*models.User
			for _, v := range tt.want.Content {
				v.ID = -1
				wantContent = append(wantContent, &v)
			}
			for _, v := range got.Content {
				v.ID = -1
				gotContent = append(gotContent, &v)
			}

			assert.EqualValues(t, wantContent, gotContent)
		})
	}
}

func TestStandardPaging_SelectAll(t *testing.T) {
	db, deferFunc := tests.ConnectToPostgres(t)
	defer deferFunc()

	usersRepo := repo.NewUsersRepo(db)
	standardPaging := NewStandardPaging(db)

	type args struct {
		ctx      context.Context
		pageSize int64
	}
	tests := []struct {
		name             string
		preparedDataSize int
		args             args
		pageSize         int64
	}{
		{
			name:             "page(1), all 7",
			preparedDataSize: 7,
			args: args{
				pageSize: 1,
			},
		},
		{
			name:             "page(3), all 10",
			preparedDataSize: 10,
			args: args{
				pageSize: 3,
			},
		},
		{
			name:             "page(10), all 100",
			preparedDataSize: 100,
			args: args{
				pageSize: 10,
			},
		},
		{
			name:             "page(100), all 1234",
			preparedDataSize: 1234,
			args: args{
				pageSize: 100,
			},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.args.ctx = ctx

			assert.NoError(t, usersRepo.DeleteAll(ctx))
			expectedData := make([]models.User, 0, tt.preparedDataSize)
			for i := 0; i < tt.preparedDataSize; i++ {
				u, err := usersRepo.Create(ctx, models.User{Name: "User #" + strconv.Itoa(i)})
				assert.NoError(t, err)

				expectedData = append(expectedData, u)
			}

			got, err := standardPaging.SelectAll(tt.args.ctx, tt.args.pageSize)
			assert.NoError(t, err)
			assert.EqualValues(t, expectedData, got)
		})
	}
}
