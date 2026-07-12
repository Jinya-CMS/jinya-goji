package migrations

import (
	"database/sql"
	"log/slog"
)

func CreateForeignKeys(conn *sql.DB) {
	slog.Info("Add foreign keys for api key")
	_, err := conn.Exec(`
select add_foreign_key_if_not_exists('api_key', 'user_id', 'users', 'id');
`)
	if err != nil {
		panic(err)
	}

	slog.Info("Add foreign keys for blog")
	_, err = conn.Exec(`
select add_foreign_key_if_not_exists('blog_category', 'parent_id', 'blog_category', 'id');

select add_foreign_key_if_not_exists('blog_post', 'category_id', 'blog_category', 'id');
select add_foreign_key_if_not_exists('blog_post', 'header_image_id', 'file', 'id');
select add_foreign_key_if_not_exists('blog_post', 'creator_id', 'users', 'id');
select add_foreign_key_if_not_exists('blog_post', 'updated_by_id', 'users', 'id');

select add_foreign_key_if_not_exists('blog_post_section', 'blog_post_id', 'blog_post', 'id');
select add_foreign_key_if_not_exists('blog_post_section', 'file_id', 'file', 'id');
select add_foreign_key_if_not_exists('blog_post_section', 'gallery_id', 'gallery', 'id');
`)
	if err != nil {
		panic(err)
	}

	slog.Info("Add foreign keys for media")
	_, err = conn.Exec(`
select add_foreign_key_if_not_exists('file', 'folder_id', 'folder', 'id');
select add_foreign_key_if_not_exists('file', 'creator_id', 'users', 'id');
select add_foreign_key_if_not_exists('file', 'updated_by_id', 'users', 'id');
select add_foreign_key_if_not_exists('folder', 'parent_id', 'folder', 'id');

select add_foreign_key_if_not_exists('gallery', 'creator_id', 'users', 'id');
select add_foreign_key_if_not_exists('gallery', 'updated_by_id', 'users', 'id');

select add_foreign_key_if_not_exists('gallery_file_position', 'file_id', 'file', 'id');
select add_foreign_key_if_not_exists('gallery_file_position', 'gallery_id', 'gallery', 'id');

select add_foreign_key_if_not_exists('uploading_file', 'file_id', 'file', 'id');

select add_foreign_key_if_not_exists('uploading_file_chunk', 'uploading_file_id', 'uploading_file', 'id');
`)
	if err != nil {
		panic(err)
	}

	slog.Info("Add foreign keys for known devices")
	_, err = conn.Exec(`
select add_foreign_key_if_not_exists('known_device', 'user_id', 'users', 'id');
`)
	if err != nil {
		panic(err)
	}

	slog.Info("Add foreign keys for menu")
	_, err = conn.Exec(`
select add_foreign_key_if_not_exists('menu', 'logo', 'file', 'id');

select add_foreign_key_if_not_exists('menu_item', 'artist_id', 'users', 'id');
select add_foreign_key_if_not_exists('menu_item', 'category_id', 'blog_category', 'id');
select add_foreign_key_if_not_exists('menu_item', 'gallery_id', 'gallery', 'id');
select add_foreign_key_if_not_exists('menu_item', 'parent_id', 'menu_item', 'id');
select add_foreign_key_if_not_exists('menu_item', 'menu_id', 'menu', 'id');
select add_foreign_key_if_not_exists('menu_item', 'page_id', 'classic_page', 'id');
select add_foreign_key_if_not_exists('menu_item', 'segment_page_id', 'modern_page', 'id');
`)
	if err != nil {
		panic(err)
	}

	slog.Info("Add foreign keys for pages")
	_, err = conn.Exec(`
select add_foreign_key_if_not_exists('classic_page', 'creator_id', 'users', 'id');
select add_foreign_key_if_not_exists('classic_page', 'updated_by_id', 'users', 'id');

select add_foreign_key_if_not_exists('modern_page_section', 'file_id', 'file', 'id');
select add_foreign_key_if_not_exists('modern_page_section', 'gallery_id', 'gallery', 'id');
select add_foreign_key_if_not_exists('modern_page_section', 'page_id', 'modern_page', 'id');

select add_foreign_key_if_not_exists('modern_page', 'creator_id', 'users', 'id');
select add_foreign_key_if_not_exists('modern_page', 'updated_by_id', 'users', 'id');
`)
	if err != nil {
		panic(err)
	}

	slog.Info("Add foreign keys for pages")
	_, err = conn.Exec(`
select add_foreign_key_if_not_exists('theme_asset', 'theme_id', 'theme', 'id');

select add_foreign_key_if_not_exists('theme_blog_category', 'blog_category_id', 'blog_category', 'id');
select add_foreign_key_if_not_exists('theme_blog_category', 'theme_id', 'theme', 'id');

select add_foreign_key_if_not_exists('theme_file', 'file_id', 'file', 'id');
select add_foreign_key_if_not_exists('theme_file', 'theme_id', 'theme', 'id');

select add_foreign_key_if_not_exists('theme_gallery', 'gallery_id', 'gallery', 'id');
select add_foreign_key_if_not_exists('theme_gallery', 'theme_id', 'theme', 'id');

select add_foreign_key_if_not_exists('theme_menu', 'menu_id', 'menu', 'id');
select add_foreign_key_if_not_exists('theme_menu', 'theme_id', 'theme', 'id');

select add_foreign_key_if_not_exists('theme_page', 'page_id', 'classic_page', 'id');
select add_foreign_key_if_not_exists('theme_page', 'theme_id', 'theme', 'id');

select add_foreign_key_if_not_exists('theme_modern_page', 'modern_page_id', 'modern_page', 'id');
select add_foreign_key_if_not_exists('theme_modern_page', 'theme_id', 'theme', 'id');
`)
	if err != nil {
		panic(err)
	}
}
