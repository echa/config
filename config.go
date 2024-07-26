// Copyright (c) 2018-2024 KIDTSUNAMI
// Author: alex@kidtsunami.com
//

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var config = NewConfig() // ReadConfig(), Set()

func ConfigName() string {
	return config.ConfigName()
}

func SetConfigName(name string) *Config {
	return config.SetConfigName(name)
}

func UseEnv(enabled bool) *Config {
	return config.UseEnv(enabled)
}

func SetEnvPrefix(p string) *Config {
	return config.SetEnvPrefix(p)
}

func ReadConfigFile() error {
	return config.ReadConfigFile()
}

func MustReadConfigFile() error {
	return config.MustReadConfigFile()
}

func ReadConfig(buf []byte) error {
	return config.ReadConfig(buf)
}

func Set(key string, val any) *Config {
	return config.Set(key, val)
}

func SetDefault(key string, val any) *Config {
	return config.SetDefault(key, val)
}

func GetString(path string) string {
	return config.GetString(path)
}

func GetStringSlice(path string) []string {
	return config.GetStringSlice(path)
}

func GetStringMap(path string) map[string]string {
	return config.GetStringMap(path)
}

func GetInterface(path string) any {
	return config.GetInterface(path)
}

func GetDuration(path string) time.Duration {
	return config.GetDuration(path)
}

func GetTime(path string) time.Time {
	return config.GetTime(path)
}

func GetBool(path string) bool {
	return config.GetBool(path)
}

func GetInt64(path string) int64 {
	return config.GetInt64(path)
}

func GetUint64(path string) uint64 {
	return config.GetUint64(path)
}

func GetInt64Slice(path string) []int64 {
	return config.GetInt64Slice(path)
}

func GetUint64Slice(path string) []uint64 {
	return config.GetUint64Slice(path)
}

func GetInt(path string) int {
	return config.GetInt(path)
}

func GetUint(path string) uint {
	return config.GetUint(path)
}

func GetIntSlice(path string) []int {
	return config.GetIntSlice(path)
}

func GetUintSlice(path string) []uint {
	return config.GetUintSlice(path)
}

func GetFloat64(path string) float64 {
	return config.GetFloat64(path)
}

func GetFloat64Slice(path string) []float64 {
	return config.GetFloat64Slice(path)
}

func All() map[string]any {
	return config.All()
}

func Unmarshal(path string, val any) error {
	return config.Unmarshal(path, val)
}

func ForEach(path string, fn func(c *Config) error) error {
	return config.ForEach(path, fn)
}

func Expand(s string) string {
	return config.Expand(s)
}

type Config struct {
	confName  string
	envPrefix string
	noEnv     bool
	data      map[string]any // read from config file or set
	merged    map[string]any // merged env, data, defaults
	defaults  map[string]any // flat 1-level key/value pairs
}

func NewConfig() *Config {
	return &Config{
		data:     make(map[string]any),
		defaults: make(map[string]any),
		merged:   nil,
	}
}

func canAccess(name string) bool {
	s, err := os.Stat(name)
	return err == nil && !s.IsDir()
}

func (c *Config) ConfigName() string {
	name := c.confName
	if name == "" || !canAccess(name) {
		name = os.Getenv(c.expandEnvKey("CONFIG_FILE"))
	}
	if name == "" || !canAccess(name) {
		name = "config.json"
	}
	return name
}

func (c *Config) SetConfigName(name string) *Config {
	c.confName = name
	return c
}

func (c *Config) SetEnvPrefix(p string) *Config {
	c.envPrefix = strings.ToUpper(strings.Replace(p, " ", "_", -1))
	c.merged = nil
	return c
}

func (c *Config) EnvPrefix() string {
	return c.envPrefix
}

func (c *Config) UseEnv(enabled bool) *Config {
	c.noEnv = !enabled
	return c
}

func (c *Config) MustReadConfigFile() error {
	return c.ReadConfigFile(true)
}

