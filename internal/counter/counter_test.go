package counter_test

import (
	"fmt"
	"sync"
	"testing"

	service "github.com/kl09/counter-api"
	"github.com/kl09/counter-api/internal/counter"
	"github.com/kl09/counter-api/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestNewStats_Sync(t *testing.T) {
	s := counter.NewStats(
		&mock.EventSenderMock{
			SendFunc: func(event service.Event) error {
				return nil
			},
		},
	)

	require.Equal(t, int64(0), s.Get("a1"))
	require.Equal(t, int64(0), s.Get("a1"))

	s.Increment("a1")
	require.Equal(t, int64(1), s.Get("a1"))

	s.Increment("a1")
	require.Equal(t, int64(2), s.Get("a1"))

	s.Decrement("a1")
	require.Equal(t, int64(1), s.Get("a1"))

	s.Decrement("a1")
	require.Equal(t, int64(0), s.Get("a1"))

	s.Decrement("a1")
	require.Equal(t, int64(-1), s.Get("a1"))

	s.Reset("a1")
	require.Equal(t, int64(0), s.Get("a1"))
}

func TestNewStats_Parallel(t *testing.T) {
	t.Parallel()

	s := counter.NewStats(
		&mock.EventSenderMock{
			SendFunc: func(event service.Event) error {
				return nil
			},
		},
	)

	keys := []string{"a1", "a2", "a3"}

	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(2)
		go func(k string) {
			defer wg.Done()

			for i := 0; i < 1000; i++ {
				s.Increment(k)
			}
		}(key)

		go func(k string) {
			defer wg.Done()

			for i := 0; i < 500; i++ {
				s.Decrement(k)
			}
		}(key)
	}

	wg.Wait()

	for _, key := range keys {
		require.Equal(t, int64(500), s.Get(key), fmt.Sprintf("key; %s", key))
	}
}
