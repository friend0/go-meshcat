package internal

func (s *Server) Routes() {
	s.Router.GET("app*", s.StaticHandler("web/meshcat/dist"))
	s.Router.GET("/", s.WsHandler())
}
