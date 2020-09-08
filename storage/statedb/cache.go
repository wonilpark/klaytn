// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package statedb

import (
	"errors"
	"strings"
)

type TrieNodeCacheType string

// TrieNodeCacheConfig contains configuration values of all TrieNodeCache.
// The validity of each value is checked in "New" function of each TrieNodeCache construction.
type TrieNodeCacheConfig struct {
	CacheType          TrieNodeCacheType
	FastCacheSizeMB    int      // Memory allowance (MB) to use for caching trie nodes in fast cache
	RedisEndpoints     []string // Endpoints of redis cache
	RedisClusterEnable bool     // Enable cluster-enabled mode of redis cache
}

// TrieNodeCache interface the cache of stateDB
type TrieNodeCache interface {
	Set(k, v []byte)
	Get(k []byte) []byte
	Has(k []byte) ([]byte, bool)
}

const (
	// Available stateDB cache types
	LocalCache  TrieNodeCacheType = "LocalCache"
	RemoteCache                   = "RemoteCache"
	HybridCache                   = "HybridCache"
)

var (
	validTrieNodeCacheTypes  = []TrieNodeCacheType{LocalCache, RemoteCache, HybridCache}
	errNotSupportedCacheType = errors.New("not supported stateDB TrieNodeCache type")
)

func (cacheType TrieNodeCacheType) ToValid() TrieNodeCacheType {
	for _, validType := range validTrieNodeCacheTypes {
		if strings.ToLower(string(cacheType)) == strings.ToLower(string(validType)) {
			return validType
		}
	}
	return ""
}

// NewTrieNodeCache creates one type of any supported trie node caches.
// NOTE: It returns (nil, nil) when the cache type is LocalCache and its size is set to zero.
func NewTrieNodeCache(config TrieNodeCacheConfig) (TrieNodeCache, error) {
	switch config.CacheType {
	case LocalCache:
		return NewFastCache(config.FastCacheSizeMB), nil
	case RemoteCache:
		return NewRedisCache(config.RedisEndpoints, config.RedisClusterEnable)
	case HybridCache:
		return NewHybridCache(config)
	default:
	}
	return nil, errNotSupportedCacheType
}
