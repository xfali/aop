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

package methodfunc

import "reflect"

func CheckMethod(method reflect.Method, instanceType reflect.Type, params []interface{}) bool {
	mt := instanceType.Method(method.Index)
	return mt.Name == method.Name && mt.Type == method.Type && mt.PkgPath == method.PkgPath
}
