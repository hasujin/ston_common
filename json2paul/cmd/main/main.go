package main

import (
	"fmt"

	"github.com/hasujin/ston_common/json2paul"

	//	"strconv"
	"encoding/json"
)

type ReplyOperation string

// TODO: check spec
type RejectMessage struct {
	Op         ReplyOperation `json:"op"`
	Reason     string         `json:"reason"`
	ReqId      uint64         `json:"reqId"`
	Identifier string         `json:"identifier"`
}

func main() {

	structExample := RejectMessage{
		Op:         "abc123",
		Reason:     "reason123",
		ReqId:      654321,
		Identifier: "myID",
	}

	marshaled, err := json.Marshal(structExample)

	data := []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": "a0844674407370955abc161"
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

	data2 := []byte(`{
  "operation": {
    "type": "122",
    "rules": [{
      "constraint": {
        "constraint_id": "OR",
        "auth_constraints": [{
            "constraint_id": "ROLE",
            "role": "0",
            "sig_count": 1,
            "need_to_be_owner": false,
            "metadata": {}
          },
          {
            "constraint_id": "ROLE",
            "role": "2",
            "sig_count": 1,
            "need_to_be_owner": true,
            "metadata": {}
          }
        ]
      },
      "field": "services",
      "auth_type": "0",
      "auth_action": "EDIT",
      "old_value": "VALIDATOR",
      "new_value": []
    }]
  },

  "identifier": "21BPzYYrFzbuECcBV3M1FH",
  "reqId": 1514304094738044,
  "protocolVersion": 1,
  "signature": "3YVzDtSxxnowVwAXZmxCG2fz1A38j1qLrwKmGEG653GZw7KJRBX57Stc1oxQZqqu9mCqFLa7aBzt4MKXk4MeunVj"
}`)

	value, dataTytpe, offset, err := json2paul.Get(data, "person", "github", "followers")
	fmt.Println("Value : ", string(value), "\nType : ", dataTytpe, "\nOffset : ", offset, "\nerr :", err)
	/*
		if dataType == json2paul.number {
			number, _ = strconv.ParseUint(string(value), 10, 64)
		}
	*/

	fmt.Println("--------------------------------")

	value, dataTytpe, offset, err = json2paul.Get(marshaled, "reqId")
	fmt.Println("Value : ", string(value), "\nType : ", dataTytpe, "\nOffset : ", offset, "\nerr :", err)
	//fmt.Println(marshaled)

	fmt.Println("--------------------------------")

	value, dataTytpe, offset, err = json2paul.Get(data2, "operation", "rules", "[0]", "constraint", "auth_constraints", "[0]", "constraint_id") // Work!!
	fmt.Println("Value : ", string(value), "\nType : ", dataTytpe, "\nOffset : ", offset, "\nerr :", err)

	fmt.Println("--------------------------------")
	value, dataTytpe, offset, err = json2paul.Get(data2, "signature") // Work!!
	fmt.Println("Value : ", string(value), "\nType : ", dataTytpe, "\nOffset : ", offset, "\nerr :", err)

	fmt.Println("--------------------------------")
	//json2paul.Set(data, value_to_replace, value_position, "avatars", "[0]", "url")
	value, err = json2paul.Set(data2, []byte("signature replaced1"), "signature2") // Work!!
	fmt.Println("Set Result : ", string(value))

	fmt.Println("--------------------------------")
	//json2paul.Set(data, value_to_replace, value_position, "avatars", "[0]", "url")
	value, err = json2paul.Set(data2, []byte("signature replaced2"), "operation", "rules", "[0]", "constraint", "auth_constraints", "[0]", "constraint_id") // Work!!
	fmt.Println("Set Result : ", string(value))

	//Overwrite
	value, err = json2paul.Set(data2, []byte(`{"type2": "122"}`), "operation", "type") // Work!!
	fmt.Println("Set Result : ", string(value))

	value = json2paul.Delete(data2, "operation", "type") // Work!!
	fmt.Println("Delete Result : ", string(value))

	//============================================================================================================================================
	//============================================================================================================================================
	//============================================================================================================================================
	fmt.Println("\n------------------------------------------------------------------------")

	jsonObj := json2paul.New()
	// or gabs.Wrap(jsonObject) to work on an existing map[string]interface{}

	//Add
	jsonObj.Set(10, "outter", "inner", "value")
	jsonObj.SetP(20, "outter.inner.value2")
	jsonObj.Set(30, "outter", "inner2", "value3")

	fmt.Println(jsonObj.String())

	jsonParsed, _ := json2paul.ParseJSON(jsonObj.Bytes())

	//Add by point.path
	jsonParsed.SetP(30, "outter.inner.value3")
	fmt.Println(jsonParsed.String())

	//Replace, if exist
	jsonParsed.SetP(33, "outter.inner.value3")
	fmt.Println(jsonParsed.String())

}
