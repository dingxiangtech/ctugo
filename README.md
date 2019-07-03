# ctugo
golang SDK for ctu product

## Example

```go
import (
	"log"
	"github.com/dingxiangtech/ctugo"
)

func ExampleCallRiskEngine() {
	conn := NewEngineConnection("http://127.0.0.1:7776/ctu/event.do", "05622f1ab6be69567d65a6e377edfef0", "b2a8a90190fff591bd93bfd99e268438")
	resp, err := conn.CallRiskEngine("marketing_evt2", map[string]interface{}{"ip": "1.2.3.4", "email": "abc@def.com"})
	log.Println(resp)
	log.Println(err)  // nil
}
```
