package base62

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	for i := range 1000 {
		fmt.Println(Encode(int64(i)))
	}
	// items := []string{}
	// for idx, r := range CODE62 {
	// 	items = append(items, fmt.Sprintf("\"%s\": %d", string(r), idx))
	// 	if (idx+1)%10 == 0 {
	// 		items = append(items, "\n")
	// 	}
	// }
	// fmt.Println(strings.Join(items, ","))
}
