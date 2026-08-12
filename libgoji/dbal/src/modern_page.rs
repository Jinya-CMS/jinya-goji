use crate::{find_by_parent, find_by_theme, find_entity};

find_entity!(modern_page, ModernPage);

find_by_parent!(
    modern_page_section,
    ModernPageSection,
    PageId,
    find_sections
);

find_by_theme!(
    modern_page,
    ModernPage,
    theme_modern_page,
    ThemeModernPage,
    find_by_theme
);
