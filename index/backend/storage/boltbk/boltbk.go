/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package boltbk implements BoltDB backed backend for standalone instances
// and test mode, you should use Etcd in production
package boltbk

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "sync"
    "time"

    "github.com/boltdb/bolt"
    "github.com/mailgun/timetools"
    "gopkg.in/vmihailenco/msgpack.v2"
)

// --- not found error ---
func notFound(format string, a ...interface{}) error {
    return &bkNotFoundError{fmt.Sprintf(format, a...)}
}
type bkNotFoundError struct {s string}
func (b *bkNotFoundError) Error() string {return b.s}
func IsNotFound(err error) bool {
    _, ok := err.(*bkNotFoundError)
    return ok
}

// --- already exists ---
func alreadyExists(format string, a ...interface{}) error {
    return &bkAlreadyExistsError{fmt.Sprintf(format, a...)}
}
type bkAlreadyExistsError struct {s string}
func (b *bkAlreadyExistsError) Error() string {return b.s}
func IsAlreadyExists(err error) bool {
    _, ok := err.(*bkAlreadyExistsError)
    return ok
}

// --- bad parameter ---
func badParameter(format string, a ...interface{}) error {
    return &bkBadParameterError{fmt.Sprintf(format, a...)}
}
type bkBadParameterError struct {s string}
func (b *bkBadParameterError) Error() string {return b.s}
func IsBadParameter(err error) bool {
    _, ok := err.(*bkBadParameterError)
    return ok
}

// --- compare failed ---
func compareFailed(format string, a ...interface{}) error {
    return &bkcompareFailedError{fmt.Sprintf(format, a...)}
}
type bkcompareFailedError struct {s string}
func (b *bkcompareFailedError) Error() string {return b.s}
func IsCompareFailed(err error) bool {
    _, ok := err.(*bkcompareFailedError)
    return ok
}

// BoltBackend is a boltdb-based backend used in tests and standalone mode
type BoltBackend struct {
    sync.Mutex

    db    *bolt.DB
    clock timetools.TimeProvider
    locks map[string]time.Time
}

// Option sets functional options for the backend
type Option func(b *BoltBackend) error

// Clock sets clock for the backend, used in tests
func Clock(clock timetools.TimeProvider) Option {
    return func(b *BoltBackend) error {
        b.clock = clock
        return nil
    }
}

// New returns a new isntance of bolt backend
func New(path string, opts ...Option) (*BoltBackend, error) {
    path, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf(err.Error() + " : failed to convert path")
    }
    dir := filepath.Dir(path)
    s, err := os.Stat(dir)
    if err != nil {
        return nil, err
    }
    if !s.IsDir() {
        return nil, badParameter("path '%v' should be a valid directory", dir)
    }
    b := &BoltBackend{
        locks: make(map[string]time.Time),
    }
    for _, option := range opts {
        if err := option(b); err != nil {
            return nil, err
        }
    }
    if b.clock == nil {
        b.clock = &timetools.RealTime{}
    }
    db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
    if err != nil {
        if err == bolt.ErrTimeout {
            return nil, fmt.Errorf("Local storage is locked. Another instance is running? (%v)", path)
        }
        return nil, err
    }
    b.db = db
    return b, nil
}

// Close closes the backend resources
func (b *BoltBackend) Close() error {
    return b.db.Close()
}

func (b *BoltBackend) GetKeys(path []string) ([]string, error) {
    keys, err := b.getKeys(path)
    if err != nil {
        if IsNotFound(err) {
            return nil, nil
        }
        return nil, err
    }
    // now do an iteration to expire keys
    for _, key := range keys {
        b.GetVal(path, key)
    }
    keys, err = b.getKeys(path)
    if err != nil {
        return nil, err
    }
    sort.Sort(sort.StringSlice(keys))
    return keys, nil
}

func (b *BoltBackend) UpsertObj(path []string, key string, val interface{}, ttl time.Duration) error {
    v, err := msgpack.Marshal(val)
    if err != nil {
        return err
    }
    return b.upsertVal(path, key, v, ttl)
}

func (b *BoltBackend) UpsertVal(path []string, key string, val []byte, ttl time.Duration) error {
    return b.upsertVal(path, key, val, ttl)
}

func (b *BoltBackend) CreateVal(bucket []string, key string, val []byte, ttl time.Duration) error {
    v := &kv{
        Created: b.clock.UtcNow(),
        Value:   val,
        TTL:     ttl,
    }
    bytes, err := msgpack.Marshal(v)
    if err != nil {
        return err
    }
    err = b.createKey(bucket, key, bytes)
    return err
}

