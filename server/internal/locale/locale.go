package locale

import (
	"bbs-go/internal/config"
	"fmt"
)

func init() {
	if err := InitLocale(config.Instance().Language); err != nil {
		panic(fmt.Errorf("failed to init locale. %w", err))
	}
}
