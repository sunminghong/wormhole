/*=============================================================================
#     FileName: interface.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-13 17:23:18
#      History:
=============================================================================*/


/*
定义wormhole所有的接口，枚举
*/
package wormhole

import (
    "reflect"
)


type Inherit struct {
    sub interface{}
    subMethodsMap map[string]reflect.Value

    defaultMethods []string
}

func NewInherit(defaultMethods ...string) *Inherit {
    i:= &Inherit{
        defaultMethods: defaultMethods,
        subMethodsMap : make(map[string]reflect.Value),
    }
    return i
}


func (i *Inherit) register(p interface{}, methods []string) {
    i.sub = p
    if len(methods) > 0 {
        sub := reflect.ValueOf(p)
        for _,mname := range methods {
            method := sub.MethodByName(mname)
            if method.IsValid() {
                i.subMethodsMap[mname] = method
            }
        }
    }
}

func (i *Inherit) RegisterSub(p interface{}, methods ...string) {
    i.register(p, i.defaultMethods)
    i.register(p, methods)
}


func (i *Inherit) CallSub(method string, params ...interface{}) {
    if method,ok := i.subMethodsMap[method]; ok {
        args := make([]reflect.Value, len(params))
        for i,param := range params {
            args[i] = reflect.ValueOf(param)
        }
        method.Call(args)
        return
    }
    panic("panic:'s Method ProcessPackets need override write by sub object")
}


