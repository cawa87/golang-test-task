// +build swipe

package transport

import (
	"github.com/kkucherenkov/golang-test-task/controller"
	. "github.com/swipe-io/swipe/v2"
)

func Swipe() {
	Build(
		Service(
			HTTPServer(),

			Interface((*controller.ScrapperController)(nil), ""),

			ClientsEnable([]string{"go"}),

			JSONRPCEnable(),

			OpenapiEnable(),
			OpenapiOutput("./docs"),
			OpenapiInfo("Service", "Sites availability checker.", "v1.0.0"),

			MethodDefaultOptions(
				Logging(true),
				Instrumenting(true),
			),
		),
	)
}
