# ctugo
golang SDK for ctu product

[![Build Status](https://travis-ci.org/dingxiangtech/ctugo.svg?branch=master)](https://travis-ci.org/dingxiangtech/ctugo)

## Example

```go
package main

import (
	"log"
	"github.com/dingxiangtech/ctugo"
)

func main() {
	conn := ctugo.NewEngineConnection("http://127.0.0.1:7776/ctu/event.do", "05622f1ab6be69567d65a6e377edfef0", "b2a8a90190fff591bd93bfd99e268438")
	resp, err := conn.CallRiskEngine("marketing_evt2", map[string]interface{}{"ip": "1.2.3.4", "email": "abc@def.com"})
	log.Println(resp)  // &{4fbd05d5-5fc8-44d8-8399-de0018d4e6fc INVALID_REQUEST_PARAMS {ACCEPT  []   [] [] marketing_evt map[_cost_time:1 _error_policy:[优惠券 归属地冲突注册限制] _rule_eval_error:[ruleId=128, seqNumber=1, missing params:ext_youhuiid=null  ruleId=129, seqNumber=1, missing params: ext_phonelocation=null] _success_execute:true] map[]}}
	log.Println(err)  // nil
}
```
