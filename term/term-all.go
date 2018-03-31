package term

import "fmt"

func Println(obj fmt.Stringer) {
Println(obj.String())
}

func Println(text string) {
Print(text)
Print("\n")
}

func Print(obj fmt.Stringer) {
Print(obj.String())
}
