// Copyright (c) 2018-2023 KIDTSUNAMI
// Author: alex@kidtsunami.com

package config

import (
    "fmt"
    "strconv"
    "time"
)

func ParseTime(v string) (time.Time, error) {
    if num, err := strconv.ParseInt(v, 10, 64); err == nil {
        return time.Unix(num, 0), nil
    }
    for _, f := range []string{
        time.RFC3339,
        time.DateOnly,
        time.DateTime,
        time.UnixDate,
    } {
        if tm, err := time.Parse(f, v); err == nil {
            return tm, nil
        }
    }
    return time.Time{}, fmt.Errorf("invalid time format %q", v)
}
