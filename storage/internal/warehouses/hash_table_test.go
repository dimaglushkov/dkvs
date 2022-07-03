package warehouses_test

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dimaglushkov/dkvs/storage/internal"
	"github.com/dimaglushkov/dkvs/storage/internal/warehouses"
)

const (
	maxValue  = 100
	testCases = 60
	cf        = 40 // number of goroutines for the concurrent test
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func itoa(x int) string {
	return strconv.FormatInt(int64(x), 10)
}

func TestHashTable(t *testing.T) {
	ht := warehouses.NewHashTable()
	m := make(map[string]string)

	// Setting up random variables
	for i := 0; i < testCases; i++ {
		k, v := itoa(rand.Intn(maxValue)), itoa(rand.Intn(maxValue))

		m[k] = v
		assert.NoErrorf(t, ht.Put(k, v), "error while putting: {%s: %s}", k, v)
	}

	// Checking stored values
	for i := 0; i < maxValue; i++ {
		k := itoa(i)
		expected, ok := m[k]
		actual, err := ht.Get(k)

		if !ok {
			assert.IsTypef(t, err, internal.UnknownKeyError{}, "unknown error occured: %s", err)
		} else {
			assert.NoErrorf(t, err, "unexpected error while getting the value: {%s}", k)
		}

		assert.Equal(t, expected, actual)
	}

	// Deleting values
	for i := 0; i < maxValue; i++ {
		k := itoa(i)
		err := ht.Delete(k)

		if _, ok := m[k]; !ok {
			assert.IsTypef(t, err, internal.UnknownKeyError{}, "unknown error occured: %s", err)
		} else {
			assert.NoErrorf(t, err, "unexpected error while getting the value: {%s}", k)
		}
	}

	// Checking that every value was removed successfully
	for i := 0; i < maxValue; i++ {
		k := itoa(i)
		_, err := ht.Get(k)
		assert.IsTypef(t, err, internal.UnknownKeyError{}, "unknown error occured: %s", err)
	}

}

func TestHashTableConcurrent(t *testing.T) {
	require.LessOrEqualf(t, testCases, maxValue, "value of testCases should be less than maxValue")

	ht := warehouses.NewHashTable()
	var wg sync.WaitGroup

	// Concurrently filling the storage
	for i := 0; i < cf; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			for i := 0; i < testCases; i++ {
				v := itoa(rand.Intn(maxValue))
				k := itoa(rand.Intn(maxValue) + maxValue*n)

				prev, _ := ht.Get(k)
				for prev != "" {
					k = itoa(rand.Intn(maxValue) + maxValue*n)
					prev, _ = ht.Get(k)
				}

				assert.NoErrorf(t, ht.Put(k, v), "error while putting: {%s: %s}", k, v)
			}
		}(i)
	}
	wg.Wait()

	// Concurrently getting the values
	for i := 0; i < cf; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			cnt := 0
			for i := maxValue * n; i < maxValue*(n+1); i++ {
				if v, err := ht.Get(itoa(i)); v != "" && err == nil {
					cnt++
				}
			}
			assert.Equalf(t, testCases, cnt, "missmatching number of records for goroutine %d - expected %d, got %d", n, testCases, cnt)
		}(i)
	}
	wg.Wait()

	// Concurrently deleting the values
	for i := 0; i < cf; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			cnt := 0
			for i := maxValue * n; i < maxValue*(n+1); i++ {
				if err := ht.Delete(itoa(i)); err == nil {
					cnt++
				}
			}
			assert.Equalf(t, testCases, cnt, "missmatching number of records for goroutine %d - expected %d, got %d", n, testCases, cnt)
		}(i)
	}
	wg.Wait()

}
