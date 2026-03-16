package db

import (
	"context"
	"simple_bank/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Функция для создания случайной записи (entry)
func createRandomEntry(t *testing.T) Entry {
	// Предварительное создание аккаунта
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID, // Используем ID существующего аккаунта
		Amount:    util.RandomInt(-1000000, 1000000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

// Тест создания новой записи
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

// Тест получения отдельной записи
func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// Тест получения множества записей
func TestListEntries(t *testing.T) {
	// Создадим десять записей с одинаковым аккаунтом
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntryWithAccount(t, account.ID)
	}

	// Запрашиваем последние 5 записей (со второго аккаунта)
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

// Функция для создания случайной записи с указанным идентификатором аккаунта
func createRandomEntryWithAccount(t *testing.T, accountID int64) Entry {
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomInt(-1000000, 1000000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}
