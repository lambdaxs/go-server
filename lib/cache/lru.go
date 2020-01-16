package cache

import (
    "github.com/bluele/gcache"
)

func InitLRUCacheWithExpire(size int, expireFunc gcache.LoaderExpireFunc) *gcache.Cache {
    cache := gcache.New(size).LRU().LoaderExpireFunc(expireFunc).Build()
    return &cache
}

func InitLRUCache(size int, expireFunc gcache.LoaderFunc) *gcache.Cache {
    cache := gcache.New(size).LRU().LoaderFunc(expireFunc).Build()
    return &cache
}
