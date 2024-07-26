// Copyright (c) 2018-2024 KIDTSUNAMI
// Author: alex@kidtsunami.com
//

package config

import (
	"strconv"
	"strings"
)

func setTree(walker map[string]any, key string, val any) {
	keys := strings.Split(key, ".")
	for n := 0; n < len(keys); n++ {
		v := keys[n]
		if sub, ok := walker[v]; ok {
			// recurse into subtree
			switch e := sub.(type) {
			case map[string]any:
				walker = e
			case []any:
				i, _ := strconv.ParseInt(v, 10, 64)
				walker = e[int(i)].(map[string]any)
			default:
				// append leaf if type is not a container
				walker[v] = val
			}
		} else if n < len(keys)-1 {
			// append subtree
			sub := make(map[string]any)
			walker[v] = sub
			walker = sub
		} else {
			// append leaf
			walker[v] = val
		}
	}
}

func setTreeIfEmpty(walker map[string]any, key string, val any) {
	keys := strings.Split(key, ".")
	for n := 0; n < len(keys); n++ {
		v := keys[n]
		if sub, ok := walker[v]; ok {
			// recurse into subtree
			switch e := sub.(type) {
			case map[string]any:
				walker = e
			case []any:
				i, _ := strconv.ParseInt(v, 10, 64)
				walker = e[int(i)].(map[string]any)
			default:
				// append leaf if type is not a container and key segment is last
				if n == len(keys)-1 {
					if _, ok := walker[v]; !ok {
						walker[v] = val
					}
				}
			}
		} else if n < len(keys)-1 {
			// append subtree
			sub := make(map[string]any)
			walker[v] = sub
			walker = sub
		} else {
			// append leaf
			walker[v] = val
		}
	}
}

func getTree(walker map[string]any, key string) any {
	keys := strings.Split(key, ".")
	for n := 0; n < len(keys); n++ {
		v := keys[n]
		sub, ok := walker[v]
		if !ok {
			return nil
		}
		if n == len(keys)-1 {
			return sub
		}
		switch e := sub.(type) {
		case map[string]any:
			walker = e
		case []any:
			i, _ := strconv.ParseInt(v, 10, 64)
			walker = e[int(i)].(map[string]any)
		default:
			break
		}
	}
	return nil
}

func walkTree(tree map[string]any, prefix string, fn func(key, val string) error) (err error) {
	for n, v := range tree {
		key := n
		if prefix != "" {
			key = prefix + "." + key
		}
		switch sub := v.(type) {
		case map[string]any:
			err = walkTree(sub, key, fn)
		default:
			err = fn(key, toString(v))
		}
		if err != nil {
			break
		}
	}
	return
}
