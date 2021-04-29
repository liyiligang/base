// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2020/2/16 11:14
// Description:

package errorBox

//PB
//func PBUnmarshalSuccess(err error, data interface{}) bool {
//	if err != nil {
//		Jlog.Warn("消息解码失败", "errorBox", err, "原始数据", data)
//		return false
//	}
//	return true
//}
//
//func PBUnmarshalSuccessNoNil(err error, data []byte) error {
//	if err != nil {
//		Jlog.Warn("消息解码失败", "errorBox", err, "原始数据", data)
//		return err
//	}
//
//	if len(data) == 0 {
//		Jlog.Warn("消息不能为空")
//		err = errors.New("消息不能为空")
//		return err
//	}
//
//	return nil
//}
//
//func PBMarshalSuccess(err error, data interface{}) bool {
//	if err != nil {
//		Jlog.Warn("消息编码失败", "errorBox", err, "原始数据", data)
//		return false
//	}
//	return true
//}
