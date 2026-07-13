package database

import (
	"context"
	"os"

	"github.com/DerKnerd/gorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"jinya.de/goji/database/migrations"
)

var dbMap *gorp.DbMap

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func SetupDatabase() {
	if dbMap == nil {
		pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}

		conn := stdlib.OpenDBFromPool(pool)

		dialect := gorp.PostgresDialect{}

		dbMap = &gorp.DbMap{Db: conn, Dialect: dialect}

		AddTableWithName[AnalyticsEntry]("analytics")
		AddTableWithName[ApiKey]("api_key")
		AddTableWithName[Artist]("users")
		AddTableWithName[BlogCategory]("blog_category")
		AddTableWithName[BlogPost]("blog_post")
		AddTableWithName[BlogPostSection]("blog_post_section")
		AddTableWithName[ClassicPage]("classic_page")
		AddTableWithName[File]("file")
		AddTableWithName[FileTag]("file_tag")
		AddTableWithName[Folder]("folder")
		AddTableWithName[Gallery]("gallery")
		AddTableWithName[GalleryFilePosition]("gallery_file_position")
		AddTableWithName[KnownDevice]("known_device")
		AddTableWithName[Menu]("menu")
		AddTableWithName[MenuItem]("menu_item")
		AddTableWithName[ModernPage]("modern_page")
		AddTableWithName[ModernPageSection]("modern_page_section")
		AddTableWithName[Theme]("theme")
		AddTableWithName[ThemeAsset]("theme_asset")
		AddTableWithName[ThemeBlogCategory]("theme_blog_category")
		AddTableWithName[ThemeClassicPage]("theme_page")
		AddTableWithName[ThemeFile]("theme_file")
		AddTableWithName[ThemeGallery]("theme_gallery")
		AddTableWithName[ThemeMenu]("theme_menu")
		AddTableWithName[ThemeModernPage]("theme_modern_page")
		AddTableWithName[UploadingFile]("uploading_file")
		AddTableWithName[UploadingFileChunk]("uploading_file_chunk")
		AddTableWithName[JinyaConfiguration]("jinya_configuration")

		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		migrations.BaseMigration(conn)
		migrations.CreateAdditionalTables(conn)
		migrations.CreateForeignKeys(conn)
	}
}
