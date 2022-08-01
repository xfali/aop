/*
 * Copyright (C) 2022, Xiongfa Li.
 * All rights reserved.
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
	"fmt"
	"reflect"
)

type meta struct {
	advice     Advice
	invocation Invocation
}

type defaultProxy struct {
	value       reflect.Value
	joinPoints  map[JoinPoint]*meta
	methodIndex map[string]int
}

func New(obj interface{}) *defaultProxy {
	ret := &defaultProxy{
		value:       reflect.ValueOf(obj),
		joinPoints:  make(map[JoinPoint]*meta),
		methodIndex: map[string]int{},
	}
	return ret
}

func (aop *defaultProxy) AddJoinPoint(point JoinPoint, advice Advice) error {
	m := aop.value.MethodByName(string(point))
	if !m.IsValid() || m.IsZero() {
		return fmt.Errorf("Cannot found type %s with JoinPoint %s ", aop.value.Type().String(), point)
	}
	aop.joinPoints[point] = &meta{
		advice:     advice,
		invocation: createInvocation(m),
	}
	return nil
}

func (aop *defaultProxy) Call(method string, params ...interface{}) (ret []interface{}, err error) {
	m := aop.findJoinPoint(method)
	if m == nil {
		return aop.invokeDefault(method, params...)
	}

	return m.advice(m.invocation, params), nil
}

func (aop *defaultProxy) findJoinPoint(method string) *meta {
	m := aop.joinPoints[JoinPoint(method)]
	if m == nil {
		return nil
	}
	return m
}

func (aop *defaultProxy) invokeDefault(method string, params ...interface{}) ([]interface{}, error) {
	if index, ok := aop.methodIndex[method]; ok {
		return call(aop.value.Method(index), params...), nil
	}

	mt, ok := aop.value.Type().MethodByName(method)
	if !ok {
		return nil, fmt.Errorf("Cannot found type %s with method %s ", aop.value.Type().String(), method)
	}
	aop.methodIndex[method] = mt.Index

	return call(aop.value.Method(mt.Index), params...), nil
}

func call(method reflect.Value, params ...interface{}) []interface{} {
	var ret []reflect.Value
	if len(params) > 0 {
		pv := make([]reflect.Value, len(params))
		for i, p := range params {
			pv[i] = reflect.ValueOf(p)
		}
		ret = method.Call(pv)
	} else {
		ret = method.Call(nil)
	}

	if len(ret) > 0 {
		results := make([]interface{}, len(ret))
		for i, v := range ret {
			results[i] = v.Interface()
		}
		return results
	}
	return nil
}

type defaultInvocation struct {
	method reflect.Value
}

func (i *defaultInvocation) Invoke(params []interface{}) (ret []interface{}) {
	return call(i.method, params...)
}

func createInvocation(method reflect.Value) Invocation {
	return &defaultInvocation{method: method}
}
