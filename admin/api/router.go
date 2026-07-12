package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupAdminApiRouter(router *mux.Router) {
	// Blog Category
	blogCategoryRouter := router.PathPrefix("/api/blog-category").Subrouter()
	blogCategoryRouter.HandleFunc("", GetAllBlogCategories).Methods(http.MethodGet)
	blogCategoryRouter.HandleFunc("", CreateBlogCategory).Methods(http.MethodPost)
	blogCategoryRouter.HandleFunc("/{id}", GetBlogCategoryById).Methods(http.MethodGet)
	blogCategoryRouter.HandleFunc("/{id}", UpdateBlogCategory).Methods(http.MethodPut)
	blogCategoryRouter.HandleFunc("/{id}", DeleteBlogCategory).Methods(http.MethodDelete)
	blogCategoryRouter.HandleFunc("/{id}/post", GetPostsByCategory).Methods(http.MethodGet)

	// Blog Post
	blogPostRouter := router.PathPrefix("/api/blog-post").Subrouter()
	blogPostRouter.HandleFunc("", GetAllBlogPosts).Methods(http.MethodGet)
	blogPostRouter.HandleFunc("", CreateBlogPost).Methods(http.MethodPost)
	blogPostRouter.HandleFunc("/{id}", GetBlogPostById).Methods(http.MethodGet)
	blogPostRouter.HandleFunc("/{id}", UpdateBlogPost).Methods(http.MethodPut)
	blogPostRouter.HandleFunc("/{id}", DeleteBlogPost).Methods(http.MethodDelete)
	blogPostRouter.HandleFunc("/{id}/section", GetBlogPostSections).Methods(http.MethodGet)
	blogPostRouter.HandleFunc("/{id}/section", UpdateBlogPostSections).Methods(http.MethodPut)

	// Artist (User)
	userRouter := router.PathPrefix("/api/user").Subrouter()
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
	fileRouter := router.PathPrefix("/api/file").Subrouter()
	fileRouter.HandleFunc("", GetAllFiles).Methods(http.MethodGet)
	fileRouter.HandleFunc("", CreateFile).Methods(http.MethodPost)
	fileRouter.HandleFunc("/{id}", GetFileById).Methods(http.MethodGet)
	fileRouter.HandleFunc("/{id}", UpdateFile).Methods(http.MethodPut)
	fileRouter.HandleFunc("/{id}", DeleteFile).Methods(http.MethodDelete)
	fileRouter.HandleFunc("/{id}/content", StartUpload).Methods(http.MethodPut)
	fileRouter.HandleFunc("/{id}/content/{position}", UploadChunk).Methods(http.MethodPut)
	fileRouter.HandleFunc("/{id}/content/finish", FinishUpload).Methods(http.MethodPut)

	// Gallery
	galleryRouter := router.PathPrefix("/api/gallery").Subrouter()
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
	classicPageRouter := router.PathPrefix("/api/classic-page").Subrouter()
	classicPageRouter.HandleFunc("", GetAllClassicPages).Methods(http.MethodGet)
	classicPageRouter.HandleFunc("", CreateClassicPage).Methods(http.MethodPost)
	classicPageRouter.HandleFunc("/{id}", GetClassicPageById).Methods(http.MethodGet)
	classicPageRouter.HandleFunc("/{id}", UpdateClassicPage).Methods(http.MethodPut)
	classicPageRouter.HandleFunc("/{id}", DeleteClassicPage).Methods(http.MethodDelete)

	// Modern Page
	modernPageRouter := router.PathPrefix("/api/modern-page").Subrouter()
	modernPageRouter.HandleFunc("", GetAllModernPages).Methods(http.MethodGet)
	modernPageRouter.HandleFunc("", CreateModernPage).Methods(http.MethodPost)
	modernPageRouter.HandleFunc("/{id}", GetModernPageById).Methods(http.MethodGet)
	modernPageRouter.HandleFunc("/{id}", UpdateModernPage).Methods(http.MethodPut)
	modernPageRouter.HandleFunc("/{id}", DeleteModernPage).Methods(http.MethodDelete)
	modernPageRouter.HandleFunc("/{id}/section", GetModernPageSections).Methods(http.MethodGet)
	modernPageRouter.HandleFunc("/{id}/section", UpdateModernPageSections).Methods(http.MethodPut)

	// Menu
	menuRouter := router.PathPrefix("/api/menu").Subrouter()
	menuRouter.HandleFunc("", GetAllMenus).Methods(http.MethodGet)
	menuRouter.HandleFunc("", CreateMenu).Methods(http.MethodPost)
	menuRouter.HandleFunc("/{id}", GetMenuById).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{id}", UpdateMenu).Methods(http.MethodPut)
	menuRouter.HandleFunc("/{id}", DeleteMenu).Methods(http.MethodDelete)
	menuRouter.HandleFunc("/{id}/item", GetMenuItems).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{id}/item", AddMenuItem).Methods(http.MethodPost)
	menuRouter.HandleFunc("/{id}/item/{position}", UpdateMenuItem).Methods(http.MethodPut)
	menuRouter.HandleFunc("/{id}/item/{position}", DeleteMenuItem).Methods(http.MethodDelete)

	// Form
	formRouter := router.PathPrefix("/api/form").Subrouter()
	formRouter.HandleFunc("", GetAllForms).Methods(http.MethodGet)
	formRouter.HandleFunc("", CreateForm).Methods(http.MethodPost)
	formRouter.HandleFunc("/{id}", GetFormById).Methods(http.MethodGet)
	formRouter.HandleFunc("/{id}", UpdateForm).Methods(http.MethodPut)
	formRouter.HandleFunc("/{id}", DeleteForm).Methods(http.MethodDelete)
	formRouter.HandleFunc("/{id}/item", GetFormItems).Methods(http.MethodGet)
	formRouter.HandleFunc("/{id}/item", AddFormItem).Methods(http.MethodPost)
	formRouter.HandleFunc("/{id}/item/{position}", UpdateFormItem).Methods(http.MethodPut)
	formRouter.HandleFunc("/{id}/item/{position}", DeleteFormItem).Methods(http.MethodDelete)

	// Theme
	themeRouter := router.PathPrefix("/api/theme").Subrouter()
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
	themeRouter.HandleFunc("/{id}/preview", GetThemePreview).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/styling", GetThemeStyling).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/styling", UpdateThemeStyling).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/blog-category", GetThemeBlogCategories).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/blog-category/{name}", UpdateThemeBlogCategory).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/classic-page", GetThemeClassicPages).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/classic-page/{name}", UpdateThemeClassicPage).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/file", GetThemeFiles).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/file/{name}", UpdateThemeFile).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/form", GetThemeForms).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/form/{name}", UpdateThemeForm).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/gallery", GetThemeGalleries).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/gallery/{name}", UpdateThemeGallery).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/menu", GetThemeMenus).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/menu/{name}", UpdateThemeMenu).Methods(http.MethodPut)
	themeRouter.HandleFunc("/{id}/modern-page", GetThemeModernPages).Methods(http.MethodGet)
	themeRouter.HandleFunc("/{id}/modern-page/{name}", UpdateThemeModernPage).Methods(http.MethodPut)

	// API Key
	apiKeyRouter := router.PathPrefix("/api/api-key").Subrouter()
	apiKeyRouter.HandleFunc("", GetAllApiKeys).Methods(http.MethodGet)
	apiKeyRouter.HandleFunc("/{key}", DeleteApiKey).Methods(http.MethodDelete)

	// Configuration
	configurationRouter := router.PathPrefix("/api/configuration").Subrouter()
	configurationRouter.HandleFunc("", GetConfiguration).Methods(http.MethodGet)
	configurationRouter.HandleFunc("", UpdateConfiguration).Methods(http.MethodPut)

	// File Tag
	fileTagRouter := router.PathPrefix("/api/file-tag").Subrouter()
	fileTagRouter.HandleFunc("", GetAllFileTags).Methods(http.MethodGet)
	fileTagRouter.HandleFunc("", CreateFileTag).Methods(http.MethodPost)
	fileTagRouter.HandleFunc("/{id}", GetFileTagById).Methods(http.MethodGet)
	fileTagRouter.HandleFunc("/{id}", UpdateFileTag).Methods(http.MethodPut)
	fileTagRouter.HandleFunc("/{id}", DeleteFileTag).Methods(http.MethodDelete)

	// Folder
	folderRouter := router.PathPrefix("/api/folder").Subrouter()
	folderRouter.HandleFunc("", GetAllFolders).Methods(http.MethodGet)
	folderRouter.HandleFunc("", CreateFolder).Methods(http.MethodPost)
	folderRouter.HandleFunc("/{id}", GetFolderById).Methods(http.MethodGet)
	folderRouter.HandleFunc("/{id}", UpdateFolder).Methods(http.MethodPut)
	folderRouter.HandleFunc("/{id}", DeleteFolder).Methods(http.MethodDelete)

	// IP Location
	ipLocationRouter := router.PathPrefix("/api/ip-location").Subrouter()
	ipLocationRouter.HandleFunc("{ip}", GetIpLocation).Methods(http.MethodGet)

	// Known Device
	knownDeviceRouter := router.PathPrefix("/api/known-device").Subrouter()
	knownDeviceRouter.HandleFunc("", GetAllKnownDevices).Methods(http.MethodGet)
	knownDeviceRouter.HandleFunc("/{key}", DeleteKnownDevice).Methods(http.MethodDelete)
	knownDeviceRouter.HandleFunc("/{key}", CheckKnownDevice).Methods(http.MethodHead)

	// Maintenance
	maintenanceDatabaseRouter := router.PathPrefix("/api/maintenance/database").Subrouter()
	maintenanceDatabaseRouter.HandleFunc("analyze", DatabaseAnalyze).Methods(http.MethodPost)
	maintenanceDatabaseRouter.HandleFunc("query", DatabaseQuery).Methods(http.MethodPost)

	// Me
	meRouter := router.PathPrefix("/api/me").Subrouter()
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
	mediaRouter := router.PathPrefix("/api/media").Subrouter()
	mediaRouter.HandleFunc("", GetMedia).Methods(http.MethodGet)

	// Statistics
	statisticsRouter := router.PathPrefix("/api/statistics").Subrouter()
	statisticsRouter.HandleFunc("/{group}/{months}/month", GetGroupStatisticsByMonth).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/{group}/{years}/year", GetGroupStatisticsByYear).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{interval}/month/{entityType}/{entityId}", GetEntityStatisticsByMonth).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{interval}/year/{entityType}/{entityId}", GetEntityStatisticsByYear).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{months}/month", GetVisitStatisticsByMonth).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/visits/{years}/year", GetVisitStatisticsByYear).Methods(http.MethodGet)
	statisticsRouter.HandleFunc("/{years}/year", GetStatisticsByYear).Methods(http.MethodGet)

	// Account
	accountRouter := router.PathPrefix("/api/account").Subrouter()
	accountRouter.HandleFunc("password", UpdatePassword).Methods(http.MethodPut)

	// Version
	versionRouter := router.PathPrefix("/api/version").Subrouter()
	versionRouter.HandleFunc("", GetVersion).Methods(http.MethodGet)

	// Authentication
	router.HandleFunc("/api/login", Login).Methods(http.MethodPost)
	router.HandleFunc("/api/login", CheckApiKey).Methods(http.MethodHead)
	router.HandleFunc("/api/logout", Logout).Methods(http.MethodDelete)
	router.HandleFunc("/api/2fa", TwoFactorEntry).Methods(http.MethodPost)
}
