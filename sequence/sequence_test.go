package sequence

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(s.Close)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { rdb.Close() })
	return rdb
}

func TestNext_Format(t *testing.T) {
	rdb := newTestRedis(t)
	g := New(rdb, "test")
	ctx := context.Background()

	id, err := g.Next(ctx, "SH")
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}

	date := time.Now().Format("20060102")
	if !strings.HasPrefix(id, "SH"+date) {
		t.Errorf("id = %q, want prefix SH%s", id, date)
	}
	if id != "SH"+date+"001" {
		t.Errorf("first id = %q, want SH%s001", id, date)
	}
}

func TestNext_Sequential(t *testing.T) {
	rdb := newTestRedis(t)
	g := New(rdb, "test")
	ctx := context.Background()

	id1, _ := g.Next(ctx, "SH")
	id2, _ := g.Next(ctx, "SH")
	id3, _ := g.Next(ctx, "SH")

	date := time.Now().Format("20060102")
	if id1 != "SH"+date+"001" {
		t.Errorf("id1 = %q", id1)
	}
	if id2 != "SH"+date+"002" {
		t.Errorf("id2 = %q", id2)
	}
	if id3 != "SH"+date+"003" {
		t.Errorf("id3 = %q", id3)
	}
}

func TestNext_DifferentPrefixes(t *testing.T) {
	rdb := newTestRedis(t)
	g := New(rdb, "test")
	ctx := context.Background()

	sh1, _ := g.Next(ctx, "SH")
	ord1, _ := g.Next(ctx, "ORD")
	sh2, _ := g.Next(ctx, "SH")

	date := time.Now().Format("20060102")
	if sh1 != "SH"+date+"001" {
		t.Errorf("sh1 = %q", sh1)
	}
	if ord1 != "ORD"+date+"001" {
		t.Errorf("ord1 = %q", ord1)
	}
	if sh2 != "SH"+date+"002" {
		t.Errorf("sh2 = %q", sh2)
	}
}
