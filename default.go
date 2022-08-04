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
	"github.com/xfali/aop/methodfunc"
	"reflect"
	"sync"
)

type meta struct {
	advice     Advice
	invocation Invocation
	lock       sync.Mutex
}

type simpleProxy struct {
	t           reflect.Type
	value       reflect.Value
	pointCuts   map[PointCut]*meta
	methodIndex map[string]reflect.Method
}

func NewSimple(obj interface{}) *simpleProxy {
	ret := &simpleProxy{
		t:           reflect.TypeOf(obj),
		value:       reflect.ValueOf(obj),
		pointCuts:   make(map[PointCut]*meta),
		methodIndex: make(map[string]reflect.Method),
	}
	return ret
}

func (aop *simpleProxy) AddAdvisor(pointCut PointCut, advice Advice) Proxy {
	aop.pointCuts[pointCut] = &meta{advice: advice}
	return aop
}

func (aop *simpleProxy) Call(method string, params ...interface{}) (ret []interface{}, err error) {
	mt, err := aop.findMethod(method)
	if err != nil {
		return nil, err
	}
	m, err := aop.findJoinPoint(mt, params...)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return call(aop.value.Method(mt.Index), params...)
	}
	return m.advice(m.invocation, params), nil
}

func (aop *simpleProxy) findJoinPoint(method reflect.Method, params ...interface{}) (*meta, error) {
	for k, v := range aop.pointCuts {
		if k.Matches(method, aop.t, params...) {
			v.lock.Lock()
			if v.invocation == nil {
				v.invocation = createInvocation(aop.value.Method(method.Index))
			}
			v.lock.Unlock()
			return v, nil
		}
	}
	return nil, nil
}

func (aop *simpleProxy) findMethod(method string) (reflect.Method, error) {
	if mt, ok := aop.methodIndex[method]; ok {
		return mt, nil
	}

	mt, ok := aop.value.Type().MethodByName(method)
	if !ok {
		return reflect.Method{}, fmt.Errorf("Cannot found type %s with method %s ", aop.value.Type().String(), method)
	}
	aop.methodIndex[method] = mt
	return mt, nil
}

func (aop *simpleProxy) invokeDefault(method string, params ...interface{}) ([]interface{}, error) {
	v, err := aop.findMethod(method)
	if err != nil {
		return nil, err
	}
	return call(aop.value.Method(v.Index), params...)
}

func call(method reflect.Value, params ...interface{}) ([]interface{}, error) {
	pn := method.Type().NumIn()
	if pn != len(params) {
		return nil, fmt.Errorf("Method expect param size: %d but get %d ", pn, len(params))
	}
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
		return results, nil
	}
	return nil, nil
}

type defaultInvocation struct {
	method reflect.Value
}

func (i *defaultInvocation) Invoke(params []interface{}) []interface{} {
	ret, err := call(i.method, params...)
	if err != nil {
		panic(err)
	}
	return ret
}

func createInvocation(method reflect.Value) Invocation {
	return &defaultInvocation{method: method}
}

type defaultPointCut string

func (p defaultPointCut) Matches(method reflect.Method, instanceType reflect.Type, params ...interface{}) bool {
	if !methodfunc.CheckMethod(method, instanceType, params) {
		return false
	}
	return string(p) == method.Name
}

func PointCutMethodName(method string) PointCut {
	return defaultPointCut(method)
}
