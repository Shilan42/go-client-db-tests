package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// Тест проверяет корректность работы функции selectClient при успешном выполнении
func Test_SelectClient_WhenOk(t *testing.T) {
	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err, "database connection error: %v", err)
	// Закрытие соединения после завершения теста
	defer db.Close()

	// ID клиента для тестирования
	clientID := 1

	// Получение данных клиента из базы
	client, err := selectClient(db, clientID)
	// Проверка, что при получении данных клиента из БД не было ошибок
	require.NoError(t, err, "error retrieving client with ID %d: %v", clientID, err)

	// Подтест для проверки полей клиента на корректность ID и заполненность всех строковых полей
	t.Run("CheckClientFields", func(t *testing.T) {
		// Проверка совпадения ID
		assert.Equal(t, client.ID, clientID, "ID mismatch: expected %d, got %d", clientID, client.ID)
		// Проверка обязательных полей
		assert.NotEmpty(t, client.Birthday, "birthday field should not be empty for client ID %d", clientID)
		assert.NotEmpty(t, client.Email, "email field should not be empty for client ID %d", clientID)
		assert.NotEmpty(t, client.FIO, "FIO field should not be empty for client ID %d", clientID)
		assert.NotEmpty(t, client.Login, "login field should not be empty for client ID %d", clientID)
	})
}

// Тест проверяет корректность обработки кейсов, когда клиент с указанным ID отсутствует в БД
func Test_SelectClient_WhenNoClient(t *testing.T) {
	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err, "database connection error: %v", err)
	// Закрытие соединения после завершения теста
	defer db.Close()

	// Невалидный ID клиента для тестирования (несуществующий в базе)
	clientID := -1

	// Попытка получения данных несуществующего клиента
	client, err := selectClient(db, clientID)
	// Проверка возникновения ошибки и проверка типа ошибки
	require.Error(t, err, "expected error when selecting non-existent client with ID %d", clientID)
	require.Equal(t, sql.ErrNoRows, err, "expected sql.ErrNoRows error when selecting client with ID %d", clientID)

	// Подтест для проверки состояния объекта клиента на отсутствии данных в БД
	t.Run("CheckClientFields", func(t *testing.T) {
		// Проверка, что все поля пустые
		assert.Empty(t, client.ID, "ID field should be empty for non-existent client with ID %d", clientID)
		assert.Empty(t, client.Birthday, "birthday field should be empty for non-existent client with ID %d", clientID)
		assert.Empty(t, client.Email, "email field should be empty for non-existent client with ID %d", clientID)
		assert.Empty(t, client.FIO, "FIO field should be empty for non-existent client with ID %d", clientID)
		assert.Empty(t, client.Login, "login field should be empty for non-existent client with ID %d", clientID)
	})
}

// Тест проверяет корректность вставки нового клиента в базу данных
func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err, "database connection error: %v", err)
	// Закрытие соединения после завершения теста
	defer db.Close()

	// Создание тестового объекта клиента с тестовыми данными
	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}
	// Вставка нового клиента в базу данных
	cl.ID, err = insertClient(db, cl)
	// Проверка, что у клиента появилось ID и не было ошибок при вставке
	assert.NotEmpty(t, cl.ID, "ID should not be empty after client insertion: %v", cl)
	require.NoError(t, err, "error inserting client: %v, error: %v", cl, err)

	// Получение вставленного клиента из базы
	client, err := selectClient(db, cl.ID)
	require.NoError(t, err, "error retrieving client with ID %d: %v", cl.ID, err)

	// Проверка соответствия полученных данных исходным
	assert.Equal(t, client.ID, cl.ID, "ID mismatch: expected %v, actual %v", cl.ID, client.ID)
	assert.Equal(t, client.FIO, cl.FIO, "FIO mismatch: expected %v, actual %v", cl.FIO, client.FIO)
	assert.Equal(t, client.Login, cl.Login, "Login mismatch: expected %v, actual %v", cl.Login, client.Login)
	assert.Equal(t, client.Birthday, cl.Birthday, "birthday mismatch: expected %v, actual %v", cl.Birthday, client.Birthday)
	assert.Equal(t, client.Email, cl.Email, "email mismatch: expected %v, actual %v", cl.Email, client.Email)

	// Очистка тестовых данных
	err = deleteClient(db, cl.ID)
	require.NoError(t, err, "Error deleting client with ID %d: %v", cl.ID, err)
}

// Тест проверяет корректность удаления нового клиента из БД
func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err, "database connection error: %v", err)
	// Закрытие соединения после завершения теста
	defer db.Close()

	// Создание тестового объекта клиента с тестовыми данными
	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// Вставка нового клиента в базу данных
	cl.ID, err = insertClient(db, cl)
	// Проверка, что у клиента появилось ID и не было ошибок при вставке
	require.NotEmpty(t, cl.ID, "ID should not be empty after client insertion: %v", cl)
	require.NoError(t, err, "error inserting client: %v, error: %v", cl, err)

	// Получение вставленного клиента из базы
	client, err := selectClient(db, cl.ID)
	require.NoError(t, err, "error retrieving client with ID %d: %v", cl.ID, err)

	// Проверка соответствия полученных данных исходным
	assert.Equal(t, client.ID, cl.ID, "ID mismatch: expected %v, actual %v", cl.ID, client.ID)
	assert.Equal(t, client.FIO, cl.FIO, "FIO mismatch: expected %v, actual %v", cl.FIO, client.FIO)
	assert.Equal(t, client.Login, cl.Login, "login mismatch: expected %v, actual %v", cl.Login, client.Login)
	assert.Equal(t, client.Birthday, cl.Birthday, "birthday mismatch: expected %v, actual %v", cl.Birthday, client.Birthday)
	assert.Equal(t, client.Email, cl.Email, "email mismatch: expected %v, actual %v", cl.Email, client.Email)

	// Удаление клиента из базы данных
	err = deleteClient(db, client.ID)
	require.NoError(t, err, "error deleting client with ID %d: %v", cl.ID, err)

	// Проверка того, что клиент действительно удален
	_, err = selectClient(db, client.ID)
	require.Error(t, err, "expected error when trying to retrieve deleted client with ID %d", client.ID)
	require.Equal(t, sql.ErrNoRows, err, "expected specific sql.ErrNoRows error when searching for deleted client with ID %d", client.ID)
}
