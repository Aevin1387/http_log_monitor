package main

import (
  "fmt"
)

// StringCount is used to contain a count
// of some type of string. Used for sorting
// data in a map[string][int]
type StringCount struct {
  name string
  count int
}

// StringCounts is an array of StringCount
type StringCounts []StringCount

// ToPrintString creates a printable version of a StringCount
func (stringCount StringCount) ToPrintString() string {
  return fmt.Sprintf("%v: %v", stringCount.name, stringCount.count)
}

// Len returns the number of StringCount in a StringCounts
// required for sorting
func (slice StringCounts) Len() int {
    return len(slice)
}

// Less determines which StringCount belongs before another.
// If the count of the StringCount is the same, it sorts by name.
func (slice StringCounts) Less(i, j int) bool {
  if(slice[i].count == slice[j].count) {
    return slice[i].name < slice[j].name
  }

  return slice[i].count < slice[j].count
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
