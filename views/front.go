package views

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func FrontPage() g.Node {
	return Page(
		"Canvas",
		"/",
		H1(g.Text(`Solutions to problems.`)),
		P(g.Text(`Do you have problems? We also had problems`)),
		P(g.Raw(`Then we created the <em>Canvas</em> app, and now we don't 😬`)))
}
