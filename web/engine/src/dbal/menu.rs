use crate::{find_by_parent, find_by_theme, find_entity};

find_entity!(menu, Menu);

find_by_parent!(menu_item, MenuItem, MenuId, find_menu_items);

find_by_theme!(menu, Menu, theme_menu, ThemeMenu, find_by_theme);
