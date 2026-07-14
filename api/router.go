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
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllBlogCategories))).
		Methods(http.MethodGet)
	blogCategoryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createBlogCategory))).
		Methods(http.MethodPost)
	blogCategoryRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(getBlogCategoryById))).
		Methods(http.MethodGet)
	blogCategoryRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateBlogCategory))).
		Methods(http.MethodPut)
	blogCategoryRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteBlogCategory))).
		Methods(http.MethodDelete)
	blogCategoryRouter.
		Path("/{id:\\d+}/post").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getPostsByCategory))).
		Methods(http.MethodGet)

	// Blog Post
	blogPostRouter := apiRouter.PathPrefix("/blog-post").Subrouter()
	blogPostRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllBlogPosts))).
		Methods(http.MethodGet)
	blogPostRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createBlogPost))).
		Methods(http.MethodPost)
	blogPostRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getBlogPostById))).
		Methods(http.MethodGet)
	blogPostRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateBlogPost))).
		Methods(http.MethodPut)
	blogPostRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteBlogPost))).
		Methods(http.MethodDelete)
	blogPostRouter.
		Path("/{id:\\d+}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getBlogPostSections))).
		Methods(http.MethodGet)
	blogPostRouter.
		Path("/{id:\\d+}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateBlogPostSections))).
		Methods(http.MethodPut)

	// Artist (User)
	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllArtists))).
		Methods(http.MethodGet)
	userRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(createArtist))).
		Methods(http.MethodPost)
	userRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getArtistById))).
		Methods(http.MethodGet)
	userRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(updateArtist))).
		Methods(http.MethodPut)
	userRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(deleteArtist))).
		Methods(http.MethodDelete)
	userRouter.
		Path("/{id:\\d+}/activation").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(activateUser))).
		Methods(http.MethodPut)
	userRouter.
		Path("/{id:\\d+}/activation").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(deactivateUser))).
		Methods(http.MethodDelete)
	userRouter.
		Path("/{id:\\d+}/profilepicture").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(updateUserProfilePicture))).
		Methods(http.MethodPut)
	userRouter.
		Path("/{id:\\d+}/profilepicture").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(deleteUserProfilePicture))).
		Methods(http.MethodDelete)
	userRouter.
		Path("/{id:\\d+}/totp").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(resetTotp))).
		Methods(http.MethodPut)

	// File
	fileRouter := apiRouter.PathPrefix("/file").Subrouter()
	fileRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createFile))).
		Methods(http.MethodPost)
	fileRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(getFile))).
		Methods(http.MethodGet)
	fileRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateFile))).
		Methods(http.MethodPut)
	fileRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteFile))).
		Methods(http.MethodDelete)
	fileRouter.
		Path("/{id:\\d+}/content").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(startUpload))).
		Methods(http.MethodPut)
	fileRouter.
		Path("/{id:\\d+}/content/{position:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(uploadChunk))).
		Methods(http.MethodPut)
	fileRouter.
		Path("/{id:\\d+}/content/finish").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(finishUpload))).
		Methods(http.MethodPut)

	// Gallery
	galleryRouter := apiRouter.PathPrefix("/gallery").Subrouter()
	galleryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllGalleries))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createGallery))).
		Methods(http.MethodPost)
	galleryRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getGalleryById))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateGallery))).
		Methods(http.MethodPut)
	galleryRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteGallery))).
		Methods(http.MethodDelete)
	galleryRouter.
		Path("/{id:\\d+}/file").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getGalleryFiles))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("/{id:\\d+}/file").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(addFileToGallery))).
		Methods(http.MethodPost)
	galleryRouter.
		Path("/{id:\\d+}/file/{position:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getGalleryFileByPosition))).
		Methods(http.MethodGet)
	galleryRouter.
		Path("/{id:\\d+}/file/{position:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateGalleryFile))).
		Methods(http.MethodPut)
	galleryRouter.
		Path("/{id:\\d+}/file/{position:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteGalleryFile))).
		Methods(http.MethodDelete)

	// Classic Page
	classicPageRouter := apiRouter.PathPrefix("/classic-page").Subrouter()
	classicPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllClassicPages))).
		Methods(http.MethodGet)
	classicPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createClassicPage))).
		Methods(http.MethodPost)
	classicPageRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getClassicPageById))).
		Methods(http.MethodGet)
	classicPageRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateClassicPage))).
		Methods(http.MethodPut)
	classicPageRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteClassicPage))).
		Methods(http.MethodDelete)

	// Modern Page
	modernPageRouter := apiRouter.PathPrefix("/modern-page").Subrouter()
	modernPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllModernPages))).
		Methods(http.MethodGet)
	modernPageRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createModernPage))).
		Methods(http.MethodPost)
	modernPageRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getModernPageById))).
		Methods(http.MethodGet)
	modernPageRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateModernPage))).
		Methods(http.MethodPut)
	modernPageRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteModernPage))).
		Methods(http.MethodDelete)
	modernPageRouter.
		Path("/{id:\\d+}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getModernPageSections))).
		Methods(http.MethodGet)
	modernPageRouter.
		Path("/{id:\\d+}/section").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateModernPageSections))).
		Methods(http.MethodPut)

	// Menu
	menuRouter := apiRouter.PathPrefix("/menu").Subrouter()
	menuRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllMenus))).
		Methods(http.MethodGet)
	menuRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createMenu))).
		Methods(http.MethodPost)
	menuRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getMenuById))).
		Methods(http.MethodGet)
	menuRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateMenu))).
		Methods(http.MethodPut)
	menuRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteMenu))).
		Methods(http.MethodDelete)
	menuRouter.
		Path("/{id:\\d+}/item").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getMenuItems))).
		Methods(http.MethodGet)
	menuRouter.
		Path("/{id:\\d+}/item").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(addMenuItem))).
		Methods(http.MethodPost)
	menuRouter.
		Path("/{id:\\d+}/item/{position:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateMenuItem))).
		Methods(http.MethodPut)
	menuRouter.
		Path("/{id:\\d+}/item/{position:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteMenuItem))).
		Methods(http.MethodDelete)

	// Theme
	themeRouter := apiRouter.PathPrefix("/theme").Subrouter()
	themeRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getAllThemes))).
		Methods(http.MethodGet)
	themeRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createTheme))).
		Methods(http.MethodPost)
	themeRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeById))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateTheme))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteTheme))).
		Methods(http.MethodDelete)
	themeRouter.
		Path("/{id:\\d+}/active").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(activateTheme))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/assets").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeAssets))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/configuration").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeConfiguration))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/configuration").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeConfiguration))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/configuration/default").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeDefaultConfiguration))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/configuration/structure").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeConfigurationStructure))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/styling").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeStyling))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/styling").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeStyling))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/blog-category").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeBlogCategories))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/blog-category/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeBlogCategory))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/classic-page").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeClassicPages))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/classic-page/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeClassicPage))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/file").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeFiles))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/file/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeFile))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/gallery").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeGalleries))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/gallery/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeGallery))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/menu").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeMenus))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/menu/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeMenu))).
		Methods(http.MethodPut)
	themeRouter.
		Path("/{id:\\d+}/modern-page").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getThemeModernPages))).
		Methods(http.MethodGet)
	themeRouter.
		Path("/{id:\\d+}/modern-page/{name}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateThemeModernPage))).
		Methods(http.MethodPut)

	// API Key
	apiKeyRouter := apiRouter.PathPrefix("/api-key").Subrouter()
	apiKeyRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(getAllApiKeys))).
		Methods(http.MethodGet)
	apiKeyRouter.
		Path("/{key}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(deleteApiKey))).
		Methods(http.MethodDelete)

	// Configuration
	configurationRouter := apiRouter.PathPrefix("/configuration").Subrouter()
	configurationRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(getConfiguration))).
		Methods(http.MethodGet)
	configurationRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(updateConfiguration))).
		Methods(http.MethodPut)

	// File Tag
	fileTagRouter := apiRouter.PathPrefix("/file-tag").Subrouter()
	fileTagRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createFileTag))).
		Methods(http.MethodPost)
	fileTagRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateFileTag))).
		Methods(http.MethodPut)
	fileTagRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteFileTag))).
		Methods(http.MethodDelete)

	// Folder
	folderRouter := apiRouter.PathPrefix("/folder").Subrouter()
	folderRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(createFolder))).
		Methods(http.MethodPost)
	folderRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(updateFolder))).
		Methods(http.MethodPut)
	folderRouter.
		Path("/{id:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleWriter)(http.HandlerFunc(deleteFolder))).
		Methods(http.MethodDelete)

	// IP Location
	ipLocationRouter := apiRouter.PathPrefix("/ip-location").Subrouter()
	ipLocationRouter.
		Path("{ip}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(getIpLocation))).
		Methods(http.MethodGet)

	// Known Device
	knownDeviceRouter := apiRouter.PathPrefix("/known-device").Subrouter()
	knownDeviceRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(getAllKnownDevices))).
		Methods(http.MethodGet)
	knownDeviceRouter.
		Path("/{key}").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(deleteKnownDevice))).
		Methods(http.MethodDelete)
	knownDeviceRouter.
		Path("/{key}").
		Handler(http.HandlerFunc(checkKnownDevice)).
		Methods(http.MethodHead)

	// Me
	meRouter := apiRouter.PathPrefix("/me").Subrouter()
	meRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(getMe))).
		Methods(http.MethodGet)
	meRouter.
		Path("").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(updateMe))).
		Methods(http.MethodPut)
	meRouter.
		Path("/colorscheme").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(updateMeColorScheme))).
		Methods(http.MethodPut)
	meRouter.
		Path("/preferences").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(updateMePreferences))).
		Methods(http.MethodPut)
	meRouter.
		Path("/profilepicture").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(updateMeProfilePicture))).
		Methods(http.MethodPut)
	meRouter.
		Path("/otp/app").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(setupOtpApp))).
		Methods(http.MethodPost)
	meRouter.
		Path("/otp/app/verify").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(verifyOtpApp))).
		Methods(http.MethodPost)

	// Media
	mediaRouter := apiRouter.PathPrefix("/media").Subrouter()
	mediaRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getMedia))).
		Methods(http.MethodGet)

	// Statistics
	statisticsRouter := apiRouter.PathPrefix("/statistics").Subrouter()
	statisticsRouter.
		Path("/{group:country|os|browser|os-version|browser-version|referer|brand|type|language|date}/{months:[1-9]\\d*}/month").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getGroupStatisticsByMonth))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/{group:country|os|browser|os-version|browser-version|referer|brand|type|language|date}/{years:[1-9]\\d*}/year").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getGroupStatisticsByYear))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{months:[1-9]\\d*}/month/{entityType:blog-post|blog-category|gallery|modern-page|classic-page|artist}/{entityId:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getEntityStatisticsByMonth))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{years:[1-9]\\d*}/year/{entityType:blog-post|blog-category|gallery|modern-page|classic-page|artist}/{entityId:\\d+}").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getEntityStatisticsByYear))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{months:[1-9]\\d*}/month").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getVisitStatisticsByMonth))).
		Methods(http.MethodGet)
	statisticsRouter.
		Path("/visits/{years:[1-9]\\d*}/year").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleReader)(http.HandlerFunc(getVisitStatisticsByYear))).
		Methods(http.MethodGet)

	// Account
	accountRouter := apiRouter.PathPrefix("/account").Subrouter()
	accountRouter.
		Path("/password").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(updatePassword))).
		Methods(http.MethodPut)

	// Version
	versionRouter := apiRouter.PathPrefix("/version").Subrouter()
	versionRouter.
		Path("").
		Handler(middleware.CheckAuthenticationAndRole(auth.RoleAdmin)(http.HandlerFunc(getVersion))).
		Methods(http.MethodGet)

	// Authentication
	apiRouter.
		Path("/login").
		HandlerFunc(login).
		Methods(http.MethodPost)
	apiRouter.
		Path("/login").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(checkApiKey))).
		Methods(http.MethodHead)
	apiRouter.
		Path("/logout").
		HandlerFunc(logout).
		Methods(http.MethodDelete)
	apiRouter.
		Path("/2fa").
		HandlerFunc(request2fa).
		Methods(http.MethodPost)

	// Theme preview
	router.
		Path("/api/theme/{id:\\d+}/preview").
		Handler(middleware.CheckAuthentication(http.HandlerFunc(GetThemePreview))).
		Methods(http.MethodGet)

	apiRouter.Use(middleware.ContentTypeJson)

	router.Use(localization.SetLanguageMiddleware)
}