func (b *BoltBackend) TouchVal(bucket []string, key string, ttl time.Duration) error {
    err := b.db.Update(func(tx *bolt.Tx) error {
        bkt, err := UpsertBucket(tx, bucket)
        if err != nil {
            return err
        }
        val := bkt.Get([]byte(key))
        if val == nil {
            return notFound("'%v' already exists", key)
        }
        var k *kv
        if err := msgpack.Unmarshal(val, &k); err != nil {
            return err
        }
        k.TTL = ttl
        k.Created = b.clock.UtcNow()
        bytes, err := msgpack.Marshal(k)
        if err != nil {
            return err
        }
        return bkt.Put([]byte(key), bytes)
    })
    return err
}

func (b *BoltBackend) upsertVal(path []string, key string, val []byte, ttl time.Duration) error {
    v := &kv{
        Created: b.clock.UtcNow(),
        Value:   val,
        TTL:     ttl,
    }
    bytes, err := msgpack.Marshal(v)
    if err != nil {
        return err
    }
    return b.upsertKey(path, key, bytes)
}

func (b *BoltBackend) CompareAndSwap(path []string, key string, val []byte, ttl time.Duration, prevVal []byte) ([]byte, error) {
    b.Lock()
    defer b.Unlock()

    storedVal, err := b.GetVal(path, key)
    if err != nil {
        if IsNotFound(err) && len(prevVal) != 0 {
            return nil, err
        }
    }
    if len(prevVal) == 0 && err == nil {
        return nil, alreadyExists("key '%v' already exists", key)
    }
    if string(prevVal) == string(storedVal) {
        err = b.upsertVal(path, key, val, ttl)
        if err != nil {
            return nil, err
        }
        return storedVal, nil
    }
    return storedVal, compareFailed("expected: %v, got: %v", string(prevVal), string(storedVal))
}

func (b *BoltBackend) GetVal(path []string, key string) ([]byte, error) {
    var val []byte
    if err := b.getKey(path, key, &val); err != nil {
        return nil, err
    }
    var k *kv
    if err := msgpack.Unmarshal(val, &k); err != nil {
        return nil, err
    }
    if k.TTL != 0 && b.clock.UtcNow().Sub(k.Created) > k.TTL {
        if err := b.deleteKey(path, key); err != nil {
            return nil, err
        }
        return nil, notFound("%v: %v not found", path, key)
    }
    return k.Value, nil
}

func (b *BoltBackend) GetObj(path []string, key string, obj interface{}) error {
    var val []byte
    if err := b.getKey(path, key, &val); err != nil {
        return err
    }
    var k *kv
    if err := msgpack.Unmarshal(val, &k); err != nil {
        return err
    }
    if k.TTL != 0 && b.clock.UtcNow().Sub(k.Created) > k.TTL {
        if err := b.deleteKey(path, key); err != nil {
            return err
        }
        return notFound("%v: %v not found", path, key)
    }
    if err := msgpack.Unmarshal(k.Value, obj); err != nil {
        return err
    }
    return nil
}

func (b *BoltBackend) GetValAndTTL(path []string, key string) ([]byte, time.Duration, error) {
    var val []byte
    if err := b.getKey(path, key, &val); err != nil {
        return nil, 0, err
    }
    var k *kv
    if err := msgpack.Unmarshal(val, &k); err != nil {
        return nil, 0, err
    }
    if k.TTL != 0 && b.clock.UtcNow().Sub(k.Created) > k.TTL {
        if err := b.deleteKey(path, key); err != nil {
            return nil, 0, err
        }
        return nil, 0, notFound("%v: %v not found", path, key)
    }
    var newTTL time.Duration
    newTTL = 0
    if k.TTL != 0 {
        newTTL = k.Created.Add(k.TTL).Sub(b.clock.UtcNow())
    }
    return k.Value, newTTL, nil
}

func (b *BoltBackend) GetObjAndTTL(path []string, key string, obj interface{}) (time.Duration, error) {
    var val []byte
    if err := b.getKey(path, key, &val); err != nil {
        return 0, err
    }
    var k *kv
    if err := msgpack.Unmarshal(val, &k); err != nil {
        return 0, err
    }
    if k.TTL != 0 && b.clock.UtcNow().Sub(k.Created) > k.TTL {
        if err := b.deleteKey(path, key); err != nil {
            return 0, err
        }
        return 0, notFound("%v: %v not found", path, key)
    }
    var newTTL time.Duration
    newTTL = 0
    if k.TTL != 0 {
        newTTL = k.Created.Add(k.TTL).Sub(b.clock.UtcNow())
    }
    if err := msgpack.Unmarshal(k.Value, obj); err != nil {
        return 0, err
    }
    return newTTL, nil
}


