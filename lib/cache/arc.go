package cache

import (
    "github.com/bluele/gcache"
)

func InitARCCacheWithExpire(size int, expireFunc gcache.LoaderExpireFunc) *gcache.Cache {
    cache := gcache.New(size).ARC().LoaderExpireFunc(expireFunc).Build()
    return &cache
}

func InitARCCache(size int, expireFunc gcache.LoaderFunc) *gcache.Cache {
    cache := gcache.New(size).ARC().LoaderFunc(expireFunc).Build()
    return &cache
}