func (c *Config) ReadConfigFile(failNonExist ...bool) error {
	// determine config name from
	// - local variable
	// - environment
	// - fallback: use config.json
	name := c.ConfigName()

	// be resilient to non existent config file
	_, err := os.Stat(name)
	if err != nil {
		if len(failNonExist) > 0 && failNonExist[0] {
			return err
		}
		return nil
	}

	// read config file
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Errorf("reading config file: %v", err)
	}
	return c.ReadConfig(buf)
}

func (c *Config) ReadConfig(buf []byte) error {
	// unpack config from JSON into Go map
	if err := json.Unmarshal(buf, &c.data); err != nil {
		return fmt.Errorf("parsing config file: %v", err)
	}
	c.merged = nil
	// parse env for any defined value
	_ = c.All()
	return nil
}

func (c *Config) expandEnvKey(key string) string {
	key = strings.ToUpper(key)
	key = strings.Replace(key, ".", "_", -1)
	if c.envPrefix != "" {
		return c.envPrefix + "_" + key
	}
	return key
}

func (c *Config) Set(key string, val any) *Config {
	setTree(c.data, key, val)
	c.merged = nil
	return c
}

func (c *Config) Use(val map[string]any) *Config {
	c.data = val
	c.merged = nil
	return c
}

func (c *Config) SetDefault(key string, val any) *Config {
	c.defaults[key] = val // flat
	c.merged = nil
	return c
}

