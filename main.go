package main

import (
	"bytes"
	"fmt"
	"math/rand"
)

func main()  {
	//var ch []byte = []byte{'\x03', '\x23'}
	//str := string(ch)
	//
	//fmt.Print(str)

	//array := [7]int{1,2,3,4,5,6,7}
	//fmt.Printf("%v\n",array)
	//array1 := array[7:]
	//fmt.Printf("%v\n",array1)

	//vect := make([]int, 10)
	////1. 赋初值
	//for i, _ := range vect {
	//	vect[i] = i
	//}
	//for i := 0; i < 3; i++ {
	//	k := rand.Intn(10)
	//	vect[k] = -1 //将其中三个值设为-1
	//}
	//vect[9] = -1 //***注：将最后一个元素置为－1
	//fmt.Println("before:", vect)
	//
	////2.遍历切片，并删除值为-1的元素
	//for i, v := range vect {
	//	if v == -1 {
	//		vect = append(vect[:i], vect[i+1:]...)
	//	}
	//}
	//fmt.Println("after:", vect)

	buffer := [6]byte{'\x01','\x01','\x01','\x01','\x01','\x01'}
	reader := bytes.NewReader(buffer)
}
