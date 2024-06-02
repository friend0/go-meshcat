package internal

import "github.com/labstack/echo/v4"

func (s *Server) StaticHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if err := next(ctx); err != nil {
			ctx.Error(err)
			return err
		} else {
			ctx.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		return nil
	}
}

func (s *Server) StaticHandler(fsRoot string) echo.HandlerFunc {
	subFs := echo.MustSubFS(s.Router.Filesystem, fsRoot)
	return echo.StaticDirectoryHandler(subFs, false)
}
