package database

import (
	"context"
	"os"

	"github.com/DerKnerd/gorp/v4"
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

		dbMap.AddTableWithName[AnalyticsEntry]("analytics")
		dbMap.AddTableWithName[ApiKey]("api_key")
		dbMap.AddTableWithName[ApiToken]("api_token")
		dbMap.AddTableWithName[Artist]("users")
		dbMap.AddTableWithName[BlogCategory]("blog_category")
		dbMap.AddTableWithName[BlogPost]("blog_post")
		dbMap.AddTableWithName[BlogPostSection]("blog_post_section")
		dbMap.AddTableWithName[ClassicPage]("classic_page")
		dbMap.AddTableWithName[File]("file").
			SetUniqueTogether("folder_id", "name")
		dbMap.AddTableWithName[FileTag]("file_tag")
		dbMap.AddTableWithName[Folder]("folder").
			SetUniqueTogether("parent_id", "name")
		dbMap.AddTableWithName[Gallery]("gallery")
		dbMap.AddTableWithName[GalleryFilePosition]("gallery_file_position")
		dbMap.AddTableWithName[KnownDevice]("known_device")
		dbMap.AddTableWithName[Menu]("menu")
		dbMap.AddTableWithName[MenuItem]("menu_item")
		dbMap.AddTableWithName[ModernPage]("modern_page")
		dbMap.AddTableWithName[ModernPageSection]("modern_page_section")
		dbMap.AddTableWithName[Theme]("theme")
		dbMap.AddTableWithName[ThemeAsset]("theme_asset").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[ThemeBlogCategory]("theme_blog_category").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[ThemeClassicPage]("theme_classic_page").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[ThemeFile]("theme_file").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[ThemeGallery]("theme_gallery").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[ThemeMenu]("theme_menu").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[ThemeModernPage]("theme_modern_page").
			SetUniqueTogether("theme_id", "name")
		dbMap.AddTableWithName[UploadingFile]("uploading_file")
		dbMap.AddTableWithName[UploadingFileChunk]("uploading_file_chunk")
		dbMap.AddTableWithName[JinyaConfiguration]("jinya_configuration").
			SetUniqueTogether("key", "group")

		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		migrations.BaseMigration(conn)
		migrations.CreateAdditionalTables(conn)
		migrations.CreateForeignKeys(conn)
	}
}
