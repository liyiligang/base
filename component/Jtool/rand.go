/*
 * Copyright 2021 liyiligang.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package Jtool

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRandInt 获取随机整数
func GetRandInt(low int, high int) (int, error) {
	if low < 0 || high < 0 {
		return 0, errors.New("random number interval cannot be negative")
	}
	if high < low {
		return 0, errors.New("the maximum random number cannot be less than the minimum random number")
	}
	if low == high {
		return low, nil
	}
	return rand.Intn(high-low) + low, nil
}

// GetRandChinese 获取随机中文
func GetRandChinese(min int, max int) string {
	length, _ := GetRandInt(min, max)
	c := make([]rune, length)
	for i := range c {
		h ,_ := GetRandInt(19968,40869)
		c[i]=rune(int64(h))
	}
	return string(c)
}

// GetRandString 获取随机字符串
func GetRandString(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}