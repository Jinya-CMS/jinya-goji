# <a id="goji_theme"></a>World goji-theme



 - Imports:
    - interface `jinya:goji/types`
    - interface `jinya:goji/db`
    - interface `jinya:goji/http`
    - interface `jinya:goji/theme`
    - type `blog-category`
    - type `blog-post`
    - type `gallery`
    - type `modern-page`
    - type `classic-page`
    - type `user`
    - type `theme-response`
 - Exports:
    - function `handle-blog-category`
    - function `handle-blog-post`
    - function `handle-blog-homepage`
    - function `handle-homepage`
    - function `handle-gallery`
    - function `handle-modern-page`
    - function `handle-classic-page`
    - function `handle-profile`
    - function `handle-error`

## <a id="jinya_goji_types"></a>Import interface jinya:goji/types



----

### Types

#### <a id="blog_category"></a>`record blog-category`



##### Record Fields

- <a id="blog_category.name"></a>`name`: option<`string`>
- <a id="blog_category.description"></a>`description`: option<`string`>

#### <a id="gallery_type"></a>`enum gallery-type`



##### Enum Cases

- <a id="gallery_type.masonry"></a>`masonry`
- <a id="gallery_type.sequence"></a>`sequence`

#### <a id="gallery_orientation"></a>`enum gallery-orientation`



##### Enum Cases

- <a id="gallery_orientation.horizontal"></a>`horizontal`
- <a id="gallery_orientation.vertical"></a>`vertical`

#### <a id="menu_item"></a>`resource menu-item`


#### <a id="menu"></a>`record menu`



##### Record Fields

