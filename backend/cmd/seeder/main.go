package main

import (
	"os"

	"github.com/Kar-Su/uas-mobile.git/internal/package/script"
	"github.com/Kar-Su/uas-mobile.git/internal/providers"

	"github.com/samber/do/v2"
)

func args(injector do.Injector) bool {
	if len(os.Args) > 1 {
		flag := script.Commands(injector)
		return flag
	}

	return true
}

func main() {
	var (
		injector = do.New()
	)

	providers.RegisterProviders(injector)

	if !args(injector) {
		return
	}

}
