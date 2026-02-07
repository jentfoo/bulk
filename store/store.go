package store

import (
	"errors"
	"strings"

	"github.com/go-analyze/bulk"
)

// ErrClosed is returned when an operation is attempted on a closed Storage.
var ErrClosed = errors.New("storage closed")

// Storage defines persistence methods for key-value blob storage.
type Storage interface {
	// Set stores blob under key, overwriting any previous value.
	Set(key string, blob []byte) error
	// Get retrieves the blob for key. Returns (nil, false, nil) when the key
	// does not exist and (blob, true, nil) on success.
	Get(key string) ([]byte, bool, error)
	// Delete removes a single key. No-op if the key does not exist.
	Delete(key string) error
	// DeleteAll removes every key from the store.
	DeleteAll() error
	// ContainsKey reports whether key exists in the store.
	ContainsKey(key string) bool
	// KeySet returns all keys currently in the store.
	KeySet() []string
	// Size reports how many entries are stored in the map.
	Size() int
	// Close releases any resources held by the store.
	Close() error
}

// KeyPrefixStore wraps another Storage, prepending a fixed prefix
// (with a ";" separator) to every key. KeySet strips the prefix before
// returning, and DeleteAll only removes keys that belong to this prefix.
func KeyPrefixStore(s Storage, prefix string) Storage {
	if prefix == "" {
		return s
	}
	return &prefixStorage{
		store:  s,
		prefix: prefix + ";",
	}
}

type prefixStorage struct {
	store  Storage
	prefix string
}

func (p *prefixStorage) Set(key string, blob []byte) error {
	return p.store.Set(p.prefix+key, blob)
}

func (p *prefixStorage) Get(key string) ([]byte, bool, error) {
	return p.store.Get(p.prefix + key)
}

func (p *prefixStorage) Delete(key string) error {
	return p.store.Delete(p.prefix + key)
}

func (p *prefixStorage) DeleteAll() error {
	var errs []error
	for _, k := range p.store.KeySet() {
		if strings.HasPrefix(k, p.prefix) {
			if err := p.store.Delete(k); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (p *prefixStorage) ContainsKey(key string) bool {
	return p.store.ContainsKey(p.prefix + key)
}

func (p *prefixStorage) KeySet() []string {
	return bulk.SliceFilter(func(k string) bool {
		return strings.HasPrefix(k, p.prefix)
	}, p.store.KeySet())
}

func (p *prefixStorage) Size() int {
	var result int
	for _, k := range p.store.KeySet() {
		if strings.HasPrefix(k, p.prefix) {
			result++
		}
	}
	return result
}

// Close is a no-op. The caller that created the underlying Storage is responsible for closing it.
func (p *prefixStorage) Close() error {
	return nil
}
