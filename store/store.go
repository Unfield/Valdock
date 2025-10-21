package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Unfield/Valdock/models"
	"github.com/valkey-io/valkey-go"
)

type Store struct {
	kv valkey.Client
}

func NewStore(kv valkey.Client) *Store {
	return &Store{
		kv: kv,
	}
}

func (s *Store) DeleteKey(key string) error {
	ctx := context.Background()
	return s.kv.Do(ctx,
		s.kv.B().Del().Key(key).Build(),
	).Error()
}

func (s *Store) SetString(key string, value string) error {
	ctx := context.Background()
	return s.kv.Do(ctx,
		s.kv.B().Set().Key(key).Value(value).Build(),
	).Error()
}

func (s *Store) GetString(key string) (string, error) {
	ctx := context.Background()
	str, err := s.kv.Do(ctx, s.kv.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", errors.New("not found")
		}
		return "", fmt.Errorf("get failed: %w", err)
	}
	return str, nil
}

func (s *Store) SetInt(key string, value int) error {
	ctx := context.Background()
	strValue := strconv.Itoa(value)

	return s.kv.Do(ctx,
		s.kv.B().Set().Key(key).Value(strValue).Build(),
	).Error()
}

func (s *Store) GetInt(key string) (int, error) {
	ctx := context.Background()

	str, err := s.kv.Do(ctx, s.kv.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return 0, errors.New("not found")
		}
		return 0, fmt.Errorf("get failed: %w", err)
	}

	v, convErr := strconv.Atoi(str)
	if convErr != nil {
		return 0, fmt.Errorf("parse failed: %w", convErr)
	}

	return v, nil
}

func (s *Store) SetJSON(key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	ctx := context.Background()
	return s.kv.Do(ctx,
		s.kv.B().Set().Key(key).Value(string(b)).Build(),
	).Error()
}

func (s *Store) GetJSON(key string, target any) error {
	ctx := context.Background()
	str, err := s.kv.Do(ctx, s.kv.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return errors.New("not found")
		}
		return fmt.Errorf("get failed: %w", err)
	}

	if err := json.Unmarshal([]byte(str), target); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	return nil
}

func (s *Store) ListInstances() ([]models.InstanceModel, error) {
	var results []models.InstanceModel
	var cursor uint64 = 0
	ctx := context.Background()

	for {
		res, err := s.kv.Do(ctx, s.kv.B().Scan().Cursor(cursor).
			Match("instances:*").Count(100).Build()).AsScanEntry()
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		keys := res.Elements
		for _, key := range keys {
			var currentVal models.InstanceModel
			s.GetJSON(key, &currentVal)
			results = append(results, currentVal)
		}

		cursor = res.Cursor
		if cursor == 0 {
			break
		}
	}

	return results, nil
}
