package badger

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/pkg/errors"
	"github.com/brezzgg/sub/internal/repo/repo"
	"github.com/dgraph-io/badger/v2"
)

type Badger struct {
	db *badger.DB
}

func Constructor(options map[string]any) (repo.StorageProvider, error) {
	file, err := repo.GetOption[string](options, "file")
	if err != nil {
		return nil, err
	}
	if *file == "" {
		*file = ".db"
	}

	opts := badger.DefaultOptions(*file)
	opts.Logger = &Logger{}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	return &Badger{db: db}, nil
}

// Get implements [repo.StorageProvider].
func (b *Badger) Get(id string) (*entity.Subscription, error) {
	var buf []byte
	err := b.db.View(func(txn *badger.Txn) error {
		// get item
		item, err := txn.Get([]byte(id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return errors.Code(errors.CodeNotFound)
			}
			return errors.Errorf(errors.CodeInternal, "get: %s", err)
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return errors.Errorf(errors.CodeInternal, "value copy: %s", err)
		}
		buf = val

		return nil
	})
	if err != nil {
		return nil, err
	}

	var res entity.Subscription
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	if err := dec.Decode(&res); err != nil {
		return nil, errors.Errorf(errors.CodeMarshalError, "decode: %s", err)
	}

	return &res, nil
}

// Remove implements [repo.StorageProvider].
func (b *Badger) Remove(id string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(id)); err != nil {
			return errors.Errorf(errors.CodeInternal, "delete: %s", err)
		}
		return nil
	})
}

// Set implements [repo.StorageProvider].
func (b *Badger) Set(id string, s *entity.Subscription) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		return errors.Errorf(errors.CodeMarshalError, "encode: %s", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(id), buf.Bytes()); err != nil {
			return errors.Errorf(errors.CodeInternal, "update: %s", err)
		}
		return nil
	})
}

// WhereFn implements [repo.StorageProvider].
func (b *Badger) WhereFn(fn repo.WhereFunc) (map[string]*entity.Subscription, error) {
	res := make(map[string]*entity.Subscription)
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())

			val, err := item.ValueCopy(nil)
			if err != nil {
				lg.Error("internal", err)
				continue
			}

			var sub entity.Subscription
			if err := gob.NewDecoder(bytes.NewReader(val)).Decode(&sub); err != nil {
				lg.Error("internal", err)
				continue
			}

			if fn != nil && !fn(&sub) {
				continue
			}

			res[key] = &sub
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// WhereMeta implements [repo.StorageProvider].
func (b *Badger) WhereMeta(equals map[string]any) (map[string]*entity.Subscription, error) {
	res := make(map[string]*entity.Subscription)
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())

			val, err := item.ValueCopy(nil)
			if err != nil {
				lg.Error("internal", err)
				continue
			}

			var sub entity.Subscription
			if err := gob.NewDecoder(bytes.NewReader(val)).Decode(&sub); err != nil {
				lg.Error("internal", err)
				continue
			}

			if !repo.MetaEqual(sub.Metadata, equals) {
				continue
			}

			res[key] = &sub
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

var _ repo.StorageProvider = (*Badger)(nil)
