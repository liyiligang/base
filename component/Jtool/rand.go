// Copyright 2021 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2021/02/18 7:38
// Description: 随机数工具

package Jtool

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandInt(low int, high int) (int, error) {
	if low < 0 || high < 0 {
		return 0, errors.New("随机数区间不能为负数")
	}
	if high < low {
		return 0, errors.New("最大随机数不能小于最少随机数")
	}
	if low == high {
		return low, nil
	}
	return rand.Intn(int(high-low)) + low, nil
}