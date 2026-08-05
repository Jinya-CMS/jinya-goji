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
		AddTableWithName[ApiToken]("api_token")
		AddTableWithName[Artist]("users")
		AddTableWithName[BlogCategory]("blog_category")
		AddTableWithName[BlogPost]("blog_post")
		AddTableWithName[BlogPostSection]("blog_post_section")
		AddTableWithName[ClassicPage]("classic_page")
		AddTableWithName[File]("file").
			SetUniqueTogether("folder_id", "name")
		AddTableWithName[FileTag]("file_tag")
		AddTableWithName[Folder]("folder").
			SetUniqueTogether("parent_id", "name")
		AddTableWithName[Gallery]("gallery")
		AddTableWithName[GalleryFilePosition]("gallery_file_position")
		AddTableWithName[KnownDevice]("known_device")
		AddTableWithName[Menu]("menu")
		AddTableWithName[MenuItem]("menu_item")
		AddTableWithName[ModernPage]("modern_page")
		AddTableWithName[ModernPageSection]("modern_page_section")
		AddTableWithName[Theme]("theme")
		AddTableWithName[ThemeAsset]("theme_asset").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[ThemeBlogCategory]("theme_blog_category").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[ThemeClassicPage]("theme_classic_page").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[ThemeFile]("theme_file").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[ThemeGallery]("theme_gallery").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[ThemeMenu]("theme_menu").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[ThemeModernPage]("theme_modern_page").
			SetUniqueTogether("theme_id", "name")
		AddTableWithName[UploadingFile]("uploading_file")
		AddTableWithName[UploadingFileChunk]("uploading_file_chunk")
		AddTableWithName[JinyaConfiguration]("jinya_configuration").
			SetUniqueTogether("key", "group")

		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		migrations.BaseMigration(conn)
		migrations.CreateAdditionalTables(conn)
		migrations.CreateForeignKeys(conn)

		migrations.AddGalleryGroups(conn)
	}
}
