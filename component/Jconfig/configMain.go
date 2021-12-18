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

package Jconfig

import (
	"bytes"
	"github.com/spf13/viper"
)

func ReadConfigFromPath(config interface{}, path string) error {
	viperPath := viper.New()
	viperPath.SetConfigName("config")
	viperPath.SetConfigType("toml")
	if path == "" {
		path = "."
	}
	viperPath.AddConfigPath(path)
	if err := viperPath.ReadInConfig(); err != nil {
		return err
	}
	if err := viperPath.Unmarshal(config); err != nil {
		return err
	}
	return nil
}

func ReadConfigFromByte(config interface{}, data []byte) error {
	viperByte := viper.New()
	viperByte.SetConfigType("toml")
	if err := viperByte.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}
	if err := viperByte.Unmarshal(config); err != nil {
		return err
	}
	return nil
}
