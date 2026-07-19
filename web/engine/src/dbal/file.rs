use crate::{find_by_theme, find_entity};

find_entity!(file, File);

find_by_theme!(file, File, theme_file, find_by_theme);
