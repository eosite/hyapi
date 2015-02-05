package main

type ResultCode struct {
	Code     string
	Name     string
	Category string
}
type Result struct {
	Codes []ResultCode
	Hint  []string
}
