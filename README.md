### go-eventbus

A simple eventbus written in go.

### Importing

    import github.com/LightOfReason/go-eventbus

### Example
```go
import (
	"fmt"
	"github.com/LightOfReason/go-eventbus"
)

type SimpleStruct struct{}

//go-eventbus uses a naming convention to detect eventhandler.
//Every method starting with EventHandler... and exactly one parameter
//is an handler.
func (st *SimpleStruct) EventHandlerSayHello(s string) {
	fmt.Println("Hello ", s)
}

func main() {
	//shutdown
	defer eventbus.Shutdown()

	simple := &SimpleStruct{}

	eventbus.Subscribe(simple)
	eventbus.Publish("World")
	eventbus.Unsubscribe(simple)
}
```

### Documentation

Visit the docs on [gopkgdoc](http://godoc.org/github.com/LightOfReason/go-eventbus)

### License
This project is licensed under the [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).