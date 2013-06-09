package eventbus

import (
	"fmt"
	"testing"
)

type MyString struct {
	s string
}

func (s *MyString) EventHandlerTestit(b MyString) {
	fmt.Printf("Testit says: %s\n", b.s)
}

func (s *MyString) EventHandlerTestit2(b int) {
	fmt.Printf("Testit2 says: %d\n", b)
}

func Test1(t *testing.T) {
	defer func() {
		fmt.Println("Shutdown...")
		Shutdown()
	}()

	fmt.Println("Started...")

	fmt.Println("Subscribe...")
	var v *MyString = &MyString{"test"}
	Subscribe(v)

	fmt.Println("Publish...")
	for i := 0; i < 2; i++ {
		Publish(MyString{"A"})
		Publish(i)
	}

	fmt.Println("Unsubscribe...")
	Unsubscribe(v)
}
