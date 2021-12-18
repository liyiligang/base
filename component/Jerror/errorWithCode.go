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

package Jerror

import (
	"github.com/liyiligang/base/component/Jtool"
)

// ErrWithCode 带错误码的错误类型
type ErrWithCode struct {
	Code  		int
	Err   		error
}

func (err *ErrWithCode) Error() string {
	return err.Err.Error() + " with error code " + Jtool.IntToString(err.Code)
}

func NewErrWithCode(code int, err error) *ErrWithCode {
	return &ErrWithCode{
		Code:  code,
		Err:   err,
	}
}

func IsErrWithCode(err error) bool {
	_ ,ok := err.(*ErrWithCode)
	return ok
}

func AssertErrWithCode(err error) *ErrWithCode {
	errWithCode , _ := err.(*ErrWithCode)
	return errWithCode
}
