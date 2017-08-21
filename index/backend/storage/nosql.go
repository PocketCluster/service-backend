package storage

import (
    "time"
)

// Forever means that object TTL will not expire unless deleted
const Forever time.Duration = 0

// Backend implements abstraction over local or remote storage backend
//
// Storage is modeled after BoltDB:
//  * bucket is a slice []string{"a", "b"}
//  * buckets contain key value pairs
//
type Nosql interface {
    // GetKeys returns a list of keys for a given path
    GetKeys(bucket []string) ([]string, error)
    // CreateVal creates value with a given TTL and key in the bucket
    // if the value already exists, returns AlreadyExistsError
    CreateVal(bucket []string, key string, val []byte, ttl time.Duration) error
    // TouchVal updates the TTL of the key without changing the value
    TouchVal(bucket []string, key string, ttl time.Duration) error
    //UpsertVal updates or inserts object with a given TTL into a bucket
    // ForeverTTL for no TTL
    UpsertObj(bucket []string, key string, val interface{}, ttl time.Duration) error
    // UpsertVal updates or inserts value with a given TTL into a bucket
    // ForeverTTL for no TTL
    UpsertVal(bucket []string, key string, val []byte, ttl time.Duration) error
    // GetVal return a value for a given key in the bucket
    GetVal(path []string, key string) ([]byte, error)
    // GetVal return a value for a given key in the bucket
    GetObj(path []string, key string, obj interface{}) (error)
    // GetValAndTTL returns value and TTL for a key in bucket
    GetValAndTTL(bucket []string, key string) ([]byte, time.Duration, error)
    // GetValAndTTL returns value and TTL for a key in bucket
    GetObjAndTTL(bucket []string, key string, obj interface{}) (time.Duration, error)
    // DeleteKey deletes a key in a bucket
    DeleteKey(bucket []string, key string) error
    // DeleteBucket deletes the bucket by a given path
    DeleteBucket(path []string, bkt string) error
    // AcquireLock grabs a lock that will be released automatically in TTL
    AcquireLock(token string, ttl time.Duration) error
    // ReleaseLock forces lock release before TTL
    ReleaseLock(token string) error
    // CompareAndSwap implements compare ans swap operation for a key
    CompareAndSwap(bucket []string, key string, val []byte, ttl time.Duration, prevVal []byte) ([]byte, error)
    // Close releases the resources taken up by this backend
    Close() error
}
