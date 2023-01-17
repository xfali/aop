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

import "reflect"

type Invocation interface {
	// Invoke 指定方法调用
	// params： 调用方法参数
	// ret： 返回调用后的结果
	Invoke(params []interface{}) (ret []interface{})

	// MethodName 返回方法名
	MethodName() string
}

type PointCut interface {
	Matches(method reflect.Method, instanceType reflect.Type, params ...interface{}) bool
}

type Advice func(invocation Invocation, params []interface{}) (ret []interface{})

type Proxy interface {
	// AddJoinPoint 增加连接点
	// joinPoint： 如方法名
	// advice： 在连接点触发的动作
	// 添加成功返回nil，否则返回错误
	AddAdvisor(pointCut PointCut, advice Advice) Proxy

	// Call 调用方法
	// method： 方法名
	// params： 方法参数
	// ret： 调用后返回的结果
	// err： 调用成功返回nil，失败返回错误
	Call(method string, params ...interface{}) (ret []interface{}, err error)
}
