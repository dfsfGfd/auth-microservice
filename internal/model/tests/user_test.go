package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"auth-microservice/internal/model"
	"auth-microservice/pkg/bcrypt"
)

func TestNewUser(t *testing.T) {
	hasher := bcrypt.NewService()

	createPasswordHash := func(t *testing.T, password string) *model.PasswordHash {
		t.Helper()
		hash, err := hasher.Hash(password, 0)
		require.NoError(t, err)
		passwordHash, err := model.NewPasswordHash(hash)
		require.NoError(t, err)
		return passwordHash
	}

	createEmail := func(t *testing.T, email string) *model.Email {
		t.Helper()
		e, err := model.NewEmail(email)
		require.NoError(t, err)
		return e
	}

	createUsername := func(t *testing.T, username string) *model.Username {
		t.Helper()
		u, err := model.NewUsername(username)
		require.NoError(t, err)
		return u
	}

	t.Run("успешное создание пользователя", func(t *testing.T) {
		email := createEmail(t, "test@example.com")
		username := createUsername(t, "testuser")
		passwordHash := createPasswordHash(t, "Password123")

		user, err := model.NewUser(email, username, passwordHash)

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotEmpty(t, user.ID())
		assert.Equal(t, "test@example.com", user.Email().String())
		assert.Equal(t, "testuser", user.Username().String())
		assert.NotEmpty(t, user.PasswordHash().String())
		assert.NotEmpty(t, user.CreatedAt())
		assert.NotEmpty(t, user.UpdatedAt())
	})

	t.Run("nil email", func(t *testing.T) {
		username := createUsername(t, "testuser")
		passwordHash := createPasswordHash(t, "Password123")

		user, err := model.NewUser(nil, username, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("nil username", func(t *testing.T) {
		email := createEmail(t, "test@example.com")
		passwordHash := createPasswordHash(t, "Password123")

		user, err := model.NewUser(email, nil, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("nil passwordHash", func(t *testing.T) {
		email := createEmail(t, "test@example.com")
		username := createUsername(t, "testuser")

		user, err := model.NewUser(email, username, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUser_UpdatePasswordHash(t *testing.T) {
	hasher := bcrypt.NewService()

	createUser := func(t *testing.T) *model.User {
		t.Helper()
		email, _ := model.NewEmail("test@example.com")
		username, _ := model.NewUsername("testuser")
		hash, _ := hasher.Hash("Password123", 0)
		passwordHash, _ := model.NewPasswordHash(hash)
		user, _ := model.NewUser(email, username, passwordHash)
		return user
	}

	t.Run("успешное обновление хеша пароля", func(t *testing.T) {
		user := createUser(t)

		oldUpdatedAt := user.UpdatedAt()
		oldHash := user.PasswordHash().Value()

		newHash, err := hasher.Hash("NewPassword456", 0)
		require.NoError(t, err)
		newPasswordHash, err := model.NewPasswordHash(newHash)
		require.NoError(t, err)

		err = user.UpdatePasswordHash(newPasswordHash)

		require.NoError(t, err)
		assert.NotEqual(t, oldHash, user.PasswordHash().Value())
		assert.Greater(t, user.UpdatedAt(), oldUpdatedAt)
	})

	t.Run("обновление nil возвращает ошибку", func(t *testing.T) {
		user := createUser(t)

		oldHash := user.PasswordHash().Value()
		oldUpdatedAt := user.UpdatedAt()

		err := user.UpdatePasswordHash(nil)

		assert.Error(t, err)
		assert.Equal(t, oldHash, user.PasswordHash().Value())
		assert.Equal(t, oldUpdatedAt, user.UpdatedAt())
	})
}

func TestUser_UpdateEmail(t *testing.T) {
	createUser := func(t *testing.T) *model.User {
		t.Helper()
		email, _ := model.NewEmail("test@example.com")
		username, _ := model.NewUsername("testuser")
		hash, _ := bcrypt.NewService().Hash("Password123", 0)
		passwordHash, _ := model.NewPasswordHash(hash)
		user, _ := model.NewUser(email, username, passwordHash)
		return user
	}

	t.Run("успешное обновление email", func(t *testing.T) {
		user := createUser(t)
		oldUpdatedAt := user.UpdatedAt()

		newEmail, err := model.NewEmail("new@example.com")
		require.NoError(t, err)

		err = user.UpdateEmail(newEmail)

		require.NoError(t, err)
		assert.Equal(t, "new@example.com", user.Email().Value())
		assert.Greater(t, user.UpdatedAt(), oldUpdatedAt)
	})

	t.Run("обновление nil email возвращает ошибку", func(t *testing.T) {
		user := createUser(t)
		oldEmail := user.Email().Value()

		err := user.UpdateEmail(nil)

		assert.Error(t, err)
		assert.Equal(t, oldEmail, user.Email().Value())
	})

	t.Run("обновление невалидным email возвращает ошибку", func(t *testing.T) {
		user := createUser(t)
		oldEmail := user.Email().Value()

		invalidEmail, _ := model.NewEmail("invalid-email")
		err := user.UpdateEmail(invalidEmail)

		assert.Error(t, err)
		assert.Equal(t, oldEmail, user.Email().Value())
	})
}

func TestUser_UpdateUsername(t *testing.T) {
	createUser := func(t *testing.T) *model.User {
		t.Helper()
		email, _ := model.NewEmail("test@example.com")
		username, _ := model.NewUsername("testuser")
		hash, _ := bcrypt.NewService().Hash("Password123", 0)
		passwordHash, _ := model.NewPasswordHash(hash)
		user, _ := model.NewUser(email, username, passwordHash)
		return user
	}

	t.Run("успешное обновление username", func(t *testing.T) {
		user := createUser(t)

		newUsername, err := model.NewUsername("newuser")
		require.NoError(t, err)

		err = user.UpdateUsername(newUsername)

		require.NoError(t, err)
		assert.Equal(t, "newuser", user.Username().Value())
	})

	t.Run("обновление nil username возвращает ошибку", func(t *testing.T) {
		user := createUser(t)
		oldUsername := user.Username().Value()

		err := user.UpdateUsername(nil)

		assert.Error(t, err)
		assert.Equal(t, oldUsername, user.Username().Value())
	})
}

func TestUser_Getters(t *testing.T) {
	t.Run("CreatedAt и UpdatedAt установлены", func(t *testing.T) {
		email, _ := model.NewEmail("test@example.com")
		username, _ := model.NewUsername("testuser")
		hash, _ := bcrypt.NewService().Hash("Password123", 0)
		passwordHash, _ := model.NewPasswordHash(hash)
		user, _ := model.NewUser(email, username, passwordHash)

		assert.NotEmpty(t, user.CreatedAt())
		assert.NotEmpty(t, user.UpdatedAt())
		assert.Equal(t, user.CreatedAt(), user.UpdatedAt())
	})
}
