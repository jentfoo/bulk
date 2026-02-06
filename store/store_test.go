package store

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyPrefixStorage_EmptyPrefix(t *testing.T) {
	t.Parallel()

	s := NewMemStore()
	ps := KeyPrefixStore(s, "")
	assert.Equal(t, s, ps, "empty prefix should return the original store")
}

func TestKeyPrefixStorage_SetGet(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	ps := KeyPrefixStore(underlying, "pfx")

	require.NoError(t, ps.Set("key", []byte("val")))

	// read through prefix storage
	got, ok, err := ps.Get("key")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("val"), got)

	// verify underlying key is prefixed with separator
	got2, ok2, err2 := underlying.Get("pfx;key")
	require.NoError(t, err2)
	assert.True(t, ok2)
	assert.Equal(t, []byte("val"), got2)

	// unprefixed key should not exist in underlying
	_, ok3, _ := underlying.Get("key")
	assert.False(t, ok3)
}

func TestKeyPrefixStorage_ContainsKey(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	ps := KeyPrefixStore(underlying, "ns")

	require.NoError(t, ps.Set("k", []byte("v")))
	assert.True(t, ps.ContainsKey("k"))
	assert.False(t, ps.ContainsKey("other"))

	// underlying should see the prefixed key
	assert.True(t, underlying.ContainsKey("ns;k"))
	assert.False(t, underlying.ContainsKey("k"))
}

func TestKeyPrefixStorage_Delete(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	ps := KeyPrefixStore(underlying, "ns")

	require.NoError(t, ps.Set("k", []byte("v")))
	ps.Delete("k")

	assert.False(t, ps.ContainsKey("k"))
	assert.False(t, underlying.ContainsKey("ns;k"))
}

func TestKeyPrefixStorage_KeySet(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	// add some keys outside the prefix
	require.NoError(t, underlying.Set("other;x", []byte("1")))
	require.NoError(t, underlying.Set("unrelated", []byte("2")))

	ps := KeyPrefixStore(underlying, "ns")
	require.NoError(t, ps.Set("a", []byte("3")))
	require.NoError(t, ps.Set("b", []byte("4")))

	keys := ps.KeySet()
	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b"}, keys)
}

func TestKeyPrefixStorage_KeySetEmpty(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	require.NoError(t, underlying.Set("other", []byte("v")))

	ps := KeyPrefixStore(underlying, "ns")
	assert.Empty(t, ps.KeySet())
}

func TestKeyPrefixStorage_DeleteAll(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	// add keys outside the prefix that should survive
	require.NoError(t, underlying.Set("keep", []byte("safe")))

	ps := KeyPrefixStore(underlying, "ns")
	require.NoError(t, ps.Set("a", []byte("1")))
	require.NoError(t, ps.Set("b", []byte("2")))

	ps.DeleteAll()

	assert.Empty(t, ps.KeySet())
	// underlying key outside prefix should remain
	assert.True(t, underlying.ContainsKey("keep"))
}

func TestKeyPrefixStorage_Size(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	// add keys outside the prefix
	require.NoError(t, underlying.Set("other", []byte("1")))

	ps := KeyPrefixStore(underlying, "ns")
	assert.Equal(t, 0, ps.Size())

	require.NoError(t, ps.Set("a", []byte("2")))
	assert.Equal(t, 1, ps.Size())

	require.NoError(t, ps.Set("b", []byte("3")))
	assert.Equal(t, 2, ps.Size())

	// underlying has 3 total, prefix only sees 2
	assert.Equal(t, 3, underlying.Size())

	ps.Delete("a")
	assert.Equal(t, 1, ps.Size())

	ps.DeleteAll()
	assert.Equal(t, 0, ps.Size())

	// underlying still has the non-prefixed key
	assert.Equal(t, 1, underlying.Size())
}

func TestKeyPrefixStorage_Close(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	require.NoError(t, underlying.Set("root", []byte("keep")))

	ps := KeyPrefixStore(underlying, "ns")
	require.NoError(t, ps.Set("a", []byte("1")))
	require.NoError(t, ps.Set("b", []byte("2")))

	err := ps.Close()
	require.NoError(t, err)

	// prefixed keys should be removed
	assert.Empty(t, ps.KeySet())

	// underlying store should still be usable and retain non-prefixed keys
	assert.True(t, underlying.ContainsKey("root"))
	require.NoError(t, underlying.Set("new", []byte("works")))
	got, ok, err := underlying.Get("new")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("works"), got)
}

func TestKeyPrefixStorage_Isolation(t *testing.T) {
	t.Parallel()

	underlying := NewMemStore()
	a := KeyPrefixStore(underlying, "a")
	b := KeyPrefixStore(underlying, "b")

	require.NoError(t, a.Set("key", []byte("from_a")))
	require.NoError(t, b.Set("key", []byte("from_b")))

	gotA, okA, errA := a.Get("key")
	require.NoError(t, errA)
	assert.True(t, okA)
	assert.Equal(t, []byte("from_a"), gotA)

	gotB, okB, errB := b.Get("key")
	require.NoError(t, errB)
	assert.True(t, okB)
	assert.Equal(t, []byte("from_b"), gotB)

	// deleting from a should not affect b
	a.Delete("key")
	assert.False(t, a.ContainsKey("key"))
	assert.True(t, b.ContainsKey("key"))

	// DeleteAll on b should not affect underlying non-prefixed keys
	require.NoError(t, underlying.Set("root", []byte("r")))
	b.DeleteAll()
	assert.Empty(t, b.KeySet())
	assert.True(t, underlying.ContainsKey("root"))
}
