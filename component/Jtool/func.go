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
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(callFunc interface{}) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(callFunc).Pointer()).Name()
	funcName = strings.TrimSuffix(funcName, "-fm")
	index := strings.LastIndex(funcName, ".")
	funcName = funcName[index+1:]
	return funcName
}

