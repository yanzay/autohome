package main

import (
  "fmt"
  "testing"
)

type TestObject struct {
  Some    int
  Another string
}

func TestCreate(t *testing.T) {
  st := NewSqliteStorage("test")
  to := TestObject{Some: 1, Another: "helo"}
  st.CreateTable(to)
  err := st.Create(to)
  if err != nil {
    t.Errorf("Got an error %v", err)
  }
  var rows []TestObject
  st.All(TestObject{}).Rows(&rows)
  fmt.Printf("One: %v\n", to)
  fmt.Printf("All: %v\n", rows)
}
