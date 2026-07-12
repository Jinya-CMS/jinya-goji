package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/DerKnerd/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var dbMap *gorp.DbMap

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func SetupDatabase() {
	if dbMap == nil {
		db, err := sql.Open("mysql", os.Getenv("DATABASE_URL")+"?parseTime=true")
		if err != nil {
			panic(err)
		}

		db.SetConnMaxLifetime(3 * time.Minute)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		dialect := gorp.MySQLDialect{
			Engine:   "innodb",
			Encoding: "utf8mb4",
		}

		dbMap = &gorp.DbMap{Db: db, Dialect: dialect}

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
		AddTableWithName[JinyaConfiguration]("jinya_configuration")
		AddTableWithName[IpAddress]("ip_address")

		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		_, err = dbMap.Exec(`
create table if not exists configuration
(
    id                        int auto_increment primary key,
    current_frontend_theme_id int null,
    constraint UNIQ_configuration_current_frontend_theme_id
        unique (current_frontend_theme_id),
    constraint FK_theme_configuration
        foreign key (current_frontend_theme_id) references jinya.theme (id)
)`)
		if err != nil {
			panic(err)
		}

		_, err = dbMap.Exec(`
create table if not exists file_tag_file
(
    file_tag_id int not null,
    file_id     int not null,
    primary key (file_id, file_tag_id),
    constraint FK_file_tag_file_file_id
        foreign key (file_id) references jinya.file (id)
            on delete cascade,
    constraint FK_file_tag_file_file_tag_id
        foreign key (file_tag_id) references jinya.file_tag (id)
            on delete cascade
)`)
		if err != nil {
			panic(err)
		}

		if err := addForeignKey("api_key", "FK_users_api_key", "user_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_category", "FK_category_category_parent_id", "parent_id", "blog_category", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post", "FK_blog_post_category_category_id", "category_id", "blog_category", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post", "FK_blog_post_file_header_image_id", "header_image_id", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post", "FK_blog_post_users_creator_id", "creator_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post", "FK_users_blog_post_updated_by_id", "updated_by_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post_segment", "FK_blog_post_segment_blog_post_blog_post_id", "blog_post_id", "blog_post", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post_segment", "FK_blog_post_segment_file_file_id", "file_id", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("blog_post_segment", "FK_blog_post_segment_gallery_gallery_id", "gallery_id", "gallery", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("file", "FK_file_folder", "folder_id", "folder", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("file", "FK_users_file_creator_id", "creator_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("file", "FK_users_file_updated_by_id", "updated_by_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("folder", "FK_folder_parent_folder", "parent_id", "folder", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("form", "FK_users_form_creator_id", "creator_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("form", "FK_users_form_last_updated_by_id", "updated_by_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("form_item", "FK_form_form_item", "form_id", "form", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("gallery", "FK_users_gallery_creator_id", "creator_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("gallery", "FK_users_gallery_last_updated_by_id", "updated_by_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("gallery_file_position", "FK_file_gallery_file_position", "file_id", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("gallery_file_position", "FK_gallery_gallery_file_position", "gallery_id", "gallery", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("known_device", "FK_users_known_device", "user_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu", "FK_file_menu", "logo", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_artist_menu_item", "artist_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_category_menu_item", "category_id", "blog_category", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_form_menu_item", "form_id", "form", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_gallery_menu_item", "gallery_id", "gallery", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_menu_item_menu_item_parent", "parent_id", "menu_item", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_menu_menu_item", "menu_id", "menu", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_page_menu_item", "page_id", "page", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("menu_item", "FK_segment_page_menu_item", "segment_page_id", "segment_page", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("page", "FK_users_page_creator_id", "creator_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("page", "FK_users_page_last_updated_by_id", "updated_by_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("segment", "FK_file_segment", "file_id", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("segment", "FK_gallery_segment", "gallery_id", "gallery", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("segment", "FK_segment_page_segment", "page_id", "segment_page", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("segment_page", "FK_users_segment_page_creator_id", "creator_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("segment_page", "FK_users_segment_page_last_updated_by_id", "updated_by_id", "users", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_asset", "FK_theme_asset_theme", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_blog_category", "FK_blog_category_theme_blog_category", "blog_category_id", "blog_category", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_blog_category", "FK_theme_theme_blog_category", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_file", "FK_file_theme_file", "file_id", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_file", "FK_theme_theme_file", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_form", "FK_form_theme_form", "form_id", "form", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_form", "FK_theme_theme_form", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_gallery", "FK_gallery_theme_gallery", "gallery_id", "gallery", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_gallery", "FK_theme_theme_gallery", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_menu", "FK_menu_theme_menu", "menu_id", "menu", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_menu", "FK_theme_theme_menu", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_page", "FK_page_theme_page", "page_id", "page", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_page", "FK_theme_theme_page", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_segment_page", "FK_segment_page_theme_segment_page", "segment_page_id", "segment_page", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("theme_segment_page", "FK_theme_theme_segment_page", "theme_id", "theme", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("uploading_file", "FK_uploading_file_file", "file_id", "file", "id"); err != nil {
			panic(err)
		}
		if err := addForeignKey("uploading_file_chunk", "FK_uploading_file_uploading_file_chunk", "uploading_file_id", "uploading_file", "id"); err != nil {
			panic(err)
		}
	}
}

func addForeignKey(table, name, column, referencedTable, referencedColumn string) error {
	var count int64
	err := dbMap.Db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.TABLE_CONSTRAINTS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND CONSTRAINT_NAME = ?
		  AND CONSTRAINT_TYPE = 'FOREIGN KEY'
	`, table, name).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		query := fmt.Sprintf(`
ALTER TABLE %s
    ADD CONSTRAINT %s FOREIGN KEY IF NOT EXISTS (%s) REFERENCES %s (%s)`, table, name, column, referencedTable, referencedColumn)
		_, err = dbMap.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
