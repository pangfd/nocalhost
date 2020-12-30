/*
Copyright 2020 The Nocalhost Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package routers

import (
	"nocalhost/pkg/nocalhost-api/app/api/v1/application_cluster"
	"nocalhost/pkg/nocalhost-api/app/api/v1/applications"
	"nocalhost/pkg/nocalhost-api/app/api/v1/cluster"
	"nocalhost/pkg/nocalhost-api/app/api/v1/cluster_user"
	"nocalhost/pkg/nocalhost-api/napp"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	ginSwagger "github.com/swaggo/gin-swagger" //nolint: goimports
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"nocalhost/pkg/nocalhost-api/app/api"
	"nocalhost/pkg/nocalhost-api/app/api/v1/user"

	// import swagger handler
	_ "nocalhost/docs" // docs is generated by Swag CLI, you have to import it.
	"nocalhost/pkg/nocalhost-api/app/router/middleware"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// 使用中间件
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(mw...)

	// 404 Handler.
	g.NoRoute(api.RouteNotFound)
	g.NoMethod(api.RouteNotFound)
	g.Use(api.Recover)

	// Static resources
	//g.Static("/static", "./static")

	// Open only in the test environment, close online
	if viper.GetString("app.run_mode") == napp.ModeDebug {
		// swagger api docs
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// pprof router Performance analysis routing
		// Closed by default, can be opened in development environment
		// interview method: HOST/debug/pprof
		// Generate profile through HOST/debug/pprof/profile
		// View analysis graph go tool pprof -http=:5000 profile
		// see: https://github.com/gin-contrib/pprof
		pprof.Register(g)
	} else {
		// disable swagger docs for release  env=release
		g.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "env"))
	}

	g.POST("/v1/register", user.Register)
	g.POST("/v1/login", user.Login)

	u := g.Group("/v1/users")
	u.Use(middleware.AuthMiddleware())
	{
		u.GET("/:id", user.Get)
		u.GET("", user.GetList)
		u.POST("", user.Create)
		u.PUT("/:id", user.Update)
		u.DELETE("/:id", user.Delete)
	}

	m := g.Group("/v1/me")
	m.Use(middleware.AuthMiddleware())
	{
		m.GET("", user.GetMe)
	}

	// Clusters
	c := g.Group("/v1/cluster")
	c.Use(middleware.AuthMiddleware())
	{
		c.POST("", cluster.Create)
		c.GET("", cluster.GetList)
		c.GET("/:id/dev_space", cluster.GetSpaceList)
		c.GET("/:id/dev_space/:space_id/detail", cluster.GetSpaceDetail)
		c.GET("/:id/detail", cluster.GetDetail)
		c.DELETE("/:id", cluster.Delete)
		c.GET("/:id/storage_class", cluster.GetStorageClass)
		c.POST("/:id/storage_class", cluster.GetStorageClassByKubeConfig)
		c.PUT("/:id", cluster.Update)
	}

	// Applications
	a := g.Group("/v1/application")
	a.Use(middleware.AuthMiddleware())
	{
		a.POST("", applications.Create)
		a.GET("", applications.Get)
		a.GET("/:id", applications.GetDetail)
		a.DELETE("/:id", applications.Delete)
		a.PUT("/:id", applications.Update)
		a.PUT("/:id/dev_space/:spaceId/plugin_sync", applications.UpdateApplicationInstall)
		a.POST("/:id/bind_cluster", application_cluster.Create)
		a.GET("/:id/bound_cluster", application_cluster.GetBound)
		a.POST("/:id/create_space", cluster_user.Create)
		a.GET("/:id/dev_space", cluster_user.GetFirst)
		a.GET("/:id/dev_space/:space_id/detail", cluster_user.GetDevSpaceDetail)
		a.GET("/:id/dev_space_list", cluster_user.GetList)
		a.GET("/:id/cluster/:clusterId", applications.GetSpaceDetail)
	}

	// nocalhost
	n := g.Group("/v1/nocalhost")
	n.Use(middleware.AuthMiddleware())
	{
		n.GET("/templates", applications.GetNocalhostConfigTemplate)
	}

	// DevSpace
	dv := g.Group("v1/dev_space")
	dv.Use(middleware.AuthMiddleware())
	{
		dv.DELETE("/:id", cluster_user.Delete)
		dv.PUT("/:id", cluster_user.Update)
		dv.POST("/:id/recreate", cluster_user.ReCreate)
	}

	// Plug-in
	pa := g.Group("/v1/plugin")
	pa.Use(middleware.AuthMiddleware())
	{
		pa.GET("/applications", applications.PluginGet)
		pa.POST("/:id/recreate", cluster_user.PluginReCreate)
	}

	return g
}
