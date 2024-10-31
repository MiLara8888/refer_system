package refersrest

import "net/http"

var (
	postopt = []string{http.MethodOptions, http.MethodPost}
	get     = []string{http.MethodGet}
)

func (s *Rest) initializeRoutes() {
	s.Routes.UseRawPath = true
	s.Routes.UnescapePathValues = false

	app := s.Routes.Group("", CORSMiddleware())

	user := app.Group("/user")
	{
		//реиcтрация нового пользователя
		user.Match(postopt, "/register", s.Register)

		//авторизация пользователя, создание свежего токена
		user.Match(postopt, "/auth", s.Auth)
	}

	res_code := app.Group("/refer", s.AuthRequired())
	{
		//создание кода
		res_code.Match(postopt, "/update_code", s.CreateRef)

		//удаление кода
		res_code.Match(postopt, "/delete", s.DeleteRef)

		//получение кода по email
		res_code.Match(get, "/getcode", s.GetCode)
	}

	subscribers := app.Group("/referals")
	{
		//создание подписчика на рефер
		subscribers.Match(postopt, "/register", s.CreateReferals)

		//список подписчиков на рефер код
		subscribers.Match(get, "/list", s.ReferalsList)
	}

}
