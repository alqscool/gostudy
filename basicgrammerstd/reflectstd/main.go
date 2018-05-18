package main

import (
	"fmt"
	"reflect"
)

type user struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

type room struct {
	F string `species:"gopher" color:"blue"`
}

func main()  {
	user := &user{"lqs","cool"}
	s:= reflect.TypeOf(user).Elem()  //通过反射获取type定义

	for i:=0; i<s.NumField(); i++ {
		fmt.Println(s.Field(i).Tag)  //将tag输出出来
	}

	room := room{}
	rt := reflect.TypeOf(room)
	field := rt.Field(0)
	fmt.Println(field.Tag.Get("color"),field.Tag.Get("species"))
}



