package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache menampung operasi cache khusus product.
// Kita pisahkan agar service tidak penuh dengan detail key Redis.
type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		client: client,
		ttl:    ttl,
	}
}

func (c *Cache) GetList(ctx context.Context, tenantID string) ([]Response, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}

	raw, err := c.client.Get(ctx, c.listKey(tenantID)).Result()
	if err != nil || raw == "" {
		return nil, false
	}

	var data []Response
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, false
	}

	return data, true
}

func (c *Cache) SetList(ctx context.Context, tenantID string, data []Response) {
	if c == nil || c.client == nil {
		return
	}

	if raw, err := json.Marshal(data); err == nil {
		_ = c.client.Set(ctx, c.listKey(tenantID), string(raw), c.ttl).Err()
	}
}

func (c *Cache) GetDetail(ctx context.Context, tenantID string, productID uint64) (*Response, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}

	raw, err := c.client.Get(ctx, c.detailKey(tenantID, productID)).Result()
	if err != nil || raw == "" {
		return nil, false
	}

	var data Response
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, false
	}

	return &data, true
}

func (c *Cache) SetDetail(ctx context.Context, tenantID string, productID uint64, data Response) {
	if c == nil || c.client == nil {
		return
	}

	if raw, err := json.Marshal(data); err == nil {
		_ = c.client.Set(ctx, c.detailKey(tenantID, productID), string(raw), c.ttl).Err()
	}
}

func (c *Cache) Invalidate(ctx context.Context, tenantID string, productID uint64) error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Del(ctx, c.listKey(tenantID), c.detailKey(tenantID, productID)).Err()
}

func (c *Cache) listKey(tenantID string) string {
	return fmt.Sprintf("products:list:%s", tenantID)
}

func (c *Cache) detailKey(tenantID string, productID uint64) string {
	return fmt.Sprintf("products:get:%s:%d", tenantID, productID)
}
