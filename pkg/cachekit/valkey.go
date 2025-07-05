package cachekit

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/mock"
	"github.com/wasay-usmani/go-boilerplate/pkg/configkit"
)

type valkeyClient struct {
	client   valkey.Client
	keySpace string
}

func NewCache(ctx context.Context, opt *configkit.Cache) (Cache, error) {
	options := valkey.ClientOption{
		InitAddress:  opt.Addresses,
		ClientName:   opt.ClientName,
		Username:     opt.Username,
		Password:     opt.Password,
		DisableCache: true,
	}

	if opt.TLSEnabled != nil && *opt.TLSEnabled {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	c, err := valkey.NewClient(options)
	if err != nil {
		return nil, err
	}

	client := &valkeyClient{client: c, keySpace: opt.KeySpace}
	err = client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMockClient(c *mock.Client) Cache {
	return &valkeyClient{client: c}
}

func (v *valkeyClient) Save(ctx context.Context, key, value string, ttl time.Duration) error {
	cmd := v.client.B().Set().Key(v.formatKey(key)).Value(value)
	if ttl > 0 {
		cmd.Ex(ttl)
	}

	return v.client.Do(ctx, cmd.Build()).Error()
}

func (v *valkeyClient) SaveKeepTTL(ctx context.Context, key, value string) error {
	cmd := v.client.B().Set().Key(v.formatKey(key)).Value(value).Keepttl()
	return v.client.Do(ctx, cmd.Build()).Error()
}

func (v *valkeyClient) RemoveKeys(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	for i := range keys {
		keys[i] = v.formatKey(keys[i])
	}

	cmd := v.client.B().Del().Key(keys...).Build()
	return v.client.Do(ctx, cmd).Error()
}

func (v *valkeyClient) Get(ctx context.Context, key string) (string, error) {
	cmd := v.client.B().Get().Key(v.formatKey(key)).Build()
	return v.client.Do(ctx, cmd).ToString()
}

func (v *valkeyClient) IsKeyActive(ctx context.Context, key string) (bool, error) {
	cmd := v.client.B().Exists().Key(v.formatKey(key)).Build()
	count, err := v.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (v *valkeyClient) AddSet(ctx context.Context, set string, value ...string) (int64, error) {
	cmd := v.client.B().Sadd().Key(v.formatKey(set)).Member(value...).Build()
	return v.client.Do(ctx, cmd).ToInt64()
}

func (v *valkeyClient) RemoveSetValues(ctx context.Context, set string, values ...string) (int64, error) {
	cmd := v.client.B().Srem().Key(v.formatKey(set)).Member(values...).Build()
	return v.client.Do(ctx, cmd).ToInt64()
}

func (v *valkeyClient) RemoveSet(ctx context.Context, set string) error {
	cmd := v.client.B().Del().Key(v.formatKey(set)).Build()
	return v.client.Do(ctx, cmd).Error()
}

func (v *valkeyClient) ContainsSet(ctx context.Context, set, value string) (bool, error) {
	cmd := v.client.B().Sismember().Key(v.formatKey(set)).Member(value).Build()
	res, err := v.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func (v *valkeyClient) Ping(ctx context.Context) error {
	cmd := v.client.B().Ping().Build()
	_, err := v.client.Do(ctx, cmd).ToString()
	return err
}

func (v *valkeyClient) Close() {
	v.client.Close()
}

func (v *valkeyClient) formatKey(key string) string {
	return fmt.Sprintf("%s::%s", v.keySpace, key)
}