func (b *BoltBackend) DeleteKey(path []string, key string) error {
    b.Lock()
    defer b.Unlock()
    return b.deleteKey(path, key)
}

func (b *BoltBackend) DeleteBucket(path []string, bucket string) error {
    b.Lock()
    defer b.Unlock()
    return b.deleteBucket(path, bucket)
}

func (b *BoltBackend) deleteBucket(buckets []string, bucket string) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        bkt, err := GetBucket(tx, buckets)
        if err != nil {
            return err
        }
        if bkt.Bucket([]byte(bucket)) == nil {
            return notFound("%v not found", bucket)
        }
        return bkt.DeleteBucket([]byte(bucket))
    })
}

func (b *BoltBackend) deleteKey(buckets []string, key string) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        bkt, err := GetBucket(tx, buckets)
        if err != nil {
            return err
        }
        if bkt.Get([]byte(key)) == nil {
            return notFound("%v is not found", key)
        }
        return bkt.Delete([]byte(key))
    })
}

func (b *BoltBackend) upsertKey(buckets []string, key string, bytes []byte) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        bkt, err := UpsertBucket(tx, buckets)
        if err != nil {
            return err
        }
        return bkt.Put([]byte(key), bytes)
    })
}

func (b *BoltBackend) createKey(buckets []string, key string, bytes []byte) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        bkt, err := UpsertBucket(tx, buckets)
        if err != nil {
            return err
        }
        val := bkt.Get([]byte(key))
        if val != nil {
            return alreadyExists("'%v' already exists", key)
        }
        return bkt.Put([]byte(key), bytes)
    })
}

func (b *BoltBackend) getKey(buckets []string, key string, val *[]byte) error {
    return b.db.View(func(tx *bolt.Tx) error {
        bkt, err := GetBucket(tx, buckets)
        if err != nil {
            return err
        }
        bytes := bkt.Get([]byte(key))
        if bytes == nil {
            _, err := GetBucket(tx, append(buckets, key))
            if err == nil {
                return badParameter("key '%v 'is a bucket", key)
            }
            return notFound("%v %v not found", buckets, key)
        }
        *val = make([]byte, len(bytes))
        copy(*val, bytes)
        return nil
    })
}

func (b *BoltBackend) getKeys(buckets []string) ([]string, error) {
    out := []string{}
    err := b.db.View(func(tx *bolt.Tx) error {
        bkt, err := GetBucket(tx, buckets)
        if err != nil {
            return err
        }
        c := bkt.Cursor()
        for k, _ := c.First(); k != nil; k, _ = c.Next() {
            out = append(out, string(k))
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return out, nil
}

func UpsertBucket(b *bolt.Tx, buckets []string) (*bolt.Bucket, error) {
    bkt, err := b.CreateBucketIfNotExists([]byte(buckets[0]))
    if err != nil {
        return nil, err
    }
    for _, key := range buckets[1:] {
        bkt, err = bkt.CreateBucketIfNotExists([]byte(key))
        if err != nil {
            return nil, err
        }
    }
    return bkt, nil
}

func GetBucket(b *bolt.Tx, buckets []string) (*bolt.Bucket, error) {
    bkt := b.Bucket([]byte(buckets[0]))
    if bkt == nil {
        return nil, notFound("bucket %v not found", buckets[0])
    }
    for _, key := range buckets[1:] {
        bkt = bkt.Bucket([]byte(key))
        if bkt == nil {
            return nil, notFound("bucket %v not found", key)
        }
    }
    return bkt, nil
}

func (b *BoltBackend) AcquireLock(token string, ttl time.Duration) error {
    for {
        b.Lock()
        expires, ok := b.locks[token]
        if ok && (expires.IsZero() || expires.After(b.clock.UtcNow())) {
            b.Unlock()
            b.clock.Sleep(100 * time.Millisecond)
        } else {
            if ttl == 0 {
                b.locks[token] = time.Time{}
            } else {
                b.locks[token] = b.clock.UtcNow().Add(ttl)
            }
            b.Unlock()
            return nil
        }
    }
}

func (b *BoltBackend) ReleaseLock(token string) error {
    b.Lock()
    defer b.Unlock()

    expires, ok := b.locks[token]
    if !ok || (!expires.IsZero() && expires.Before(b.clock.UtcNow())) {
        return notFound("lock %v is deleted or expired", token)
    }
    delete(b.locks, token)
    return nil
}

type kv struct {
    Created time.Time     `msgpack:"created"`
    TTL     time.Duration `msgpack:"ttl"`
    Value   []byte        `msgpack:"val"`
}
