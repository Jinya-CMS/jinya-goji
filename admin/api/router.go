package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/auth"
	"jinya.de/goji/localization"
)

func SetupAdminApiRouter(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Blog Category
	blogCategoryRouter := apiRouter.PathPrefix("/blog-category").Subrouter()
	blogCategoryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllBlogCategories))).
		Methods(http.MethodGet)
	blogCategoryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateBlogCategory))).
		Methods(http.MethodPost)
	blogCategoryRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(GetBlogCategoryById))).
		Methods(http.MethodGet)
	blogCategoryRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateBlogCategory))).
		Methods(http.MethodPut)
	blogCategoryRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteBlogCategory))).
		Methods(http.MethodDelete)
	blogCategoryRouter.
		Path("/{id}/post").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetPostsByCategory))).
		Methods(http.MethodGet)

	// Blog Post
	blogPostRouter := apiRouter.PathPrefix("/blog-post").Subrouter()
	blogPostRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllBlogPosts))).
		Methods(http.MethodGet)
	blogPostRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateBlogPost))).
		Methods(http.MethodPost)
	blogPostRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetBlogPostById))).
		Methods(http.MethodGet)
	blogPostRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateBlogPost))).
		Methods(http.MethodPut)
	blogPostRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteBlogPost))).
		Methods(http.MethodDelete)
	blogPostRouter.
		Path("/{id}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetBlogPostSections))).
		Methods(http.MethodGet)
	blogPostRouter.
		Path("/{id}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateBlogPostSections))).
		Methods(http.MethodPut)

	// Artist (User)
	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllArtists))).
		Methods(http.MethodGet)
	userRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(CreateArtist))).
		Methods(http.MethodPost)
	userRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetArtistById))).
		Methods(http.MethodGet)
	userRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(UpdateArtist))).
		Methods(http.MethodPut)
	userRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(DeleteArtist))).
		Methods(http.MethodDelete)
	userRouter.
		Path("/{id}/activation").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(ActivateUser))).
		Methods(http.MethodPost)
	userRouter.
		Path("/{id}/activation").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(DeactivateUser))).
		Methods(http.MethodDelete)
	userRouter.
		Path("/{id}/profilepicture").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(UpdateUserProfilePicture))).
		Methods(http.MethodPut)
	userRouter.
		Path("/{id}/profilepicture").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(DeleteUserProfilePicture))).
		Methods(http.MethodDelete)
	userRouter.
		Path("/{id}/totp").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(ResetTotp))).
		Methods(http.MethodPut)

	// File
	fileRouter := apiRouter.PathPrefix("/file").Subrouter()
	fileRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllFiles))).
		Methods(http.MethodGet)
	fileRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateFile))).
		Methods(http.MethodPost)
	fileRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetFileById))).
		Methods(http.MethodGet)
	fileRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateFile))).
		Methods(http.MethodPut)
	fileRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteFile))).
		Methods(http.MethodDelete)
	fileRouter.
		Path("/{id}/content").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(StartUpload))).
		Methods(http.MethodPut)
	fileRouter.
		Path("/{id}/content/{position}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UploadChunk))).
		Methods(http.MethodPut)
	fileRouter.
		Path("/{id}/content/finish").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(FinishUpload))).
		Methods(http.MethodPut)

	// Gallery
	galleryRouter := apiRouter.PathPrefix("/gallery").Subrouter()
	galleryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllGalleries))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateGallery))).
		Methods(http.MethodPost)
	galleryRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetGalleryById))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateGallery))).
		Methods(http.MethodPut)
	galleryRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteGallery))).
		Methods(http.MethodDelete)
	galleryRouter.
		Path("/{id}/file").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetGalleryFiles))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("/{id}/file").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(AddFileToGallery))).
		Methods(http.MethodPost)
	galleryRouter.
		Path("/{id}/file/{position}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetGalleryFileByPosition))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("/{id}/file/{position}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateGalleryFile))).
		Methods(http.MethodPut)
	galleryRouter.
		Path("/{id}/file/{position}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteGalleryFile))).
		Methods(http.MethodDelete)

	// Classic Page
	classicPageRouter := apiRouter.PathPrefix("/classic-page").Subrouter()
	classicPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllClassicPages))).
		Methods(http.MethodGet)
	classicPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateClassicPage))).
		Methods(http.MethodPost)
	classicPageRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetClassicPageById))).
		Methods(http.MethodGet)
	classicPageRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateClassicPage))).
		Methods(http.MethodPut)
	classicPageRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteClassicPage))).
		Methods(http.MethodDelete)

	// Modern Page
	modernPageRouter := apiRouter.PathPrefix("/modern-page").Subrouter()
	modernPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllModernPages))).
		Methods(http.MethodGet)
	modernPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateModernPage))).
		Methods(http.MethodPost)
	modernPageRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetModernPageById))).
		Methods(http.MethodGet)
	modernPageRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateModernPage))).
		Methods(http.MethodPut)
	modernPageRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteModernPage))).
		Methods(http.MethodDelete)
	modernPageRouter.
		Path("/{id}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetModernPageSections))).
		Methods(http.MethodGet)
	modernPageRouter.
		Path("/{id}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateModernPageSections))).
		Methods(http.MethodPut)

	// Menu
	menuRouter := apiRouter.PathPrefix("/menu").Subrouter()
	menuRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllMenus))).
		Methods(http.MethodGet)
	menuRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateMenu))).
		Methods(http.MethodPost)
	menuRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetMenuById))).
		Methods(http.MethodGet)
	menuRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateMenu))).
		Methods(http.MethodPut)
	menuRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteMenu))).
		Methods(http.MethodDelete)
	menuRouter.
		Path("/{id}/item").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetMenuItems))).
		Methods(http.MethodGet)
	menuRouter.
		Path("/{id}/item").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(AddMenuItem))).
		Methods(http.MethodPost)
	menuRouter.
		Path("/{id}/item/{position}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateMenuItem))).
		Methods(http.MethodPut)
	menuRouter.
		Path("/{id}/item/{position}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteMenuItem))).
		Methods(http.MethodDelete)

	// Theme
	themeRouter := apiRouter.PathPrefix("/theme").Subrouter()
	themeRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllThemes))).
		Methods(http.MethodGet)
	themeRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateTheme))).
		Methods(http.MethodPost)
	themeRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeById))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateTheme))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteTheme))).
		Methods(http.MethodDelete)
	themeRouter.
		Path("/{id}/active").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(ActivateTheme))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/assets").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeAssets))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/configuration").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeConfiguration))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/configuration").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeConfiguration))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/configuration/default").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeDefaultConfiguration))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/configuration/structure").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeConfigurationStructure))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/styling").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeStyling))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/styling").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeStyling))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/blog-category").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeBlogCategories))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/blog-category/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeBlogCategory))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/classic-page").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeClassicPages))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/classic-page/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeClassicPage))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/file").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeFiles))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/file/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeFile))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/gallery").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeGalleries))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/gallery/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeGallery))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/menu").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeMenus))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/menu/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeMenu))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id}/modern-page").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetThemeModernPages))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id}/modern-page/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateThemeModernPage))).
		Methods(http.MethodPut)

	// API Key
	apiKeyRouter := apiRouter.PathPrefix("/api-key").Subrouter()
	apiKeyRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetAllApiKeys))).
		Methods(http.MethodGet)
	apiKeyRouter.
		Path("/{key}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(DeleteApiKey))).
		Methods(http.MethodDelete)

	// Configuration
	configurationRouter := apiRouter.PathPrefix("/configuration").Subrouter()
	configurationRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(GetConfiguration))).
		Methods(http.MethodGet)
	configurationRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(UpdateConfiguration))).
		Methods(http.MethodPut)

	// File Tag
	fileTagRouter := apiRouter.PathPrefix("/file-tag").Subrouter()
	fileTagRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllFileTags))).
		Methods(http.MethodGet)
	fileTagRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateFileTag))).
		Methods(http.MethodPost)
	fileTagRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetFileTagById))).
		Methods(http.MethodGet)
	fileTagRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateFileTag))).
		Methods(http.MethodPut)
	fileTagRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteFileTag))).
		Methods(http.MethodDelete)

	// Folder
	folderRouter := apiRouter.PathPrefix("/folder").Subrouter()
	folderRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetAllFolders))).
		Methods(http.MethodGet)
	folderRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(CreateFolder))).
		Methods(http.MethodPost)
	folderRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetFolderById))).
		Methods(http.MethodGet)
	folderRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(UpdateFolder))).
		Methods(http.MethodPut)
	folderRouter.
		Path("/{id}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(DeleteFolder))).
		Methods(http.MethodDelete)

	// IP Location
	ipLocationRouter := apiRouter.PathPrefix("/ip-location").Subrouter()
	ipLocationRouter.
		Path("{ip}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetIpLocation))).
		Methods(http.MethodGet)

	// Known Device
	knownDeviceRouter := apiRouter.PathPrefix("/known-device").Subrouter()
	knownDeviceRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetAllKnownDevices))).
		Methods(http.MethodGet)
	knownDeviceRouter.
		Path("/{key}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(DeleteKnownDevice))).
		Methods(http.MethodDelete)
	knownDeviceRouter.
		Path("/{key}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(CheckKnownDevice))).
		Methods(http.MethodHead)

	// Me
	meRouter := apiRouter.PathPrefix("/me").Subrouter()
	meRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetMe))).
		Methods(http.MethodGet)
	meRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(UpdateMe))).
		Methods(http.MethodPut)
	meRouter.
		Path("/colorscheme").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetMeColorScheme))).
		Methods(http.MethodGet)
	meRouter.
		Path("/colorscheme").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(UpdateMeColorScheme))).
		Methods(http.MethodPut)
	meRouter.
		Path("/preferences").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetMePreferences))).
		Methods(http.MethodGet)
	meRouter.
		Path("/preferences").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(UpdateMePreferences))).
		Methods(http.MethodPut)
	meRouter.
		Path("/profilepicture").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetMeProfilePicture))).
		Methods(http.MethodGet)
	meRouter.
		Path("/profilepicture").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(UpdateMeProfilePicture))).
		Methods(http.MethodPut)
	meRouter.
		Path("/otp/app").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(SetupOtpApp))).
		Methods(http.MethodPost)
	meRouter.
		Path("/otp/app/verify").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(VerifyOtpApp))).
		Methods(http.MethodPost)

	// Media
	mediaRouter := apiRouter.PathPrefix("/media").Subrouter()
	mediaRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetMedia))).
		Methods(http.MethodGet)

	// Statistics
	statisticsRouter := apiRouter.PathPrefix("/statistics").Subrouter()
	statisticsRouter.
		Path("/{group}/{months}/month").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetGroupStatisticsByMonth))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/{group}/{years}/year").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetGroupStatisticsByYear))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{interval}/month/{entityType}/{entityId}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetEntityStatisticsByMonth))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{interval}/year/{entityType}/{entityId}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetEntityStatisticsByYear))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{months}/month").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetVisitStatisticsByMonth))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{years}/year").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetVisitStatisticsByYear))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/{years}/year").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(GetStatisticsByYear))).
		Methods(http.MethodGet)

	// Account
	accountRouter := apiRouter.PathPrefix("/account").Subrouter()
	accountRouter.
		Path("/password").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(UpdatePassword))).
		Methods(http.MethodPut)

	// Version
	versionRouter := apiRouter.PathPrefix("/version").Subrouter()
	versionRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(GetVersion))).
		Methods(http.MethodGet)

	// Authentication
	apiRouter.
		Path("/login").
		HandlerFunc(Login).
		Methods(http.MethodPost)
	apiRouter.
		Path("/login").
		HandlerFunc(CheckApiKey).
		Methods(http.MethodHead)
	apiRouter.
		Path("/logout").
		HandlerFunc(Logout).
		Methods(http.MethodDelete)
	apiRouter.
		Path("/2fa").
		HandlerFunc(VerifyOtp).
		Methods(http.MethodPost)

	// Theme preview
	router.
		Path("/api/theme/{id}/preview").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetThemePreview))).
		Methods(http.MethodGet)

	apiRouter.Use(middleware.ContentTypeJson)

	router.Use(localization.SetLanguageMiddleware)
}
