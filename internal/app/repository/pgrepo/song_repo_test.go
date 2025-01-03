package pgrepo

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"songs/internal/app/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *SongRepo) {
	// Create a new SQL mock
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}

	repo := NewSongRepo(db)
	return mockDB, mock, repo
}

func TestGetSong(t *testing.T) {
	mockDB, mock, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()
	now := time.Now()
	expectedSong := &domain.Song{
		ID:          1,
		GroupID:     1,
		Title:       "Test Song",
		ReleaseDate: now,
		Text:        "Test lyrics",
		Link:        "http://example.com",
	}

	rows := sqlmock.NewRows([]string{
		"id", "group_id", "title", "release_date", "text", "link",
	}).AddRow(
		expectedSong.ID,
		expectedSong.GroupID,
		expectedSong.Title,
		expectedSong.ReleaseDate,
		expectedSong.Text,
		expectedSong.Link,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "songs" WHERE "songs"."id" = $1 ORDER BY "songs"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	song, err := repo.GetSong(ctx, 1)

	if assert.NoError(t, err) {
		assert.NotNil(t, song)
		assert.Equal(t, expectedSong.ID, song.ID)
		assert.Equal(t, expectedSong.Title, song.Title)
		assert.Equal(t, expectedSong.GroupID, song.GroupID)
		assert.Equal(t, expectedSong.Text, song.Text)
		assert.Equal(t, expectedSong.Link, song.Link)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSong_NotFound(t *testing.T) {
	mockDB, mock, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "songs" WHERE "songs"."id" = $1 ORDER BY "songs"."id" LIMIT $2`)).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "group_id", "title", "release_date", "text", "link", "created_at", "updated_at",
		}))

	song, err := repo.GetSong(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, song)
	assert.Contains(t, err.Error(), "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSongs(t *testing.T) {
	mockDB, mock, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()
	filter := map[string]string{"title": "Test"}
	now := time.Now()

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "songs"`)).
		WillReturnRows(countRows)

	// Mock data query
	rows := sqlmock.NewRows([]string{"id", "group_id", "title", "release_date", "text", "link", "created_at", "updated_at"}).
		AddRow(1, 1, "Test Song 1", now, "Lyrics 1", "link1", now, now).
		AddRow(2, 1, "Test Song 2", now, "Lyrics 2", "link2", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "songs"`)).
		WillReturnRows(rows)

	songs, total, err := repo.GetSongs(ctx, filter, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, songs, 2)
	assert.Equal(t, "Test Song 1", songs[0].Title)
	assert.Equal(t, "Test Song 2", songs[1].Title)
}

func TestCreateSong(t *testing.T) {
	mockDB, mock, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()
	now := time.Now()
	newSong := &domain.Song{
		GroupID:     1,
		Title:       "New Song",
		ReleaseDate: now,
		Text:        "New lyrics",
		Link:        "http://example.com/new",
	}

	// Expect Begin transaction
	mock.ExpectBegin()

	// Expect the INSERT query with RETURNING clause
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "songs" ("group_id","title","release_date","text","link") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			newSong.GroupID,
			newSong.Title,
			newSong.ReleaseDate,
			newSong.Text,
			newSong.Link,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Expect Commit transaction
	mock.ExpectCommit()

	createdSong, err := repo.CreateSong(ctx, newSong)

	if assert.NoError(t, err) {
		assert.NotNil(t, createdSong)
		assert.Equal(t, newSong.Title, createdSong.Title)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSong_InvalidID(t *testing.T) {
	mockDB, _, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()

	song, err := repo.GetSong(ctx, 0)

	assert.Error(t, err)
	assert.Nil(t, song)
	assert.Contains(t, err.Error(), "invalid song ID")
}

func TestGetSongs_EmptyResult(t *testing.T) {
	mockDB, mock, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()
	filter := map[string]string{"title": "Nonexistent"}

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "songs"`)).
		WillReturnRows(countRows)

	// Mock empty data query
	rows := sqlmock.NewRows([]string{"id", "group_id", "title", "release_date", "text", "link", "created_at", "updated_at"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "songs"`)).
		WillReturnRows(rows)

	songs, total, err := repo.GetSongs(ctx, filter, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, songs)
}

func TestCreateSong_Error(t *testing.T) {
	mockDB, mock, repo := setupTest(t)
	defer func() {
		_ = mockDB.Close()
	}()

	ctx := context.Background()
	newSong := &domain.Song{
		GroupID:     1,
		Title:       "New Song",
		ReleaseDate: time.Now(),
		Text:        "New lyrics",
		Link:        "http://example.com/new",
	}

	// Expect Begin transaction
	mock.ExpectBegin()

	// Expect the INSERT query to fail
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "songs" ("group_id","title","release_date","text","link") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			newSong.GroupID,
			newSong.Title,
			newSong.ReleaseDate,
			newSong.Text,
			newSong.Link,
		).
		WillReturnError(sql.ErrConnDone)

	// Expect Rollback since the operation failed
	mock.ExpectRollback()

	createdSong, err := repo.CreateSong(ctx, newSong)

	assert.Error(t, err)
	assert.Nil(t, createdSong)
	assert.Contains(t, err.Error(), "failed to create song")
	assert.NoError(t, mock.ExpectationsWereMet())
}
