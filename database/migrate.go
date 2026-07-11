package database

import (
	"github.com/DerKnerd/gorp"
)

var dbMap *gorp.DbMap

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func CreateDbMap() {
	if dbMap == nil {
		dialect := gorp.MySQLDialect{}

		dbMap = &gorp.DbMap{Dialect: dialect}

		AddTableWithName[AnalyticsEntry]("analytics")
		AddTableWithName[ApiKey]("api_key")
		AddTableWithName[Artist]("users")
		AddTableWithName[BlogCategory]("blog_category")
		AddTableWithName[BlogPost]("blog_post")
		AddTableWithName[BlogPostSection]("blog_post_segment")
		AddTableWithName[ClassicPage]("page")
		AddTableWithName[File]("file")
		AddTableWithName[FileTag]("file_tag")
		AddTableWithName[Folder]("folder")
		AddTableWithName[Form]("form")
		AddTableWithName[FormItem]("form_item")
		AddTableWithName[Gallery]("gallery")
		AddTableWithName[GalleryFilePosition]("gallery_file_position")
		AddTableWithName[KnownDevice]("known_device")
		AddTableWithName[Menu]("menu")
		AddTableWithName[MenuItem]("menu_item")
		AddTableWithName[ModernPage]("segment_page")
		AddTableWithName[ModernPageSection]("segment")
		AddTableWithName[Theme]("theme")
		AddTableWithName[ThemeAsset]("theme_asset")
		AddTableWithName[ThemeBlogCategory]("theme_blog_category")
		AddTableWithName[ThemeClassicPage]("theme_page")
		AddTableWithName[ThemeFile]("theme_file")
		AddTableWithName[ThemeForm]("theme_form")
		AddTableWithName[ThemeGallery]("theme_gallery")
		AddTableWithName[ThemeMenu]("theme_menu")
		AddTableWithName[ThemeModernPage]("theme_segment_page")
		AddTableWithName[UploadingFile]("uploading_file")
		AddTableWithName[UploadingFileChunk]("uploading_file_chunk")
	}
}