func setTree(walker map[string]any, key string, val any) {
	keys := strings.Split(key, ".")
	for n, v := range keys {
		if sub, ok := walker[v]; ok {
			if submap, ok := sub.(map[string]any); ok {
				walker = submap
			} else if n == len(keys)-1 {
				walker[v] = val
			} else {
				log.Fatalf("config: cannot set path '%s': %s exists as type %T", key, v, sub)
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
	for n, v := range keys {
		if sub, ok := walker[v]; ok {
			if submap, ok := sub.(map[string]any); ok {
				walker = submap
			} else if n == len(keys)-1 {
				if _, ok := walker[v]; !ok {
					walker[v] = val
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
	for n, v := range keys {
		if sub, ok := walker[v]; ok {
			if n == len(keys)-1 {
				return sub
			}
			if submap, ok := sub.(map[string]any); ok {
				walker = submap
			} else {
				break
			}
		}
	}
	return nil
}

func (c *Config) getEnv(path string) (string, bool) {
	if c.noEnv {
		return "", false
	}
	path = strings.ToUpper(path)
	path = strings.Replace(path, ".", "_", -1)
	return os.LookupEnv(c.expandEnvKey(path))
}

func (c *Config) getValue(path string) any {
	// get env key when present (allows to overwrite with empty value)
	if val, ok := c.getEnv(path); ok {
		return val
	}
	// try get merged tree (env + config + defaults)
	if c.merged != nil {
		if val := getTree(c.merged, path); val != nil {
			return val
		}
	}
	// get config file data if set
	if c.data != nil {
		if val := getTree(c.data, path); val != nil {
			return val
		}
	}
	// get default if registered
	if val, ok := c.defaults[path]; ok {
		return val
	}
	return nil
}

func (c *Config) GetString(path string) string {
	val := c.getValue(path)
	if val == nil {
		return ""
	}
	return toString(val)
}

func (c *Config) GetStringSlice(path string) []string {
	val := c.getValue(path)
	if val == nil {
		return []string{}
	}
	switch s := val.(type) {
	case []string:
		return s
	case []any:
		res := make([]string, len(s))
		for i, v := range s {
			res[i] = toString(v)
		}
		return res
	default:
		return strings.Split(toString(val), ",")
	}
}

func (c *Config) GetStringMap(path string) map[string]string {
	val := c.getValue(path)
	smap := make(map[string]string)
	if val == nil {
		return smap
	}
	switch m := val.(type) {
	case map[string]string:
		smap = m
	case map[string]any:
		for k, v := range m {
			k = strings.ToLower(k)
			if s := toString(v); s != "" {
				smap[k] = s
			}
		}
	case []string:
		for _, v := range m {
			k, v, ok := strings.Cut(v, "=")
			k = strings.ToLower(k)
			if ok {
				smap[k] = v
			} else {
				smap[k] = "true"
			}
		}
	case string:
		for _, v := range strings.Split(m, ",") {
			k, v, ok := strings.Cut(v, "=")
			k = strings.ToLower(k)
			if ok {
				smap[k] = v
			} else {
				smap[k] = "true"
			}
		}
	}
	// check if value was overwritten by env
	for n := range smap {
		if ev, ok := c.getEnv(path + "." + n); ok {
			smap[n] = ev
		}
	}
	// add extra values from env
	pfx := c.expandEnvKey(path)
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, pfx) {
			continue
		}
		k, v, _ := strings.Cut(v, "=")
		k = strings.TrimPrefix(k, pfx+"_")
		k = strings.ToLower(k)
		smap[k] = v
	}
	return smap
}

func (c *Config) GetInterface(path string) any {
	return c.getValue(path)
}

func (c *Config) GetDuration(path string) time.Duration {
	val := c.getValue(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case time.Duration:
		return v
	case int:
		return time.Duration(v)
	case int32:
		return time.Duration(v)
	case uint32:
		return time.Duration(v)
	case int64:
		return time.Duration(v)
	case uint64:
		return time.Duration(v)
	case float64:
		return time.Duration(int64(v))
	case string:
		if dur, err := ParseDuration(v); err == nil {
			return dur.Duration()
		}
	default:
		if dur, err := ParseDuration(toString(v)); err == nil {
			return dur.Duration()
		}
	}
	return 0
}

func (c *Config) GetTime(path string) time.Time {
	val := c.getValue(path)
	if val == nil {
		return time.Time{}
	}
	switch v := val.(type) {
	case time.Time:
		return v
	case int:
		return time.Unix(int64(v), 0)
	case int32:
		return time.Unix(int64(v), 0)
	case uint32:
		return time.Unix(int64(v), 0)
	case int64:
		return time.Unix(int64(v), 0)
	case uint64:
		return time.Unix(int64(v), 0)
	case float64:
		return time.Unix(int64(v), 0)
	case string:
		if tm, err := ParseTime(v); err == nil {
			return tm
		}
	default:
		if tm, err := ParseTime(toString(v)); err == nil {
			return tm
		}
	}
	return time.Time{}
}

func (c *Config) GetBool(path string) bool {
	val := c.getValue(path)
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	default:
		s := toString(v)
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
	}
	return false
}

func (c *Config) GetInt(path string) int {
	val := c.getValue(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return int(n)
		}
	default:
		s := toString(v)
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return int(i)
		}
	}
	return 0
}

func (c *Config) GetUint(path string) uint {
	val := c.getValue(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint64:
		return uint(v)
	case float64:
		return uint(v)
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return uint(n)
		}
	default:
		s := toString(v)
		if i, err := strconv.ParseUint(s, 10, 64); err == nil {
			return uint(i)
		}
	}
	return 0
}

func (c *Config) GetIntSlice(path string) []int {
	is := make([]int, 0)
	for _, v := range c.GetStringSlice(path) {
		val, _ := strconv.Atoi(v)
		is = append(is, val)
	}
	return is
}

func (c *Config) GetUintSlice(path string) []uint {
	is := make([]uint, 0)
	for _, v := range c.GetStringSlice(path) {
		val, _ := strconv.Atoi(v)
		is = append(is, uint(val))
	}
	return is
}

func (c *Config) GetInt64(path string) int64 {
	val := c.getValue(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return n
		}
	default:
		s := toString(v)
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func (c *Config) GetUint64(path string) uint64 {
	val := c.getValue(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint64:
		return v
	case float64:
		return uint64(v)
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return uint64(n)
		}
	default:
		s := toString(v)
		if i, err := strconv.ParseUint(s, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func (c *Config) GetInt64Slice(path string) []int64 {
	is := make([]int64, 0)
	for _, v := range c.GetStringSlice(path) {
		val, _ := strconv.ParseInt(v, 10, 64)
		is = append(is, val)
	}
	return is
}

func (c *Config) GetUint64Slice(path string) []uint64 {
	is := make([]uint64, 0)
	for _, v := range c.GetStringSlice(path) {
		val, _ := strconv.ParseUint(v, 10, 64)
		is = append(is, val)
	}
	return is
}

func (c *Config) GetFloat64(path string) float64 {
	val := c.getValue(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case json.Number:
		if n, err := v.Float64(); err == nil {
			return n
		}
	default:
		s := toString(v)
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return 0
}

func (c *Config) GetFloat64Slice(path string) []float64 {
	is := make([]float64, 0)
	for _, v := range c.GetStringSlice(path) {
		val, _ := strconv.ParseFloat(v, 64)
		is = append(is, val)
	}
	return is
}

func (c *Config) All() map[string]any {
	if c.merged != nil {
		return c.merged
	}
	c.merged = make(map[string]any)

	// load data map into merged
	buf, _ := json.Marshal(&c.data)
	json.Unmarshal(buf, &c.merged)

	// add defaults for missing (nested) keys
	for key, val := range c.defaults {
		setTreeIfEmpty(c.merged, key, val)
	}

	// extend keys with matching env variables, only if env prefix is set
	if c.noEnv || c.envPrefix == "" {
		return c.merged
	}

	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, c.envPrefix) {
			continue
		}
		fields := strings.SplitN(v, "=", 2)
		key := strings.Join(strings.SplitN(strings.ToLower(fields[0]), "_", 3)[1:], ".")
		setTree(c.merged, key, fields[1])
	}
	return c.merged
}

func (c *Config) ForEach(path string, fn func(c *Config) error) error {
	// requires merged tree
	s := c.All()
	segs := strings.Split(path, ".")
	var slice []any
	for i, v := range segs[:len(segs)-1] {
		if s == nil {
			break
		}
		sub, ok := s[v]
		if !ok {
			return fmt.Errorf("missing config path '%s'", path)
		}
		s, ok = sub.(map[string]any)
		if !ok && i < len(segs)-1 {
			return fmt.Errorf("invalid type %T at config path '%s'", sub, path)
		}
	}
	// assuming the last sub-tree element is a slice
	slice, ok := s[segs[len(segs)-1]].([]any)
	if !ok {
		return fmt.Errorf("expected slice of values at path '%s'", path)
	}
	for i, v := range slice {
		err := fn(&Config{
			envPrefix: c.expandEnvKey(path + "." + strconv.Itoa(i)),
			noEnv:     c.noEnv,
			data:      v.(map[string]any),
			merged:    v.(map[string]any),
		})
		if err != nil {
			return err
		}
	}
	if c.noEnv {
		return nil
	}
	// search for more env keys starting with `$path_$num`
	more := len(slice)
	for {
		found := false
		prefix := c.expandEnvKey(path + "." + strconv.Itoa(more))
		for _, v := range os.Environ() {
			if !strings.HasPrefix(v, prefix) {
				continue
			}
			err := fn(&Config{
				envPrefix: prefix,
				data:      nil,
				merged:    nil,
			})
			if err != nil {
				return err
			}
			found = true
			more++
			break
		}
		if !found {
			break
		}
	}
	return nil
}

func (c *Config) Unmarshal(path string, val any) error {
	// requires merged tree
	s := c.All()
	for _, v := range strings.Split(path, ".") {
		if s == nil {
			break
		}
		sub, ok := s[v]
		if !ok {
			return fmt.Errorf("missing config path '%s'", path)
		}
		s, ok = sub.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid type %T at config path '%s'", sub, path)
		}
	}
	buf, _ := json.Marshal(s)
	return json.Unmarshal(buf, val)
}

var exp = regexp.MustCompile(`\${(.+?)}`)

// Resolves env/config variables embedded in a string using ${VAR}
func (c *Config) Expand(s string) string {
	for _, match := range exp.FindAllStringSubmatch(s, -1) {
		if len(match) < 2 {
			continue
		}
		path := strings.ToLower(match[1])
		path = strings.Replace(path, "_", ".", -1)
		val := c.getValue(path)
		if val == nil {
			continue
		}
		s = strings.Replace(s, match[0], toString(val), -1)
	}
	return s
}
