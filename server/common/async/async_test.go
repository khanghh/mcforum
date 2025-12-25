package async_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"bbs-go/common/async"
)

func TestAsync(t *testing.T) {
	f1 := async.Exec(func() (int, error) {
		fmt.Println("execute method 1")
		return 1, nil
	})
	f2 := async.Exec(func() (int, error) {
		fmt.Println("execute method 2")
		time.Sleep(1 * time.Second)
		return 2, errors.New("failed")
	})

	fmt.Println(f1.Await())
	fmt.Println(f2.Await())
}

func TestAsyncWithNotTimeout(t *testing.T) {
	f := async.Exec(func() (string, error) {
		fmt.Println("execute method")
		time.Sleep(2 * time.Second)
		return "success", nil
	})

	fmt.Println(f.AwaitTimeout(3 * time.Second))
}

func TestAsyncWithTimeout(t *testing.T) {
	f := async.Exec(func() (string, error) {
		fmt.Println("execute method")
		time.Sleep(5 * time.Second)
		return "success", nil
	})
	fmt.Println(f.AwaitTimeout(3 * time.Second))
}
