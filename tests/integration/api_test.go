package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/niamfauzi/go-starter/internal/modules/auth"
	"github.com/niamfauzi/go-starter/internal/modules/product"
	"github.com/niamfauzi/go-starter/internal/modules/user"
	"github.com/niamfauzi/go-starter/internal/platform/db"
	platformhttp "github.com/niamfauzi/go-starter/internal/platform/http"
	sharedvalidator "github.com/niamfauzi/go-starter/internal/shared/validator"
)

func TestAPIIntegration(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctx := context.Background()

	mysqlContainer, mysqlDSN := startMySQL(t, ctx)
	defer func() { _ = mysqlContainer.Terminate(ctx) }()

	redisContainer, redisAddr := startRedis(t, ctx)
	defer func() { _ = redisContainer.Terminate(ctx) }()

	sqlDB, err := db.NewSQL(mysqlDSN)
	require.NoError(t, err)
	defer sqlDB.Close()

	migrator := db.NewMigrator(sqlDB)
	require.NoError(t, migrator.Up(projectPath("migrations")))

	appLogger := zap.NewNop()
	gormDB, err := db.NewGorm(mysqlDSN, appLogger)
	require.NoError(t, err)

	seedTestData(t, gormDB)

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	validate := sharedvalidator.New()

	userRepo := user.NewRepository(gormDB)
	authRepo := auth.NewRepository(userRepo)
	jwtManager := auth.NewJWTManager("test-secret", "go-starter-test", 15*time.Minute)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHTTPDelivery(authService, validate, appLogger)

	productRepo := product.NewRepository(gormDB)
	productCache := product.NewCache(redisClient, 30*time.Second)
	productService := product.NewService(productRepo, productCache, appLogger)
	productHandler := product.NewHTTPDelivery(productService, validate, appLogger)

	router := platformhttp.NewRouter(appLogger, authHandler, productHandler, jwtManager)
	server := httptest.NewServer(router)
	defer server.Close()

	token := loginAndGetToken(t, server.URL, "tenant-demo", "demo@example.com", "password123")

	t.Run("login sukses", func(t *testing.T) {
		body := map[string]any{
			"email":    "demo@example.com",
			"password": "password123",
		}
		resp := doJSON(t, http.MethodPost, server.URL+"/api/v1/auth/login", "tenant-demo", "", body)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("login gagal", func(t *testing.T) {
		body := map[string]any{
			"email":    "demo@example.com",
			"password": "salah123",
		}
		resp := doJSON(t, http.MethodPost, server.URL+"/api/v1/auth/login", "tenant-demo", "", body)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	var productID uint64

	t.Run("create product dengan JWT valid", func(t *testing.T) {
		body := map[string]any{
			"name":        "Produk Test",
			"description": "dibuat dari integration test",
			"price":       12345,
			"stock":       9,
		}
		resp := doJSON(t, http.MethodPost, server.URL+"/api/v1/products", "tenant-demo", token, body)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var out map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))

		data := out["data"].(map[string]any)
		productID = uint64(data["id"].(float64))
	})

	t.Run("create product tanpa JWT", func(t *testing.T) {
		body := map[string]any{
			"name":        "Produk Gagal",
			"description": "harus gagal",
			"price":       1000,
			"stock":       1,
		}
		resp := doJSON(t, http.MethodPost, server.URL+"/api/v1/products", "tenant-demo", "", body)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get product", func(t *testing.T) {
		resp := doJSON(t, http.MethodGet, fmt.Sprintf("%s/api/v1/products/%d", server.URL, productID), "tenant-demo", token, nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("update product", func(t *testing.T) {
		body := map[string]any{
			"name":        "Produk Updated",
			"description": "sudah diupdate",
			"price":       99999,
			"stock":       20,
		}
		resp := doJSON(t, http.MethodPut, fmt.Sprintf("%s/api/v1/products/%d", server.URL, productID), "tenant-demo", token, body)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("tenant isolation", func(t *testing.T) {
		otherToken := loginAndGetToken(t, server.URL, "tenant-other", "other@example.com", "password123")
		resp := doJSON(t, http.MethodGet, fmt.Sprintf("%s/api/v1/products/%d", server.URL, productID), "tenant-other", otherToken, nil)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("delete product", func(t *testing.T) {
		resp := doJSON(t, http.MethodDelete, fmt.Sprintf("%s/api/v1/products/%d", server.URL, productID), "tenant-demo", token, nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func startMySQL(t *testing.T, ctx context.Context) (testcontainers.Container, string) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.4",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "go_starter_test",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("ready for connections"),
			wait.ForListeningPort("3306/tcp"),
		).WithDeadline(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started:          true,
		ContainerRequest: req,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/go_starter_test?parseTime=true", host, port.Port())
	return container, dsn
}

func startRedis(t *testing.T, ctx context.Context) (testcontainers.Container, string) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7.4-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started:          true,
		ContainerRequest: req,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)

	return container, fmt.Sprintf("%s:%s", host, port.Port())
}

func seedTestData(t *testing.T, gormDB *gorm.DB) {
	t.Helper()

	hash, err := auth.HashPassword("password123")
	require.NoError(t, err)

	users := []user.Entity{
		{TenantID: "tenant-demo", Name: "Demo User", Email: "demo@example.com", PasswordHash: string(hash), Role: "admin", IsActive: true},
		{TenantID: "tenant-other", Name: "Other User", Email: "other@example.com", PasswordHash: string(hash), Role: "admin", IsActive: true},
	}

	for _, u := range users {
		require.NoError(t, gormDB.Create(&u).Error)
	}
}

func loginAndGetToken(t *testing.T, baseURL string, tenantID string, email string, password string) string {
	t.Helper()

	body := map[string]any{
		"email":    email,
		"password": password,
	}

	resp := doJSON(t, http.MethodPost, baseURL+"/api/v1/auth/login", tenantID, "", body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))

	data := out["data"].(map[string]any)
	return data["access_token"].(string)
}

func doJSON(t *testing.T, method string, url string, tenantID string, token string, body any) *http.Response {
	t.Helper()

	var reqBody bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&reqBody).Encode(body))
	}

	req, err := http.NewRequest(method, url, &reqBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func projectPath(rel string) string {
	wd, _ := os.Getwd()
	return filepath.Clean(filepath.Join(wd, "..", "..", rel))
}
