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

package test

import (
	"fmt"
	"github.com/xfali/aop"
	"testing"
)

type a interface {
	AGet(a string) string
}

type b interface {
	BGet(b string) string
}

type testStruct struct {
}

func (t *testStruct) AGet(a string) string {
	fmt.Println(a)
	return a
}

func (t testStruct) BGet(b string) (string, int) {
	fmt.Println(b)
	return b, len(b)
}

func (t testStruct) Concat(a, b string) (string, int) {
	ret := a + b
	fmt.Println(ret)
	return ret, len(ret)
}

func TestAop(t *testing.T) {
	o := &testStruct{}
	p := aop.New(o)
	err := p.AddJoinPoint("xx", func(invocation aop.Invocation, params []interface{}) (ret []interface{}) {
		fmt.Println("before AGet")
		v := invocation.Invoke(params)
		fmt.Println("after AGet")
		return v
	})
	if err == nil {
		t.Fatal("expect not nil but get ", err)
	}

	err = p.AddJoinPoint("AGet", func(invocation aop.Invocation, params []interface{}) (ret []interface{}) {
		fmt.Println("before AGet")
		v := invocation.Invoke(params)
		fmt.Println("after AGet")
		return v
	})
	if err != nil {
		t.Fatal("expect nil but get ", err)
	}

	v, err := p.Call("AGet", "hello")
	if err != nil {
		t.Fatal("expect nil but get ", err)
	}

	if v[0].(string) != "hello" {
		t.Fatal("expect hello but get ", v[0].(string))
	}

	err = p.AddJoinPoint("BGet", func(invocation aop.Invocation, params []interface{}) (ret []interface{}) {
		fmt.Println("before BGet")
		params[0] = params[0].(string) + "p1"
		v := invocation.Invoke(params)
		fmt.Println("after BGet")
		v[0] = v[0].(string) + "r1"
		return v
	})
	if err != nil {
		t.Fatal("expect nil but get ", err)
	} else {
		t.Log(v[0].(string))
	}

	v, err = p.Call("BGet", "world")
	if err != nil {
		t.Fatal("expect nil but get ", err)
	}

	if v[0].(string) != "worldp1r1" {
		t.Fatal("expect worldp1r1 but get ", v[0].(string))
	} else {
		t.Log(v[0].(string))
	}
	if v[1].(int) != len("worldp1") {
		t.Fatal("expect 7 but get ", v[1].(int))
	}

	v, err = p.Call("NotExistMethod", "?")
	if err == nil {
		t.Fatal("expect nil but ", err)
	} else {
		t.Log(err)
	}

	v, err = p.Call("Concat", "hello", "world")
	if err != nil {
		t.Fatal("expect nil but get ", err)
	} else {
		t.Log(v[0].(string))
	}

	if v[0].(string) != "helloworld" {
		t.Fatal("expect helloworld but get ", v[0].(string))
	}
	if v[1].(int) != len("helloworld") {
		t.Fatal("expect 10 but get ", v[1].(int))
	}
}
