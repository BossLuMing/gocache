1.由于业务需求需要用到缓存，目前网上查到的go相关的缓存框架不太适用我的业务需求，所以个人使用map写了一个使用map相对简单和贴合业务需求的缓存功能。

2.目前有获取缓存，设置缓存，缓存数量，过期清理，更新已有缓存，但是还没有做对内存的管理，有相关需求的朋友可以自己补全。