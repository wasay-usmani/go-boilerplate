package cachekit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valkey-io/valkey-go/mock"
	"go.uber.org/mock/gomock"
)

func setupMockClient(t *testing.T) (*mock.Client, Cache) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewClient(ctrl)
	cache := &valkeyClient{client: mockClient, keySpace: "test_keyspace"}
	return mockClient, cache
}

func TestSet(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	key := "test_key"
	val := "test_value"
	exp := 10 * time.Second

	mockClient.EXPECT().
		Do(ctx, mock.Match("SET", fmt.Sprintf("%s::%s", "test_keyspace", key), val, "EX", "10")).
		Return(mock.Result(mock.ValkeyString("OK")))

	err := cache.Save(ctx, key, val, exp)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	key := "test_key"
	val := "test_value"

	mockClient.EXPECT().
		Do(ctx, mock.Match("GET", fmt.Sprintf("%s::%s", "test_keyspace", key))).
		Return(mock.Result(mock.ValkeyString(val)))

	result, err := cache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, val, result)
}

func TestRemoveKeys(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	keys := []string{"key1", "key2"}

	mockClient.EXPECT().
		Do(ctx, mock.Match("DEL", "test_keyspace::key1", "test_keyspace::key2")).
		Return(mock.Result(mock.ValkeyInt64(2)))

	err := cache.RemoveKeys(ctx, keys...)
	assert.NoError(t, err)
}

func TestIsKeyActive(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	key := "test_key"

	mockClient.EXPECT().
		Do(ctx, mock.Match("EXISTS", fmt.Sprintf("%s::%s", "test_keyspace", key))).
		Return(mock.Result(mock.ValkeyInt64(1)))

	result, err := cache.IsKeyActive(ctx, key)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestPing(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()

	mockClient.EXPECT().Do(ctx, mock.Match("PING")).Return(mock.Result(mock.ValkeyString("PONG")))

	err := cache.Ping(ctx)
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	mockClient.EXPECT().Close()

	cache.Close()
}

func TestAddSet(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	set := "my_set"
	value := "my_value"

	mockClient.EXPECT().Do(ctx, mock.Match("SADD", fmt.Sprintf("%s::%s", "test_keyspace", set), value)).Return(mock.Result(mock.ValkeyInt64(1)))

	_, err := cache.AddSet(ctx, set, value)
	assert.NoError(t, err)
}

func TestRemoveSet(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	set := "my_set"

	mockClient.EXPECT().Do(ctx, mock.Match("DEL", fmt.Sprintf("%s::%s", "test_keyspace", set))).Return(mock.Result(mock.ValkeyInt64(1)))

	err := cache.RemoveSet(ctx, set)
	assert.NoError(t, err)
}

func TestContainsSet(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	set := "my_set"
	value := "my_value"

	mockClient.EXPECT().Do(ctx, mock.Match("SISMEMBER", fmt.Sprintf("%s::%s", "test_keyspace", set), value)).Return(mock.Result(mock.ValkeyInt64(1)))

	exists, err := cache.ContainsSet(ctx, set, value)
	assert.NoError(t, err)
	assert.Equal(t, true, exists)
}

func TestRemoveSetValues(t *testing.T) {
	mockClient, cache := setupMockClient(t)
	ctx := context.Background()
	set := "my_set"
	values := []string{"value1", "value2"}

	mockClient.EXPECT().Do(ctx, mock.Match("SREM", fmt.Sprintf("%s::%s", "test_keyspace", set), values[0], values[1])).Return(mock.Result(mock.ValkeyInt64(2)))

	n, err := cache.RemoveSetValues(ctx, set, values...)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)
}