- <a id="menu.id"></a>`id`: `s32`
- <a id="menu.name"></a>`name`: option<`string`>
- <a id="menu.logo"></a>`logo`: option<`s32`>
- <a id="menu.items"></a>`items`: list<own<[`menu-item`](#menu_item)>>

#### <a id="user"></a>`record user`



##### Record Fields

- <a id="user.email"></a>`email`: option<`string`>
- <a id="user.user_name"></a>`user-name`: option<`string`>
- <a id="user.profile_picture"></a>`profile-picture`: option<`string`>
- <a id="user.about_me"></a>`about-me`: option<`string`>

#### <a id="classic_page"></a>`record classic-page`



##### Record Fields

- <a id="classic_page.creator"></a>`creator`: option<[`user`](#user)>
- <a id="classic_page.updated_by"></a>`updated-by`: option<[`user`](#user)>
- <a id="classic_page.created_at"></a>`created-at`: option<`string`>
- <a id="classic_page.last_updated_at"></a>`last-updated-at`: option<`string`>
- <a id="classic_page.content"></a>`content`: option<`string`>
- <a id="classic_page.title"></a>`title`: option<`string`>

#### <a id="file"></a>`record file`



##### Record Fields

- <a id="file.id"></a>`id`: `s32`
- <a id="file.creator"></a>`creator`: option<[`user`](#user)>
- <a id="file.updated_by"></a>`updated-by`: option<[`user`](#user)>
- <a id="file.created_at"></a>`created-at`: option<`string`>
- <a id="file.last_updated_at"></a>`last-updated-at`: option<`string`>
- <a id="file.path"></a>`path`: option<`string`>
- <a id="file.name"></a>`name`: option<`string`>

#### <a id="gallery_file_position"></a>`record gallery-file-position`



##### Record Fields

- <a id="gallery_file_position.position"></a>`position`: option<`s32`>
- <a id="gallery_file_position.file"></a>`file`: option<[`file`](#file)>

#### <a id="gallery"></a>`record gallery`



##### Record Fields

- <a id="gallery.creator"></a>`creator`: option<[`user`](#user)>
- <a id="gallery.updated_by"></a>`updated-by`: option<[`user`](#user)>
- <a id="gallery.created_at"></a>`created-at`: option<`string`>
- <a id="gallery.last_updated_at"></a>`last-updated-at`: option<`string`>
- <a id="gallery.name"></a>`name`: option<`string`>
- <a id="gallery.description"></a>`description`: option<`string`>
- <a id="gallery.type"></a>`type`: option<[`gallery-type`](#gallery_type)>
- <a id="gallery.orientation"></a>`orientation`: option<[`gallery-orientation`](#gallery_orientation)>
- <a id="gallery.files"></a>`files`: list<[`gallery-file-position`](#gallery_file_position)>

#### <a id="blog_post_section"></a>`record blog-post-section`



##### Record Fields

- <a id="blog_post_section.gallery"></a>`gallery`: option<[`gallery`](#gallery)>
- <a id="blog_post_section.file"></a>`file`: option<[`file`](#file)>
- <a id="blog_post_section.html"></a>`html`: option<`string`>
- <a id="blog_post_section.link"></a>`link`: option<`string`>
- <a id="blog_post_section.position"></a>`position`: option<`s32`>

#### <a id="blog_post"></a>`record blog-post`



##### Record Fields

- <a id="blog_post.title"></a>`title`: option<`string`>
- <a id="blog_post.slug"></a>`slug`: option<`string`>
- <a id="blog_post.header_image"></a>`header-image`: option<[`file`](#file)>
- <a id="blog_post.public"></a>`public`: option<`bool`>
- <a id="blog_post.created_at"></a>`created-at`: option<`string`>
- <a id="blog_post.creator"></a>`creator`: option<[`user`](#user)>
- <a id="blog_post.updated_by"></a>`updated-by`: option<[`user`](#user)>
- <a id="blog_post.last_updated_at"></a>`last-updated-at`: option<`string`>
- <a id="blog_post.blog_category"></a>`blog-category`: option<[`blog-category`](#blog_category)>
- <a id="blog_post.sections"></a>`sections`: list<[`blog-post-section`](#blog_post_section)>

#### <a id="modern_page_section"></a>`record modern-page-section`



##### Record Fields

- <a id="modern_page_section.position"></a>`position`: option<`s32`>
- <a id="modern_page_section.gallery"></a>`gallery`: option<[`gallery`](#gallery)>
- <a id="modern_page_section.file"></a>`file`: option<[`file`](#file)>
- <a id="modern_page_section.html"></a>`html`: option<`string`>
- <a id="modern_page_section.action"></a>`action`: option<`string`>
- <a id="modern_page_section.script"></a>`script`: option<`string`>
- <a id="modern_page_section.target"></a>`target`: option<`string`>

#### <a id="modern_page"></a>`record modern-page`



##### Record Fields

- <a id="modern_page.creator"></a>`creator`: option<[`user`](#user)>
- <a id="modern_page.updated_by"></a>`updated-by`: option<[`user`](#user)>
- <a id="modern_page.created_at"></a>`created-at`: option<`string`>
- <a id="modern_page.last_updated_at"></a>`last-updated-at`: option<`string`>
- <a id="modern_page.name"></a>`name`: option<`string`>
- <a id="modern_page.sections"></a>`sections`: list<[`modern-page-section`](#modern_page_section)>

#### <a id="db_value"></a>`variant db-value`



##### Variant Cases

- <a id="db_value.null_value"></a>`null-value`
- <a id="db_value.bool_value"></a>`bool-value`: `bool`
- <a id="db_value.s16_value"></a>`s16-value`: `s16`
- <a id="db_value.s32_value"></a>`s32-value`: `s32`
- <a id="db_value.s64_value"></a>`s64-value`: `s64`
- <a id="db_value.f32_value"></a>`f32-value`: `f32`
- <a id="db_value.f64_value"></a>`f64-value`: `f64`
- <a id="db_value.string_value"></a>`string-value`: `string`
- <a id="db_value.bytes_value"></a>`bytes-value`: list<`u8`>

#### <a id="db_column"></a>`record db-column`



##### Record Fields

- <a id="db_column.name"></a>`name`: `string`
- <a id="db_column.value"></a>`value`: [`db-value`](#db_value)

#### <a id="theme_response"></a>`record theme-response`



##### Record Fields

- <a id="theme_response.content"></a>`content`: `string`
- <a id="theme_response.status_code"></a>`status-code`: `s16`

----

### Functions

#### <a id="method_menu_item_position"></a>`[method]menu-item.position: func`



##### Params

- <a id="method_menu_item_position.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_position.0"></a> option<`s32`>

#### <a id="method_menu_item_title"></a>`[method]menu-item.title: func`



##### Params

- <a id="method_menu_item_title.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_title.0"></a> option<`string`>

#### <a id="method_menu_item_route"></a>`[method]menu-item.route: func`



##### Params

- <a id="method_menu_item_route.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_route.0"></a> option<`string`>

#### <a id="method_menu_item_highlighted"></a>`[method]menu-item.highlighted: func`



##### Params

- <a id="method_menu_item_highlighted.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_highlighted.0"></a> option<`bool`>

#### <a id="method_menu_item_artist"></a>`[method]menu-item.artist: func`



##### Params

- <a id="method_menu_item_artist.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_artist.0"></a> option<[`user`](#user)>

#### <a id="method_menu_item_gallery"></a>`[method]menu-item.gallery: func`



##### Params

- <a id="method_menu_item_gallery.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_gallery.0"></a> option<[`gallery`](#gallery)>

#### <a id="method_menu_item_modern_page"></a>`[method]menu-item.modern-page: func`



##### Params

- <a id="method_menu_item_modern_page.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_modern_page.0"></a> option<[`modern-page`](#modern_page)>

#### <a id="method_menu_item_classic_page"></a>`[method]menu-item.classic-page: func`



##### Params

- <a id="method_menu_item_classic_page.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_classic_page.0"></a> option<[`classic-page`](#classic_page)>

#### <a id="method_menu_item_blog_category"></a>`[method]menu-item.blog-category: func`



##### Params

- <a id="method_menu_item_blog_category.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_blog_category.0"></a> option<[`blog-category`](#blog_category)>

#### <a id="method_menu_item_blog_home_page"></a>`[method]menu-item.blog-home-page: func`



##### Params

- <a id="method_menu_item_blog_home_page.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_blog_home_page.0"></a> option<`bool`>

#### <a id="method_menu_item_children"></a>`[method]menu-item.children: func`



##### Params

- <a id="method_menu_item_children.self"></a>`self`: borrow<[`menu-item`](#menu_item)>

##### Return values

- <a id="method_menu_item_children.0"></a> list<own<[`menu-item`](#menu_item)>>

## <a id="jinya_goji_db"></a>Import interface jinya:goji/db



----

### Types

#### <a id="blog_category"></a>`type blog-category`
[`blog-category`](#blog_category)
<p>

#### <a id="blog_post"></a>`type blog-post`
[`blog-post`](#blog_post)
<p>

#### <a id="blog_post_section"></a>`type blog-post-section`
[`blog-post-section`](#blog_post_section)
<p>

#### <a id="classic_page"></a>`type classic-page`
[`classic-page`](#classic_page)
<p>

#### <a id="db_column"></a>`type db-column`
[`db-column`](#db_column)
<p>

#### <a id="file"></a>`type file`
[`file`](#file)
<p>

#### <a id="gallery"></a>`type gallery`
[`gallery`](#gallery)
<p>

#### <a id="gallery_file_position"></a>`type gallery-file-position`
[`gallery-file-position`](#gallery_file_position)
<p>

#### <a id="menu"></a>`type menu`
[`menu`](#menu)
<p>

#### <a id="menu_item"></a>`type menu-item`
[`menu-item`](#menu_item)
<p>

#### <a id="modern_page"></a>`type modern-page`
[`modern-page`](#modern_page)
<p>

#### <a id="modern_page_section"></a>`type modern-page-section`
[`modern-page-section`](#modern_page_section)
<p>

#### <a id="user"></a>`type user`
[`user`](#user)
<p>

----

### Functions

#### <a id="find_blog_category_by_id"></a>`find-blog-category-by-id: func`



##### Params

- <a id="find_blog_category_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_blog_category_by_id.0"></a> option<[`blog-category`](#blog_category)>

#### <a id="find_blog_post_by_id"></a>`find-blog-post-by-id: func`



##### Params

- <a id="find_blog_post_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_blog_post_by_id.0"></a> option<[`blog-post`](#blog_post)>

#### <a id="find_classic_page_by_id"></a>`find-classic-page-by-id: func`



##### Params

- <a id="find_classic_page_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_classic_page_by_id.0"></a> option<[`classic-page`](#classic_page)>

#### <a id="find_file_by_id"></a>`find-file-by-id: func`



##### Params

- <a id="find_file_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_file_by_id.0"></a> option<[`file`](#file)>

#### <a id="find_gallery_by_id"></a>`find-gallery-by-id: func`



##### Params

- <a id="find_gallery_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_gallery_by_id.0"></a> option<[`gallery`](#gallery)>

#### <a id="find_menu_by_id"></a>`find-menu-by-id: func`



##### Params

- <a id="find_menu_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_menu_by_id.0"></a> option<[`menu`](#menu)>

#### <a id="find_modern_page_by_id"></a>`find-modern-page-by-id: func`



##### Params

- <a id="find_modern_page_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_modern_page_by_id.0"></a> option<[`modern-page`](#modern_page)>

#### <a id="find_user_by_id"></a>`find-user-by-id: func`



##### Params

- <a id="find_user_by_id.id"></a>`id`: `s32`

##### Return values

- <a id="find_user_by_id.0"></a> option<[`user`](#user)>

#### <a id="find_all_blog_categories"></a>`find-all-blog-categories: func`



##### Return values

- <a id="find_all_blog_categories.0"></a> list<[`blog-category`](#blog_category)>

#### <a id="find_all_blog_posts"></a>`find-all-blog-posts: func`



##### Return values

- <a id="find_all_blog_posts.0"></a> list<[`blog-post`](#blog_post)>

#### <a id="find_all_classic_pages"></a>`find-all-classic-pages: func`



##### Return values

- <a id="find_all_classic_pages.0"></a> list<[`classic-page`](#classic_page)>

#### <a id="find_all_files"></a>`find-all-files: func`



##### Return values

- <a id="find_all_files.0"></a> list<[`file`](#file)>

#### <a id="find_all_galleries"></a>`find-all-galleries: func`



##### Return values

- <a id="find_all_galleries.0"></a> list<[`gallery`](#gallery)>

#### <a id="find_all_menus"></a>`find-all-menus: func`



##### Return values

- <a id="find_all_menus.0"></a> list<[`menu`](#menu)>

#### <a id="find_all_modern_pages"></a>`find-all-modern-pages: func`



##### Return values

- <a id="find_all_modern_pages.0"></a> list<[`modern-page`](#modern_page)>

#### <a id="find_all_users"></a>`find-all-users: func`



##### Return values

- <a id="find_all_users.0"></a> list<[`user`](#user)>

#### <a id="query_db"></a>`query-db: func`



##### Params

- <a id="query_db.query"></a>`query`: `string`

##### Return values

- <a id="query_db.0"></a> list<list<[`db-column`](#db_column)>>

## <a id="jinya_goji_http"></a>Import interface jinya:goji/http



----

### Functions

#### <a id="get"></a>`get: func`



##### Params

- <a id="get.url"></a>`url`: `string`

##### Return values

- <a id="get.0"></a> `string`

## <a id="jinya_goji_theme"></a>Import interface jinya:goji/theme



----

### Types

#### <a id="blog_category"></a>`type blog-category`
[`blog-category`](#blog_category)
<p>

#### <a id="classic_page"></a>`type classic-page`
[`classic-page`](#classic_page)
<p>

#### <a id="file"></a>`type file`
[`file`](#file)
<p>

#### <a id="gallery"></a>`type gallery`
[`gallery`](#gallery)
<p>

#### <a id="menu"></a>`type menu`
[`menu`](#menu)
<p>

#### <a id="menu_item"></a>`type menu-item`
[`menu-item`](#menu_item)
<p>

#### <a id="modern_page"></a>`type modern-page`
[`modern-page`](#modern_page)
<p>

----

### Functions

#### <a id="get_theme_config_value"></a>`get-theme-config-value: func`



##### Params

- <a id="get_theme_config_value.group"></a>`group`: `string`
- <a id="get_theme_config_value.key"></a>`key`: `string`

##### Return values

- <a id="get_theme_config_value.0"></a> option<`string`>

#### <a id="get_active_menu_item"></a>`get-active-menu-item: func`



##### Return values

- <a id="get_active_menu_item.0"></a> option<own<[`menu-item`](#menu_item)>>

#### <a id="is_menu_item_active"></a>`is-menu-item-active: func`



##### Params

- <a id="is_menu_item_active.item"></a>`item`: own<[`menu-item`](#menu_item)>

##### Return values

- <a id="is_menu_item_active.0"></a> `bool`

#### <a id="is_parent_of_active_menu_item"></a>`is-parent-of-active-menu-item: func`



##### Params

- <a id="is_parent_of_active_menu_item.item"></a>`item`: own<[`menu-item`](#menu_item)>

##### Return values

- <a id="is_parent_of_active_menu_item.0"></a> `bool`

#### <a id="find_modern_page_by_theme"></a>`find-modern-page-by-theme: func`



##### Params

- <a id="find_modern_page_by_theme.name"></a>`name`: `string`

##### Return values

- <a id="find_modern_page_by_theme.0"></a> option<[`modern-page`](#modern_page)>

#### <a id="find_classic_page_by_theme"></a>`find-classic-page-by-theme: func`



##### Params

- <a id="find_classic_page_by_theme.name"></a>`name`: `string`

##### Return values

- <a id="find_classic_page_by_theme.0"></a> option<[`classic-page`](#classic_page)>

#### <a id="find_file_by_theme"></a>`find-file-by-theme: func`



##### Params

- <a id="find_file_by_theme.name"></a>`name`: `string`

##### Return values

- <a id="find_file_by_theme.0"></a> option<[`file`](#file)>

#### <a id="find_menu_by_theme"></a>`find-menu-by-theme: func`



##### Params

- <a id="find_menu_by_theme.name"></a>`name`: `string`

##### Return values

- <a id="find_menu_by_theme.0"></a> option<[`menu`](#menu)>

#### <a id="find_gallery_by_theme"></a>`find-gallery-by-theme: func`



##### Params

- <a id="find_gallery_by_theme.name"></a>`name`: `string`

##### Return values

- <a id="find_gallery_by_theme.0"></a> option<[`gallery`](#gallery)>

#### <a id="find_blog_category_by_theme"></a>`find-blog-category-by-theme: func`



##### Params

- <a id="find_blog_category_by_theme.name"></a>`name`: `string`

##### Return values

- <a id="find_blog_category_by_theme.0"></a> option<[`blog-category`](#blog_category)>

## Exported types from world `goji-theme`

----

### Types

#### <a id="blog_category"></a>`type blog-category`
[`blog-category`](#blog_category)
<p>

#### <a id="blog_post"></a>`type blog-post`
[`blog-post`](#blog_post)
<p>

#### <a id="gallery"></a>`type gallery`
[`gallery`](#gallery)
<p>

#### <a id="modern_page"></a>`type modern-page`
[`modern-page`](#modern_page)
<p>

#### <a id="classic_page"></a>`type classic-page`
[`classic-page`](#classic_page)
<p>

#### <a id="user"></a>`type user`
[`user`](#user)
<p>

#### <a id="theme_response"></a>`type theme-response`
[`theme-response`](#theme_response)
<p>

## Exported functions from world `goji-theme`

#### <a id="handle_blog_category"></a>`handle-blog-category: func`



##### Params

- <a id="handle_blog_category.entity"></a>`entity`: [`blog-category`](#blog_category)

##### Return values

- <a id="handle_blog_category.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_blog_post"></a>`handle-blog-post: func`



##### Params

- <a id="handle_blog_post.entity"></a>`entity`: [`blog-post`](#blog_post)

##### Return values

- <a id="handle_blog_post.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_blog_homepage"></a>`handle-blog-homepage: func`



##### Return values

- <a id="handle_blog_homepage.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_homepage"></a>`handle-homepage: func`



##### Return values

- <a id="handle_homepage.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_gallery"></a>`handle-gallery: func`



##### Params

- <a id="handle_gallery.entity"></a>`entity`: [`gallery`](#gallery)

##### Return values

- <a id="handle_gallery.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_modern_page"></a>`handle-modern-page: func`



##### Params

- <a id="handle_modern_page.entity"></a>`entity`: [`modern-page`](#modern_page)

##### Return values

- <a id="handle_modern_page.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_classic_page"></a>`handle-classic-page: func`



##### Params

- <a id="handle_classic_page.entity"></a>`entity`: [`classic-page`](#classic_page)

##### Return values

- <a id="handle_classic_page.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_profile"></a>`handle-profile: func`



##### Params

- <a id="handle_profile.entity"></a>`entity`: [`user`](#user)

##### Return values

- <a id="handle_profile.0"></a> [`theme-response`](#theme_response)

#### <a id="handle_error"></a>`handle-error: func`



##### Params

- <a id="handle_error.status_code"></a>`status-code`: `s32`

##### Return values

- <a id="handle_error.0"></a> [`theme-response`](#theme_response)

