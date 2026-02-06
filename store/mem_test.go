package store

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemStorage(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	assert.NotNil(t, s)
	assert.Empty(t, s.KeySet())
}

func TestMemStorage_SetGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		value []byte
	}{
		{"simple", "key1", []byte("value1")},
		{"empty_value", "key2", []byte{}},
		{"nil_value", "key3", nil},
		{"binary_data", "bin", []byte{0x00, 0xFF, 0x80}},
		{"unicode_key", "café", []byte("latte")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStore()
			err := s.Set(tt.key, tt.value)
			require.NoError(t, err)

			got, ok, err := s.Get(tt.key)
			require.NoError(t, err)
			assert.True(t, ok)
			assert.Equal(t, tt.value, got)
		})
	}
}

func TestMemStorage_GetMissing(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	got, ok, err := s.Get("missing")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestMemStorage_SetOverwrite(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	require.NoError(t, s.Set("k", []byte("v1")))
	require.NoError(t, s.Set("k", []byte("v2")))

	got, ok, err := s.Get("k")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("v2"), got)
}

func TestMemStorage_DefensiveCopy(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	original := []byte("original")
	require.NoError(t, s.Set("k", original))

	// mutate the original slice after Set
	original[0] = 'X'

	got, ok, err := s.Get("k")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("original"), got, "Set should store a defensive copy")

	// mutate the returned slice
	got[0] = 'Y'

	got2, _, _ := s.Get("k")
	assert.Equal(t, []byte("original"), got2, "Get should return a defensive copy")
}

func TestMemStorage_ContainsKey(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	assert.False(t, s.ContainsKey("k"))

	require.NoError(t, s.Set("k", []byte("v")))
	assert.True(t, s.ContainsKey("k"))
	assert.False(t, s.ContainsKey("other"))
}

func TestMemStorage_Delete(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	require.NoError(t, s.Set("k", []byte("v")))
	assert.True(t, s.ContainsKey("k"))

	s.Delete("k")
	assert.False(t, s.ContainsKey("k"))

	got, ok, err := s.Get("k")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestMemStorage_DeleteMissing(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	// should not panic
	s.Delete("nonexistent")
}

func TestMemStorage_KeySet(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	require.NoError(t, s.Set("b", []byte("2")))
	require.NoError(t, s.Set("a", []byte("1")))
	require.NoError(t, s.Set("c", []byte("3")))

	keys := s.KeySet()
	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestMemStorage_KeySetEmpty(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	keys := s.KeySet()
	assert.Empty(t, keys)
}

func TestMemStorage_DeleteAll(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	require.NoError(t, s.Set("a", []byte("1")))
	require.NoError(t, s.Set("b", []byte("2")))

	s.DeleteAll()

	assert.Empty(t, s.KeySet())
	assert.False(t, s.ContainsKey("a"))
	assert.False(t, s.ContainsKey("b"))
}

func TestMemStorage_Size(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	assert.Equal(t, 0, s.Size())

	require.NoError(t, s.Set("a", []byte("1")))
	assert.Equal(t, 1, s.Size())

	require.NoError(t, s.Set("b", []byte("2")))
	require.NoError(t, s.Set("c", []byte("3")))
	assert.Equal(t, 3, s.Size())

	// overwrite should not change size
	require.NoError(t, s.Set("a", []byte("updated")))
	assert.Equal(t, 3, s.Size())

	s.Delete("b")
	assert.Equal(t, 2, s.Size())

	s.DeleteAll()
	assert.Equal(t, 0, s.Size())
}

func TestMemStorage_Close(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	err := s.Close()
	assert.NoError(t, err)
}
