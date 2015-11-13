package main

import (
  "fmt"
)

type StringCount struct {
  name string
  count int
}

type StringCounts []StringCount

func (stringCount StringCount) ToPrintString() string {
  return fmt.Sprintf("%v: %v", stringCount.name, stringCount.count)
}

func (slice StringCounts) Len() int {
    return len(slice)
}

func (slice StringCounts) Less(i, j int) bool {
  if(slice[i].count == slice[j].count) {
    return slice[i].name < slice[j].name
  } else {
    return slice[i].count < slice[j].count
  }
}

func (slice StringCounts) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func mappingToStringCounts(mapping map[string]int) StringCounts{
  stringCounts := make(StringCounts, 0)
  for key, val := range mapping {
    stringCount := StringCount{ name: key, count: val }
    stringCounts = append(stringCounts, stringCount)
  }
  return stringCounts
}
