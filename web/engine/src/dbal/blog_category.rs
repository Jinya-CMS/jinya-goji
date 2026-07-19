use crate::{find_by_parent, find_by_theme, find_entity};

find_entity!(blog_category, BlogCategory);

find_by_parent!(blog_post, BlogPost, CategoryId, find_posts);

find_by_theme!(blog_category, BlogCategory, theme_blog_category, find_by_theme);
