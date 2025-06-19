package main

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func pop(lst []string) (string, []string) {
	return lst[len(lst)-1], lst[:len(lst)-1]
}
