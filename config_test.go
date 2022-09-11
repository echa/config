// Copyright (c) 2018-2021 KIDTSUNAMI
// Author: alex@kidtsunami.com
//
// go test -cpuprofile cpu.prof -memprofile mem.prof -bench=. -benchmem

package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestName(T *testing.T) {
	c := NewConfig()
	if exp, got := "config.json", c.ConfigName(); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	f1, err := ioutil.TempFile("", "testconfig.json")
	if err != nil {
		T.Error(err)
	}
	defer os.RemoveAll(f1.Name())
	c.SetConfigName(f1.Name())
	if exp, got := f1.Name(), c.ConfigName(); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	c.SetConfigName("")

	f2, err := ioutil.TempFile("", "testconfig.json")
	if err != nil {
		T.Error(err)
	}
	defer os.RemoveAll(f2.Name())
	if err := os.Setenv("CONFIG_FILE", f2.Name()); err != nil {
		T.Error(err)
	}
	if exp, got := f2.Name(), c.ConfigName(); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}

	f3, err := ioutil.TempFile("", "env_with_prefix.json")
	if err != nil {
		T.Error(err)
	}
	defer os.RemoveAll(f3.Name())
	c.SetEnvPrefix("TESTPREFIX")
	if err := os.Setenv("TESTPREFIX_CONFIG_FILE", f3.Name()); err != nil {
		T.Error(err)
	}
	if exp, got := f3.Name(), c.ConfigName(); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestString(T *testing.T) {
	c := NewConfig()
	key, val := "test.string", "teststring"
	c.Set(key, val)
	if exp, got := val, c.GetString(key); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestBool(T *testing.T) {
	c := NewConfig()
	key, val := "test.bool", "true"
	c.Set(key, val)
	if exp, got := true, c.GetBool(key); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key2, val2 := "test.bool2", true
	c.Set(key2, val2)
	if exp, got := val2, c.GetBool(key2); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestStringSlice(T *testing.T) {
	c := NewConfig()
	key, val := "test.slice", []string{"one", "two"}
	c.Set(key, val)
	if exp, got := val, c.GetStringSlice(key); len(exp) != len(got) {
		T.Errorf("invalid result: expected=%#v got=%#v (%[2]T)", exp, got)
	}
	key2, val2 := "test.slice2", "lonelystring"
	c.Set(key2, val2)
	if exp, got := val2, c.GetStringSlice(key2); len(got) != 1 {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestStringMap(T *testing.T) {
	c := NewConfig()
	key, val := "test.map", map[string]interface{}{"one": "one", "two": "two"}
	c.Set(key, val)
	if exp, got := val, c.GetStringMap(key); len(exp) != len(got) {
		T.Errorf("invalid result: expected=%#v got=%#v (%[2]T)", exp, got)
	}
	key2, val2 := "test.map2", "lonelystring"
	c.Set(key2, val2)
	if exp, got := val2, c.GetStringMap(key2); len(got) != 0 {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestStringMapEnv(T *testing.T) {
	c := NewConfig()
	key, val := "test.map", map[string]interface{}{"one": "one", "two": "two"}
	os.Clearenv()
	os.Setenv("TEST_MAP_THREE", "three")
	c.Set(key, val)
	if exp, got := val, c.GetStringMap(key); len(exp)+1 != len(got) {
		T.Errorf("invalid result: expected=%#v got=%#v (%[2]T)", exp, got)
	}
	if v, ok := c.GetStringMap(key)["three"]; !ok || v != "three" {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", "three", v)
	}
	os.Clearenv()
}

func TestInt(T *testing.T) {
	c := NewConfig()
	key, val := "test.int", "10"
	c.Set(key, val)
	if exp, got := int(10), c.GetInt(key); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key2, val2 := "test.int2", int64(42)
	c.Set(key2, val2)
	if exp, got := int(val2), c.GetInt(key2); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key3, val3 := "test.int3", uint64(421)
	c.Set(key3, val3)
	if exp, got := int(val3), c.GetInt(key3); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key4, val4 := "test.int4", int32(423)
	c.Set(key4, val4)
	if exp, got := int(val4), c.GetInt(key4); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key5, val5 := "test.int5", uint32(424)
	c.Set(key5, val5)
	if exp, got := int(val5), c.GetInt(key5); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key6, val6 := "test.int6", int(425)
	c.Set(key6, val6)
	if exp, got := int(val6), c.GetInt(key6); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestInt64(T *testing.T) {
	c := NewConfig()
	key, val := "test.int", "10"
	c.Set(key, val)
	if exp, got := int64(10), c.GetInt64(key); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key2, val2 := "test.int2", int64(42)
	c.Set(key2, val2)
	if exp, got := int64(val2), c.GetInt64(key2); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key3, val3 := "test.int3", uint64(421)
	c.Set(key3, val3)
	if exp, got := int64(val3), c.GetInt64(key3); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key4, val4 := "test.int4", int32(423)
	c.Set(key4, val4)
	if exp, got := int64(val4), c.GetInt64(key4); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key5, val5 := "test.int5", uint32(424)
	c.Set(key5, val5)
	if exp, got := int64(val5), c.GetInt64(key5); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key6, val6 := "test.int6", int(425)
	c.Set(key6, val6)
	if exp, got := int64(val6), c.GetInt64(key6); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestFloat(T *testing.T) {
	c := NewConfig()
	key, val := "test.float", "1.1"
	c.Set(key, val)
	if exp, got := 1.1, c.GetFloat64(key); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key2, val2 := "test.float2", 2.2
	c.Set(key2, val2)
	if exp, got := val2, c.GetFloat64(key2); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestDuration(T *testing.T) {
	c := NewConfig()
	key, val := "test.duration", "10s"
	c.Set(key, val)
	if exp, got := 10*time.Second, c.GetDuration(key); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	key2, val2 := "test.duration2", 5*time.Second
	c.Set(key2, val2)
	if exp, got := val2, c.GetDuration(key2); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestInterface(T *testing.T) {
	testcases := map[string]interface{}{
		"test.one":   "string",
		"test.two":   10,
		"test.three": 1.1,
		"test.four":  true,
		"test.five":  time.Second,
		"test.siz":   []string{"one", "two"},
	}
	c := NewConfig()
	for n, v := range testcases {
		c.Set(n, v)
		if exp, got := toString(v), c.GetString(n); exp != got {
			T.Errorf("%s: invalid result: expected=%v got=%v (%[2]T)", n, exp, got)
		}
	}
}

var testcfg = `{
 "test": {
 	"one": "string",
 	"two": 10,
 	"three": 3.4,
 	"four": "2s",
 	"five": ["one", "two"],
 	"six": [{
 		"idx": 0,
 		"value": 10
 	},{
 		"idx": 1,
 		"value": 11
 	}]
  }
}`

func TestUnmarshal(T *testing.T) {
	c := NewConfig()
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	if exp, got := "string", c.GetString("test.one"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	if exp, got := int64(10), c.GetInt64("test.two"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	if exp, got := 3.4, c.GetFloat64("test.three"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	if exp, got := 2*time.Second, c.GetDuration("test.four"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
	if exp, got := []string{"one", "two"}, c.GetStringSlice("test.five"); len(exp) != len(got) {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestDefaults(T *testing.T) {
	c := NewConfig()
	c.SetDefault("test.default", 42)
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	if exp, got := int64(42), c.GetInt64("test.default"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestSet(T *testing.T) {
	c := NewConfig()
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	newval := "otherstring"
	c.Set("test.one", newval)
	if exp, got := newval, c.GetString("test.one"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestEnv(T *testing.T) {
	c := NewConfig()
	c.SetDefault("test.one", "defaultstring")
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	newval := "envstring"
	c.SetEnvPrefix("TESTPREFIX")
	if err := os.Setenv("TESTPREFIX_TEST_ONE", newval); err != nil {
		T.Error(err)
	}
	if exp, got := newval, c.GetString("test.one"); exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestAllSettingsFile(T *testing.T) {
	c := NewConfig()
	c.SetDefault("test.one", "defaultstring")
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	all := c.All()
	if all == nil {
		T.Fatalf("invalid result: nil map for all settings")
	}
	l1, ok := all["test"]
	if !ok {
		T.Fatalf("invalid result: missing map entry 'test'")
	}
	l1v, ok := l1.(map[string]interface{})
	if !ok {
		T.Fatalf("invalid result: invalid type %T for map entry 'test'", l1)
	}
	l2, ok := l1v["one"]
	if !ok {
		T.Fatalf("invalid result: missing map entry 'test.one'")
	}
	l2v, ok := l2.(string)
	if !ok {
		T.Fatalf("invalid result: invalid type %T for map entry 'test.one'", l2)
	}
	if exp, got := "string", l2v; exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestAllSettingsEnv(T *testing.T) {
	c := NewConfig()
	c.SetDefault("test.one", "defaultstring")
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	newval := "envstring"
	c.SetEnvPrefix("TESTPREFIX")
	if err := os.Setenv("TESTPREFIX_TEST_ONE", newval); err != nil {
		T.Error(err)
	}
	all := c.All()
	if all == nil {
		T.Fatalf("invalid result: nil map for all settings")
	}
	l1, ok := all["test"]
	if !ok {
		T.Fatalf("invalid result: missing map entry 'test'")
	}
	l1v, ok := l1.(map[string]interface{})
	if !ok {
		T.Fatalf("invalid result: invalid type %T for map entry 'test'", l1)
	}
	l2, ok := l1v["one"]
	if !ok {
		T.Fatalf("invalid result: missing map entry 'test.one'")
	}
	l2v, ok := l2.(string)
	if !ok {
		T.Fatalf("invalid result: invalid type %T for map entry 'test.one'", l2)
	}
	if exp, got := newval, l2v; exp != got {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
	}
}

func TestForEachFile(T *testing.T) {
	c := NewConfig()
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	num := 0
	c.ForEach("test.six", func(c *Config) error {
		if exp, got := int64(num), c.GetInt64("idx"); exp != got {
			T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
		}
		if exp, got := int64(num+10), c.GetInt64("value"); exp != got {
			T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
		}
		num++
		return nil
	})
	if num != 2 {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", 2, num)
	}
}

func TestForEachEnv(T *testing.T) {
	c := NewConfig()
	if err := c.ReadConfig([]byte(testcfg)); err != nil {
		T.Error(err)
	}
	c.SetEnvPrefix("TESTPREFIX")
	if err := os.Setenv("TESTPREFIX_TEST_SIX_2_IDX", "2"); err != nil {
		T.Error(err)
	}
	if err := os.Setenv("TESTPREFIX_TEST_SIX_2_VALUE", "12"); err != nil {
		T.Error(err)
	}
	num := 0
	c.ForEach("test.six", func(c *Config) error {
		if exp, got := int64(num), c.GetInt64("idx"); exp != got {
			T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
		}
		if exp, got := int64(num+10), c.GetInt64("value"); exp != got {
			T.Errorf("invalid result: expected=%v got=%v (%[2]T)", exp, got)
		}
		num++
		return nil
	})
	if num != 3 {
		T.Errorf("invalid result: expected=%v got=%v (%[2]T)", 3, num)
	}
}
