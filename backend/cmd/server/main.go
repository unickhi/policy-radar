package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"policy-radar/internal/config"
	"policy-radar/internal/handler"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
	"policy-radar/internal/service"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db := config.InitDB(cfg)
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化各层
	nationalRepo := repository.NewNationalStandardRepo(db)
	industryRepo := repository.NewIndustryStandardRepo(db)
	localRepo := repository.NewLocalStandardRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	recommendRepo := repository.NewRecommendRepo(db)
	crawlerLogRepo := repository.NewCrawlerLogRepo(db)

	nationalSvc := service.NewNationalStandardService(nationalRepo)
	industrySvc := service.NewIndustryStandardService(industryRepo)
	localSvc := service.NewLocalStandardService(localRepo)
	categorySvc := service.NewCategoryService(categoryRepo, nationalRepo, industryRepo, localRepo)
	recommendSvc := service.NewRecommendService(recommendRepo)
	crawlerSvc := service.NewCrawlerService(crawlerLogRepo, nationalRepo, industryRepo, localRepo)

	nationalHandler := handler.NewNationalStandardHandler(nationalSvc)
	industryHandler := handler.NewIndustryStandardHandler(industrySvc)
	localHandler := handler.NewLocalStandardHandler(localSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	recommendHandler := handler.NewRecommendHandler(recommendSvc)
	crawlerHandler := handler.NewCrawlerHandler(crawlerSvc)
	policyHandler := handler.NewPolicyHandler(nationalSvc, industrySvc, localSvc, recommendSvc, categorySvc)

	// 初始化路由
	r := gin.Default()

	// CORS中间件
	r.Use(config.CORS())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 国标政策
		api.GET("/national", nationalHandler.List)
		api.GET("/national/:id", nationalHandler.Get)
		api.POST("/national", nationalHandler.Create)
		api.PUT("/national/:id", nationalHandler.Update)
		api.DELETE("/national/:id", nationalHandler.Delete)
		api.POST("/national/import", nationalHandler.Import)
		api.PUT("/national/:id/check", nationalHandler.Check)

		// 行标政策
		api.GET("/industry", industryHandler.List)
		api.GET("/industry/:id", industryHandler.Get)
		api.POST("/industry", industryHandler.Create)
		api.PUT("/industry/:id", industryHandler.Update)
		api.DELETE("/industry/:id", industryHandler.Delete)
		api.POST("/industry/import", industryHandler.Import)
		api.PUT("/industry/:id/check", industryHandler.Check)

		// 地标政策
		api.GET("/local", localHandler.List)
		api.GET("/local/:id", localHandler.Get)
		api.POST("/local", localHandler.Create)
		api.PUT("/local/:id", localHandler.Update)
		api.DELETE("/local/:id", localHandler.Delete)
		api.POST("/local/import", localHandler.Import)
		api.PUT("/local/:id/check", localHandler.Check)

		// 分类管理
		api.GET("/categories", categoryHandler.List)
		api.GET("/categories/:id", categoryHandler.Get)
		api.POST("/categories", categoryHandler.Create)
		api.PUT("/categories/:id", categoryHandler.Update)
		api.DELETE("/categories/:id", categoryHandler.Delete)
		api.GET("/categories/:id/count", categoryHandler.Count)

		// 推荐政策
		api.GET("/recommends", recommendHandler.List)
		api.POST("/recommends", recommendHandler.Create)
		api.PUT("/recommends/:id", recommendHandler.Update)
		api.DELETE("/recommends/:id", recommendHandler.Delete)

		// 爬虫热更新
		api.POST("/crawler/execute", crawlerHandler.Execute)
		api.GET("/crawler/logs", crawlerHandler.Logs)
		api.POST("/crawler/import", crawlerHandler.Import)

		// 前端展示接口
		api.GET("/home", policyHandler.GetHomeData)
		api.GET("/policies/:type", policyHandler.ListByType)
		api.GET("/policies/:type/:id", policyHandler.GetDetail)
		api.GET("/search", policyHandler.Search)
	}

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})

	// 启动服务
	log.Printf("服务启动于 %s", cfg.ServerAddr)
	r.Run(cfg.ServerAddr)
}