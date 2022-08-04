/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package aop

import (
	"github.com/xfali/aop/methodfunc"
	"reflect"
	"regexp"
)

type regexpPointCut struct {
	typeRegexp   *regexp.Regexp
	methodRegexp *regexp.Regexp
}

func (p *regexpPointCut) Matches(method reflect.Method, instanceType reflect.Type, params ...interface{}) bool {
	if !methodfunc.CheckMethod(method, instanceType, params) {
		return false
	}
	tname := instanceType.String()
	if p.typeRegexp != nil && !p.typeRegexp.Match([]byte(tname)) {
		return false
	}
	if p.methodRegexp != nil && !p.methodRegexp.Match([]byte(method.Name)) {
		return false
	}
	return true
}

func PointCutRegExp(instanceType, method string) *regexpPointCut {
	ret := &regexpPointCut{
	}
	if instanceType != "" {
		ret.typeRegexp = regexp.MustCompile(instanceType)
	}
	if method != "" {
		ret.methodRegexp = regexp.MustCompile(method)
	}
	return ret
}
