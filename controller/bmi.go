package controller

import(
	"fmt"
)
func getBMI(height float32, weight float32)float32{
	BMI := float32(weight/(height*height/10000))
	return BMI
}
func main(){
	BMI:=getBMI(181,68)
	fmt.Println(BMI)
}