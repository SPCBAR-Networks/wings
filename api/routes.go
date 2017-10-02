package api

func (api *InternalAPI) RegisterRoutes() {
	// Register routes for v1 of the API. This API should be fully backwards compatable with
	// the existing Nodejs Daemon API.
	v1 := api.router.Group("/v1")
	{
		v1.GET("/", AuthHandler(""), GetIndex)
		v1.GET("/servers", AuthHandler("c:list"), handleGetServers)

		v1.PATCH("/config", AuthHandler("c:config"), PatchConfiguration)
		v1.POST("/servers", AuthHandler("c:create"), handlePostServers)

		v1ServerRoutes := v1.Group("/servers/:server")
		{
			v1ServerRoutes.GET("/", AuthHandler("s:get"), handleGetServer)
			v1ServerRoutes.GET("/log", AuthHandler("s:console"), handleGetServerLog)

			v1ServerRoutes.POST("/reinstall", AuthHandler("s:install-server"), handlePostServerReinstall)
			v1ServerRoutes.POST("/rebuild", AuthHandler("g:server:rebuild"), handlePostServerRebuild)
			v1ServerRoutes.POST("/password", AuthHandler(""), handlePostServerPassword)
			v1ServerRoutes.POST("/power", AuthHandler("s:power"), handlePostServerPower)
			v1ServerRoutes.POST("/command", AuthHandler("s:command"), handlePostServerCommand)
			v1ServerRoutes.POST("/suspend", AuthHandler(""), handlePostServerSuspend)
			v1ServerRoutes.POST("/unsuspend", AuthHandler(""), handlePostServerUnsuspend)

			v1ServerRoutes.PATCH("/", AuthHandler("s:config"), handlePatchServer)
			v1ServerRoutes.DELETE("/", AuthHandler("g:server:delete"), handleDeleteServer)
		}

		v1ServerFileRoutes := v1.Group("/servers/:server/files")
		{
			v1ServerFileRoutes.GET("/file/:file", AuthHandler("s:files:read"), ReadFileContents)
			v1ServerFileRoutes.GET("/stat/:file", AuthHandler("s:files:get"), StatFile)
			v1ServerFileRoutes.GET("/directory/:directory", AuthHandler("s:files:get"), ListDirectory)
			v1ServerFileRoutes.GET("/download/:token", DownloadFile)

			v1ServerFileRoutes.POST("/directory/:directory", AuthHandler("s:files:create"), StoreDirectory)
			v1ServerFileRoutes.POST("/file/:file", AuthHandler("s:files:post"), WriteFileContents)
			v1ServerFileRoutes.POST("/copy/:file", AuthHandler("s:files:copy"), CopyFile)
			v1ServerFileRoutes.POST("/move/:file", AuthHandler("s:files:move"), MoveFile)
			v1ServerFileRoutes.POST("/rename/:file", AuthHandler("s:files:move"), MoveFile)
			v1ServerFileRoutes.POST("/compress/:file", AuthHandler("s:files:compress"), CompressFile)
			v1ServerFileRoutes.POST("/decompress/:file", AuthHandler("s:files:decompress"), DecompressFile)

			v1ServerFileRoutes.DELETE("/file/:file", AuthHandler("s:files:delete"), DeleteFile)
		}
	}
}
