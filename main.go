package main

import (
	"fmt"
	"time"

	"github.com/HotCodeGroup/warscript-tester/tester"
)

func main() {
	d1, err := tester.InitDocker("pong", 5000, `{me.x+=3; me.y+=5;}`)
	if err != nil {
		fmt.Println("docker nort started:", err)
		return
	}
	defer d1.Kill()
	d2, err := tester.InitDocker("pong", 5001, `{me.x+=3; me.y+=5;}`)
	if err != nil {
		fmt.Println("docker nort started:", err)
		return
	}
	defer d2.Kill()
	st := []byte(`{"me":{"x":1,"y":2},"enemy":{}, "ball":{}}`)
	for {
		t := time.NewTimer(time.Second / 2)
		<-t.C
		var err error
		st, err = d1.SendState(st)
		if err != nil {
			fmt.Println("err occured:", err)
			return
		}
		fmt.Println("d1", string(st))
		st, err = d2.SendState(st)
		if err != nil {
			fmt.Println("err occured:", err)
			return
		}
		fmt.Println("d2", string(st))
	}
}
