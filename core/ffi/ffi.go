package main

import (
  "C"
  "encoding/json"
)

type Paylaod struct {
  Method string `json:"method"`
  Params json.RawMessage `json:"params"`
}

//export call
func call(payload *C.char) *C.char {
  return nil
}

func main(){}
