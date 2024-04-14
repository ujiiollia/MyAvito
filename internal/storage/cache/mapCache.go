package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *Cache) Set(id string, banner interface{}) {

	expiration := time.Now().Add(c.defaultExpiration).UnixNano()

	c.Lock()

	defer c.Unlock()

	c.items[id] = Item{
		Value:      banner,
		Expiration: expiration,
		Created:    time.Now(),
	}

}
func (c *Cache) Get(id string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[id]

	// ключ не найден
	if !found {
		return nil, false
	}

	// Проверка на установку времени истечения, в противном случае он бессрочный
	if item.Expiration > 0 {

		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *Cache) UpdateUserBannerCache(pool *pgxpool.Pool) {
	ctx := context.Background()

	query := `SELECT tag_id, feature_id, content FROM user_banner WHERE is_active = true`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Query Failed: %v", err)
	}
	defer rows.Close()

	newCacheItems := make(map[string]Item)

	for rows.Next() {
		var tagID, featureID int
		var content json.RawMessage
		if err := rows.Scan(&tagID, &featureID, &content); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
			continue
		}

		key := fmt.Sprintf("tag_%d:feature_%d", tagID, featureID)

		newCacheItems[key] = Item{
			Value:      content,
			Created:    time.Now(),
			Expiration: int64(c.defaultExpiration.Seconds()),
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Rows error: %v", err)
	}

	c.items = newCacheItems
}
