package strcase_test

import (
	"fmt"
	"testing"

	"bbs-go/common/strs/strcase"
)

func TestFuck(t *testing.T) {
	fmt.Println(strcase.ToSnake("serviceCat1"))
	fmt.Println(strcase.ToSnake("serviceCat1Id"))
	fmt.Println(strcase.ToSnake("serviceCat1Id2"))
}
