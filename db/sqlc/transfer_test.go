package db

import (
	"context"
	"simple_bank/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Вспомогательная функция для создания случайного перевода (transfer)
func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomInt(-1000000, 1000000),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

// Тест создания нового перевода
func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

// Тест получения отдельного перевода
func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

// Тест получения множества переводов
func TestListTransfers(t *testing.T) {
	// Создаём случайный аккаунт
	account := createRandomAccount(t)

	// Создаём 10 переводов с участием этого аккаунта как отправителя или получателя
	for i := 0; i < 10; i++ {
		// Каждый раз случайно выбираем: аккаунт как From или To
		if i%2 == 0 {
			createRandomTransferWithAccounts(t, account.ID, createRandomAccount(t).ID)
		} else {
			createRandomTransferWithAccounts(t, createRandomAccount(t).ID, account.ID)
		}
	}

	// Устанавливаем фильтры на основании существующего аккаунта
	arg := ListTransfersParams{
		FromAccountID: account.ID,
		ToAccountID:   account.ID,
		Limit:         5,
		Offset:        0, // начинаем с первой записи
	}

	// Получаем переводы
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)

	// Проверяем, что вернулись переводы
	require.NotEmpty(t, transfers)
	require.LessOrEqual(t, len(transfers), 5) // может быть меньше 5, если записей меньше

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		// Дополнительно проверяем, что аккаунт участвует в переводе
		require.True(t, transfer.FromAccountID == account.ID || transfer.ToAccountID == account.ID)
	}
}

// Создаёт перевод с заданными аккаунтами
func createRandomTransferWithAccounts(t *testing.T, fromID, toID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	return transfer
}
