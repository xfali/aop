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
	"sync"
)

type advisor struct {
	pointCut PointCut
	advice   Advice
}

type adviceData struct {
	advice     Advice
	invocation Invocation
}

type chainProxy struct {
	t        reflect.Type
	value    reflect.Value
	advisors []advisor

	methodLocker sync.Mutex
	methodIndex  map[string]reflect.Method

	adviceLocker sync.Mutex
	adviceDatas  map[string]*adviceData
}

func New(obj interface{}) *chainProxy {
	ret := &chainProxy{
		t:           reflect.TypeOf(obj),
		value:       reflect.ValueOf(obj),
		methodIndex: make(map[string]reflect.Method),
		adviceDatas: make(map[string]*adviceData),
	}
	return ret
}

func (aop *chainProxy) AddAdvisor(pointCut PointCut, advice Advice) Proxy {
	aop.advisors = append(aop.advisors, advisor{
		pointCut: pointCut,
		advice:   advice,
	})
	return aop
}

func (aop *chainProxy) Call(method string, params ...interface{}) (ret []interface{}, err error) {
	mt, _, err := aop.findMethod(method)
	if err != nil {
		return nil, err
	}
	advice, invocation := aop.findAdvisor(mt, params...)
	if advice == nil {
		return call(aop.value.Method(mt.Index), params...)
	}
	return advice(invocation, params), nil
}

func (aop *chainProxy) findAdvisor(method reflect.Method, params ...interface{}) (Advice, Invocation) {
	aop.adviceLocker.Lock()
	defer aop.adviceLocker.Unlock()

	key := getKey(method, aop.t)
	if d, ok := aop.adviceDatas[key]; ok {
		return d.advice, d.invocation
	}
	var advice Advice
	var invocation Invocation = newInvocation(aop.value.Method(method.Index))
	last := invocation
	for i := len(aop.advisors) - 1; i >= 0; i-- {
		v := aop.advisors[i]
		if v.pointCut.Matches(method, aop.t, params...) {
			advice = v.advice
			last = invocation
			// first
			invocation = newChainInvocation(v.advice, last)
		}
	}

	aop.adviceDatas[key] = &adviceData{
		advice:     advice,
		invocation: last,
	}
	return advice, last
}

func getKey(method reflect.Method, t reflect.Type) string {
	return fmt.Sprintf("%s.%s.%s", t.String(), method.PkgPath, method.Name)
}

func (aop *chainProxy) findMethod(method string) (reflect.Method, bool, error) {
	aop.methodLocker.Lock()
	defer aop.methodLocker.Unlock()

	if mt, ok := aop.methodIndex[method]; ok {
		return mt, true, nil
	}

	mt, ok := aop.value.Type().MethodByName(method)
	if !ok {
		return reflect.Method{}, false, fmt.Errorf("Cannot found type %s with method %s ", aop.value.Type().String(), method)
	}
	aop.methodIndex[method] = mt
	return mt, false, nil
}

func (aop *chainProxy) invokeDefault(method string, params ...interface{}) ([]interface{}, error) {
	v, _, err := aop.findMethod(method)
	if err != nil {
		return nil, err
	}
	return call(aop.value.Method(v.Index), params...)
}

type chainInvocation struct {
	advice     Advice
	invocation Invocation
}

func (i *chainInvocation) Invoke(params []interface{}) []interface{} {
	return i.advice(i.invocation, params)
}

func newChainInvocation(advice Advice, invocation Invocation) *chainInvocation {
	return &chainInvocation{
		advice:     advice,
		invocation: invocation,
	}
}
