use crate::{find_by_theme, find_entity};

find_entity!(classic_page, ClassicPage);

find_by_theme!(classic_page, ClassicPage, theme_classic_page, find_by_theme);
