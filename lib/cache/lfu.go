package cache

import "github.com/bluele/gcache"

func InitLFUCacheWithExpire(size int, expireFunc gcache.LoaderExpireFunc) *gcache.Cache {
    cache := gcache.New(size).LFU().LoaderExpireFunc(expireFunc).Build()
    return &cache
}

func InitLFUCache(size int, expireFunc gcache.LoaderFunc) *gcache.Cache {
    cache := gcache.New(size).LFU().LoaderFunc(expireFunc).Build()
    return &cache
}
