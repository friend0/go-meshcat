package internal

func (s *Server) Routes() {
	s.Router.GET("/ws", s.WsHandler())
	s.Router.GET("/data/*", s.StaticHandler("web/meshcat/data"))
	s.Router.GET("/*", s.StaticHandler("web/meshcat/dist"))
}
