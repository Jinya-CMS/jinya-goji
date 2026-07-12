package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupAdminApiRouter(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Blog Category
	blogCategoryRouter := apiRouter.PathPrefix("/blog-category").Subrouter()
	blogCategoryRouter.HandleFunc("", GetAllBlogCategories).Methods(http.MethodGet)
	blogCategoryRouter.HandleFunc("", CreateBlogCategory).Methods(http.MethodPost)
	blogCategoryRouter.HandleFunc("/{id}", GetBlogCategoryById).Methods(http.MethodGet)
	blogCategoryRouter.HandleFunc("/{id}", UpdateBlogCategory).Methods(http.MethodPut)
	blogCategoryRouter.HandleFunc("/{id}", DeleteBlogCategory).Methods(http.MethodDelete)
	blogCategoryRouter.HandleFunc("/{id}/post", GetPostsByCategory).Methods(http.MethodGet)

	// Blog Post
	blogPostRouter := apiRouter.PathPrefix("/blog-post").Subrouter()
	blogPostRouter.HandleFunc("", GetAllBlogPosts).Methods(http.MethodGet)
	blogPostRouter.HandleFunc("", CreateBlogPost).Methods(http.MethodPost)
	blogPostRouter.HandleFunc("/{id}", GetBlogPostById).Methods(http.MethodGet)
	blogPostRouter.HandleFunc("/{id}", UpdateBlogPost).Methods(http.MethodPut)
	blogPostRouter.HandleFunc("/{id}", DeleteBlogPost).Methods(http.MethodDelete)
	blogPostRouter.HandleFunc("/{id}/section", GetBlogPostSections).Methods(http.MethodGet)
	blogPostRouter.HandleFunc("/{id}/section", UpdateBlogPostSections).Methods(http.MethodPut)

	// Artist (User)
	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("", GetAllArtists).Methods(http.MethodGet)
	userRouter.HandleFunc("", CreateArtist).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id}", GetArtistById).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", UpdateArtist).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}", DeleteArtist).Methods(http.MethodDelete)
	userRouter.HandleFunc("/{id}/activation", ActivateUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id}/totp", UpdateUserTotp).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}/profilepicture", UpdateUserProfilePicture).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}/password", UpdateUserPassword).Methods(http.MethodPut)

	// File
	fileRouter := apiRouter.PathPrefix("/file").Subrouter()
	fileRouter.HandleFunc("", GetAllFiles).Methods(http.MethodGet)
	fileRouter.HandleFunc("", CreateFile).Methods(http.MethodPost)
	fileRouter.HandleFunc("/{id}", GetFileById).Methods(http.MethodGet)
	fileRouter.HandleFunc("/{id}", UpdateFile).Methods(http.MethodPut)
	fileRouter.HandleFunc("/{id}", DeleteFile).Methods(http.MethodDelete)
	fileRouter.HandleFunc("/{id}/content", StartUpload).Methods(http.MethodPut)
	fileRouter.HandleFunc("/{id}/content/{position}", UploadChunk).Methods(http.MethodPut)
	fileRouter.HandleFunc("/{id}/content/finish", FinishUpload).Methods(http.MethodPut)

	// Gallery
	galleryRouter := apiRouter.PathPrefix("/gallery").Subrouter()
	galleryRouter.HandleFunc("", GetAllGalleries).Methods(http.MethodGet)
	galleryRouter.HandleFunc("", CreateGallery).Methods(http.MethodPost)
	galleryRouter.HandleFunc("/{id}", GetGalleryById).Methods(http.MethodGet)
	galleryRouter.HandleFunc("/{id}", UpdateGallery).Methods(http.MethodPut)
	galleryRouter.HandleFunc("/{id}", DeleteGallery).Methods(http.MethodDelete)
	galleryRouter.HandleFunc("/{id}/file", GetGalleryFiles).Methods(http.MethodGet)
	galleryRouter.HandleFunc("/{id}/file", AddFileToGallery).Methods(http.MethodPost)
	galleryRouter.HandleFunc("/{id}/file/{position}", GetGalleryFileByPosition).Methods(http.MethodGet)
	galleryRouter.HandleFunc("/{id}/file/{position}", UpdateGalleryFile).Methods(http.MethodPut)
	galleryRouter.HandleFunc("/{id}/file/{position}", DeleteGalleryFile).Methods(http.MethodDelete)

	// Classic Page
	classicPageRouter := apiRouter.PathPrefix("/classic-page").Subrouter()
	classicPageRouter.HandleFunc("", GetAllClassicPages).Methods(http.MethodGet)
	classicPageRouter.HandleFunc("", CreateClassicPage).Methods(http.MethodPost)
	classicPageRouter.HandleFunc("/{id}", GetClassicPageById).Methods(http.MethodGet)
	classicPageRouter.HandleFunc("/{id}", UpdateClassicPage).Methods(http.MethodPut)
	classicPageRouter.HandleFunc("/{id}", DeleteClassicPage).Methods(http.MethodDelete)

	// Modern Page
	modernPageRouter := apiRouter.PathPrefix("/modern-page").Subrouter()
	modernPageRouter.HandleFunc("", GetAllModernPages).Methods(http.MethodGet)
	modernPageRouter.HandleFunc("", CreateModernPage).Methods(http.MethodPost)
	modernPageRouter.HandleFunc("/{id}", GetModernPageById).Methods(http.MethodGet)
	modernPageRouter.HandleFunc("/{id}", UpdateModernPage).Methods(http.MethodPut)
	modernPageRouter.HandleFunc("/{id}", DeleteModernPage).Methods(http.MethodDelete)
	modernPageRouter.HandleFunc("/{id}/section", GetModernPageSections).Methods(http.MethodGet)
	modernPageRouter.HandleFunc("/{id}/section", UpdateModernPageSections).Methods(http.MethodPut)

	// Menu
	menuRouter := apiRouter.PathPrefix("/menu").Subrouter()
	menuRouter.HandleFunc("", GetAllMenus).Methods(http.MethodGet)
	menuRouter.HandleFunc("", CreateMenu).Methods(http.MethodPost)
	menuRouter.HandleFunc("/{id}", GetMenuById).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{id}", UpdateMenu).Methods(http.MethodPut)
	menuRouter.HandleFunc("/{id}", DeleteMenu).Methods(http.MethodDelete)
	menuRouter.HandleFunc("/{id}/item", GetMenuItems).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{id}/item", AddMenuItem).Methods(http.MethodPost)
	menuRouter.HandleFunc("/{id}/item/{position}", UpdateMenuItem).Methods(http.MethodPut)
	menuRouter.HandleFunc("/{id}/item/{position}", DeleteMenuItem).Methods(http.MethodDelete)

	// Theme
	themeRouter := apiRouter.PathPrefix("/theme").Subrouter()
	themeRouter.HandleFunc("", GetAllThemes).Methods(http.MethodGet)
	themeRouter.HandleFunc("", CreateTheme).Methods(http.MethodPost)
	themeRouter.HandleFunc("/{id}", GetThemeById).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}", UpdateTheme).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}", DeleteTheme).Methods(http.MethodDelete)
	themeRouter.HandleFunc("/{id}/active", ActivateTheme).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/assets", GetThemeAssets).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/configuration", GetThemeConfiguration).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/configuration", UpdateThemeConfiguration).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/configuration/default", GetThemeDefaultConfiguration).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/configuration/structure", GetThemeConfigurationStructure).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/styling", GetThemeStyling).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/styling", UpdateThemeStyling).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/blog-category", GetThemeBlogCategories).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/blog-category/{name}", UpdateThemeBlogCategory).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/classic-page", GetThemeClassicPages).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/classic-page/{name}", UpdateThemeClassicPage).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/file", GetThemeFiles).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/file/{name}", UpdateThemeFile).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/gallery", GetThemeGalleries).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/gallery/{name}", UpdateThemeGallery).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/menu", GetThemeMenus).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/menu/{name}", UpdateThemeMenu).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/modern-page", GetThemeModernPages).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/modern-page/{name}", UpdateThemeModernPage).Methods(http.MethodPut)

	// API Key
	apiKeyRouter := apiRouter.PathPrefix("/api-key").Subrouter()
	apiKeyRouter.HandleFunc("", GetAllApiKeys).Methods(http.MethodGet)
	apiKeyRouter.HandleFunc("/{key}", DeleteApiKey).Methods(http.MethodDelete)

	// Configuration
	configurationRouter := apiRouter.PathPrefix("/configuration").Subrouter()
	configurationRouter.HandleFunc("", GetConfiguration).Methods(http.MethodGet)
	configurationRouter.HandleFunc("", UpdateConfiguration).Methods(http.MethodPut)

	// File Tag
	fileTagRouter := apiRouter.PathPrefix("/file-tag").Subrouter()
	fileTagRouter.HandleFunc("", GetAllFileTags).Methods(http.MethodGet)
	fileTagRouter.HandleFunc("", CreateFileTag).Methods(http.MethodPost)
	fileTagRouter.HandleFunc("/{id}", GetFileTagById).Methods(http.MethodGet)
	fileTagRouter.HandleFunc("/{id}", UpdateFileTag).Methods(http.MethodPut)
	fileTagRouter.HandleFunc("/{id}", DeleteFileTag).Methods(http.MethodDelete)

	// Folder
	folderRouter := apiRouter.PathPrefix("/folder").Subrouter()
	folderRouter.HandleFunc("", GetAllFolders).Methods(http.MethodGet)
	folderRouter.HandleFunc("", CreateFolder).Methods(http.MethodPost)
	folderRouter.HandleFunc("/{id}", GetFolderById).Methods(http.MethodGet)
	folderRouter.HandleFunc("/{id}", UpdateFolder).Methods(http.MethodPut)
	folderRouter.HandleFunc("/{id}", DeleteFolder).Methods(http.MethodDelete)

	// IP Location
	ipLocationRouter := apiRouter.PathPrefix("/ip-location").Subrouter()
	ipLocationRouter.HandleFunc("{ip}", GetIpLocation).Methods(http.MethodGet)

	// Known Device
	knownDeviceRouter := apiRouter.PathPrefix("/known-device").Subrouter()
	knownDeviceRouter.HandleFunc("", GetAllKnownDevices).Methods(http.MethodGet)
	knownDeviceRouter.HandleFunc("/{key}", DeleteKnownDevice).Methods(http.MethodDelete)
	knownDeviceRouter.HandleFunc("/{key}", CheckKnownDevice).Methods(http.MethodHead)

	// Me
	meRouter := apiRouter.PathPrefix("/me").Subrouter()
	meRouter.HandleFunc("", GetMe).Methods(http.MethodGet)
	meRouter.HandleFunc("", UpdateMe).Methods(http.MethodPut)
	meRouter.HandleFunc("/colorscheme", GetMeColorScheme).Methods(http.MethodGet)
	meRouter.HandleFunc("/colorscheme", UpdateMeColorScheme).Methods(http.MethodPut)
	meRouter.HandleFunc("/preferences", GetMePreferences).Methods(http.MethodGet)
	meRouter.HandleFunc("/preferences", UpdateMePreferences).Methods(http.MethodPut)
	meRouter.HandleFunc("/profilepicture", GetMeProfilePicture).Methods(http.MethodGet)
	meRouter.HandleFunc("/profilepicture", UpdateMeProfilePicture).Methods(http.MethodPut)
	meRouter.HandleFunc("/otp/app", SetupOtpApp).Methods(http.MethodPost)
	meRouter.HandleFunc("/otp/app/verify", VerifyOtpApp).Methods(http.MethodPost)

	// Media
	mediaRouter := apiRouter.PathPrefix("/media").Subrouter()
	mediaRouter.HandleFunc("", GetMedia).Methods(http.MethodGet)

	// Statistics
	statisticsRouter := apiRouter.PathPrefix("/statistics").Subrouter()
	statisticsRouter.HandleFunc("/{group}/{months}/month", GetGroupStatisticsByMonth).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/{group}/{years}/year", GetGroupStatisticsByYear).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{interval}/month/{entityType}/{entityId}", GetEntityStatisticsByMonth).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{interval}/year/{entityType}/{entityId}", GetEntityStatisticsByYear).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{months}/month", GetVisitStatisticsByMonth).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{years}/year", GetVisitStatisticsByYear).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/{years}/year", GetStatisticsByYear).Methods(http.MethodGet)

	// Account
	accountRouter := apiRouter.PathPrefix("/account").Subrouter()
	accountRouter.HandleFunc("password", UpdatePassword).Methods(http.MethodPut)

	// Version
	versionRouter := apiRouter.PathPrefix("/version").Subrouter()
	versionRouter.HandleFunc("", GetVersion).Methods(http.MethodGet)

	// Authentication
	apiRouter.HandleFunc("/login", Login).Methods(http.MethodPost)
	apiRouter.HandleFunc("/login", CheckApiKey).Methods(http.MethodHead)
	apiRouter.HandleFunc("/logout", Logout).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/2fa", VerifyOtp).Methods(http.MethodPost)

	// Theme preview
	router.HandleFunc("/api/theme/{id}/preview", GetThemePreview).Methods(http.MethodGet)
}
