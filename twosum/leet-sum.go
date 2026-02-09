package main

import (
	"fmt"
)

func main () {
	fmt.Printf("hello")

	nums := [5]int{4,8,9,11,20};
	

	result := twoSum(nums, 12);

	fmt.Println(result)
}
func twoSum (nums [5]int, target int) []int {
	for i := 0; i < len(nums); i++{
		for j := i +1 ; j <len(nums); j++{
			sum := nums[i] + nums[j]
			if sum == target{
				return []int{i,j}
			}
		}
	}

	return []int{0,0}
}