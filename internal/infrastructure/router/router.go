package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func New() *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	
	return &Router{
		engine: engine,
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) Group(relativePath string) *gin.RouterGroup {
	return r.engine.Group(relativePath)
}

func (r *Router) RegisterRoutes(register func(*gin.Engine)) {
	register(r.engine)
}

func (r *Router) RegisterGroupRoutes(relativePath string, register func(*gin.RouterGroup)) {
	group := r.engine.Group(relativePath)
	register(group)
}

func (r *Router) Start(addr string) error {
	return r.engine.Run(addr)
}