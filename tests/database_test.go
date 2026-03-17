package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseConnection проверяет подключение к PostgreSQL из контейнера.
func TestDatabaseConnection(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("DATABASE_URL")
	require.NotEmpty(t, dsn, "DATABASE_URL must be set")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создание пула подключений
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "Failed to create connection pool")
	defer pool.Close()

	// Проверка подключения
	err = pool.Ping(ctx)
	assert.NoError(t, err, "Failed to ping database")

	t.Log("✅ Database connection successful")
}

// TestRedisConnection проверяет подключение к Redis из контейнера.
func TestRedisConnection(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("REDIS_URL")
	require.NotEmpty(t, dsn, "REDIS_URL must be set")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создание Redis клиента
	opt, err := redis.ParseURL(dsn)
	require.NoError(t, err, "Failed to parse Redis URL")

	client := redis.NewClient(opt)
	defer client.Close()

	// Проверка подключения
	err = client.Ping(ctx).Err()
	assert.NoError(t, err, "Failed to ping Redis")

	t.Log("✅ Redis connection successful")
}

// TestDatabaseTables проверяет существование таблиц после миграций.
func TestDatabaseTables(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// Проверка существования таблицы accounts
	var exists bool
	err = pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'accounts'
		)
	`).Scan(&exists)

	require.NoError(t, err)
	assert.True(t, exists, "Table 'accounts' should exist")

	t.Log("✅ Table 'accounts' exists")
}

// TestDatabaseIndexes проверяет существование индексов.
func TestDatabaseIndexes(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// Проверка индекса idx_accounts_email
	var count int
	err = pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_indexes 
		WHERE schemaname = 'public' 
		AND tablename = 'accounts' 
		AND indexname = 'idx_accounts_email'
	`).Scan(&count)

	require.NoError(t, err)
	assert.Greater(t, count, 0, "Index 'idx_accounts_email' should exist")

	t.Log("✅ Index 'idx_accounts_email' exists")
}

// TestDatabaseTrigger проверяет существование триггера.
func TestDatabaseTrigger(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// Проверка триггера update_accounts_updated_at
	var count int
	err = pool.QueryRow(ctx, `
		SELECT count(*) FROM information_schema.triggers 
		WHERE trigger_schema = 'public' 
		AND trigger_name = 'update_accounts_updated_at'
	`).Scan(&count)

	require.NoError(t, err)
	assert.Greater(t, count, 0, "Trigger 'update_accounts_updated_at' should exist")

	t.Log("✅ Trigger 'update_accounts_updated_at' exists")
}

// TestDatabaseFunction проверяет существование функции.
func TestDatabaseFunction(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// Проверка функции update_updated_at_column
	var count int
	err = pool.QueryRow(ctx, `
		SELECT count(*) FROM information_schema.routines 
		WHERE routine_schema = 'public' 
		AND routine_name = 'update_updated_at_column'
	`).Scan(&count)

	require.NoError(t, err)
	assert.Greater(t, count, 0, "Function 'update_updated_at_column' should exist")

	t.Log("✅ Function 'update_updated_at_column' exists")
}

// TestRedisOperations тестирует операции с Redis.
func TestRedisOperations(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("REDIS_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt, err := redis.ParseURL(dsn)
	require.NoError(t, err)

	client := redis.NewClient(opt)
	defer client.Close()

	// Тест Set/Get
	t.Run("set_get", func(t *testing.T) {
		key := "test:key"
		value := "test_value"

		err := client.Set(ctx, key, value, time.Minute).Err()
		require.NoError(t, err)

		got, err := client.Get(ctx, key).Result()
		require.NoError(t, err)
		assert.Equal(t, value, got)

		// Очистка
		client.Del(ctx, key)
		t.Log("✅ Redis SET/GET successful")
	})

	// Тест TTL
	t.Run("ttl", func(t *testing.T) {
		key := "test:ttl"
		value := "ttl_value"

		err := client.Set(ctx, key, value, 5*time.Second).Err()
		require.NoError(t, err)

		ttl, err := client.TTL(ctx, key).Result()
		require.NoError(t, err)
		assert.Greater(t, ttl, time.Duration(0))
		assert.Less(t, ttl, 6*time.Second)

		client.Del(ctx, key)
		t.Log("✅ Redis TTL successful")
	})

	// Тест Delete
	t.Run("delete", func(t *testing.T) {
		key := "test:delete"

		err := client.Set(ctx, key, "value", time.Minute).Err()
		require.NoError(t, err)

		err = client.Del(ctx, key).Err()
		require.NoError(t, err)

		_, err = client.Get(ctx, key).Result()
		assert.Error(t, err, redis.Nil)

		t.Log("✅ Redis DELETE successful")
	})
}

// TestConcurrentDatabaseAccess тестирует конкурентный доступ к БД.
func TestConcurrentDatabaseAccess(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// Запускаем 10 горутин с запросами
	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			var result int
			err := pool.QueryRow(ctx, "SELECT $1", id).Scan(&result)
			assert.NoError(t, err)
			assert.Equal(t, id, result)
		}(i)
	}

	// Ожидаем завершения всех горутин
	for i := 0; i < goroutines; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for goroutines")
		}
	}

	t.Logf("✅ Concurrent access test passed (%d goroutines)", goroutines)
}

// TestConcurrentRedisAccess тестирует конкурентный доступ к Redis.
func TestConcurrentRedisAccess(t *testing.T) {
	t.Parallel()

	dsn := os.Getenv("REDIS_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opt, err := redis.ParseURL(dsn)
	require.NoError(t, err)

	client := redis.NewClient(opt)
	defer client.Close()

	// Запускаем 10 горутин с запросами
	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			key := fmt.Sprintf("test:concurrent:%d", id)
			value := fmt.Sprintf("value_%d", id)

			err := client.Set(ctx, key, value, time.Minute).Err()
			assert.NoError(t, err)

			got, err := client.Get(ctx, key).Result()
			assert.NoError(t, err)
			assert.Equal(t, value, got)

			client.Del(ctx, key)
		}(i)
	}

	// Ожидаем завершения всех горутин
	for i := 0; i < goroutines; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for goroutines")
		}
	}

	t.Logf("✅ Concurrent Redis access test passed (%d goroutines)", goroutines)
}
