import { Alpine } from '../../../../lib/alpine.js';
import { getClassicPages } from '../../foundation/api/classic-pages.js';
import {
  activateTheme,
  compileThemeAssets,
  getTheme,
  getThemeBlogCategories,
  getThemeClassicPages,
  getThemeConfigurationStructure,
  getThemeDefaultConfiguration,
  getThemeFiles,
  getThemeGalleries,
  getThemeMenus,
  getThemeModernPages,
  getThemes,
  getThemeStyleVariables,
  updateTheme,
  updateThemeBlogCategory,
  updateThemeClassicPage,
  updateThemeConfiguration,
  updateThemeFile,
  updateThemeGallery,
  updateThemeMenu,
  updateThemeModernPage,
  updateThemeVariables,
  uploadTheme,
} from '../../foundation/api/themes.js';
import { getGalleries } from '../../foundation/api/galleries.js';
import { getModernPages } from '../../foundation/api/modern-pages.js';
import { getMenus } from '../../foundation/api/menus.js';
import { getBlogCategories } from '../../foundation/api/blog-categories.js';
import localize from '../../foundation/utils/localize.js';
import confirm from '../../foundation/ui/confirm.js';
import { getMediaDatabase } from '../../foundation/database/media.js';
import filePicker from '../../foundation/ui/filePicker.js';
import { dismissToast, negativeToast, positiveToast } from '../../foundation/ui/toast.js';

const mediaDatabase = getMediaDatabase();

Alpine.data('themesData', () => ({
  themes: [],
  linksLoaded: false,
  selectedTheme: null,
  themeSelected: false,
  activeTab: 'details',
  galleries: [],
  files: [],
  classicPages: [],
  modernPages: [],
  menus: [],
  blogCategories: [],
  forms: [],
  toastId: null,
  getValueForCurrentLanguage(data) {
    const currentLanguage = navigator.language.substring(0, 2);
    if (Object.prototype.hasOwnProperty.call(data, currentLanguage)) {
      return data[currentLanguage];
    }
    if (Object.prototype.hasOwnProperty.call(data, 'en')) {
      return data.en;
    }

    return data;
  },
  getVariableLabel(variable) {
    return variable
      .replace('$', '')
      .replaceAll('-', ' ')
      .replace(/(^\w)|(\s+\w)/g, (letter) => letter.toUpperCase());
  },
  get title() {
    return this.selectedThemeDisplayName;
  },
  get selectedThemeDisplayName() {
    return this.getValueForCurrentLanguage(this.selectedTheme.displayName);
  },
  get selectedThemePreviewImage() {
    return `/api/theme/${this.selectedTheme.id}/preview`;
  },
  get selectedThemeDescription() {
    return this.getValueForCurrentLanguage(this.selectedTheme.description);
  },
  async init() {
    Promise.all([getGalleries(), getClassicPages(), getModernPages(), getMenus(), getBlogCategories()]).then(
      ([galleries, classicPages, modernPages, menus, blogCategories]) => {
        this.galleries = galleries;
        this.classicPages = classicPages;
        this.modernPages = modernPages;
        this.menus = menus;
        this.blogCategories = blogCategories;
      },
    );

    mediaDatabase.cacheMedia().then(() => {});

    this.themes = await getThemes();
    if (this.themes.length > 0) {
      await this.selectTheme(this.themes[0]);
    }
    this.$watch('selectedTheme', () => dismissToast(this.toastId));
  },
  destroy() {
    dismissToast(this.toastId);
  },
  selectedThemeData: {
    modernPages: {},
    classicPages: {},
    menus: {},
    galleries: {},
    files: {},
    blogCategories: {},
    styleVariables: {},
    scssVariables: {},
    defaultConfiguration: {},
    configurationStructure: {},
    configuration: {},
  },
  openCreateDialog() {
    this.upload.error.reset();
    this.upload.file = null;
    this.upload.open = true;
  },
  openEditDialog() {
    this.edit.error.reset();
    this.edit.file = null;
    this.edit.open = true;
  },
  loadThemeConfiguration() {
    const configuration = Array.isArray(this.selectedTheme.configuration) ? {} : this.selectedTheme.configuration;

    for (const group of this.selectedThemeData.configurationStructure.groups) {
      configuration[group.name] = configuration[group.name] ?? {};
      for (const field of group.fields) {
        configuration[group.name][field.name] = configuration[group.name][field.name] ?? null;
      }
    }

    this.selectedThemeData.configuration = configuration;
  },
  loadThemeVariables() {
    this.selectedThemeData.scssVariables = Array.isArray(this.selectedTheme.scssVariables)
      ? {}
      : this.selectedTheme.scssVariables;
  },
  async getSelectedFileName(id) {
    if (!id) {
      return localize({ key: 'design.themes.links.no_file' });
    }

    return (await mediaDatabase.getFileById(parseInt(id))).name ?? localize({ key: 'design.themes.links.no_file' });
  },
  async getFilePath(id) {
    if (!id) {
      return null;
    }

    return (await mediaDatabase.getFileById(parseInt(id))).path;
  },
  async selectFile(fileLink) {
    const id = this.selectedThemeData.files[fileLink];
    const file = await filePicker({
      title: localize({ key: 'file_picker.title' }),
      selectedFileId: id,
      cancelLabel: localize({ key: 'file_picker.dismiss' }),
      pickLabel: localize({ key: 'file_picker.pick' }),
    });
    if (file) {
      this.selectedThemeData.files[fileLink] = file.id;
    }
  },
  async loadLinks() {
    const [modernPages, classicPages, menus, galleries, files, blogCategories] = await Promise.all([
      getThemeModernPages(this.selectedTheme.id),
      getThemeClassicPages(this.selectedTheme.id),
      getThemeMenus(this.selectedTheme.id),
      getThemeGalleries(this.selectedTheme.id),
      getThemeFiles(this.selectedTheme.id),
      getThemeBlogCategories(this.selectedTheme.id),
    ]);

    this.selectedThemeData.modernPages = {};
    for (const [key, item] of Object.entries(modernPages)) {
      this.selectedThemeData.modernPages[key] = item.id;
    }

    this.selectedThemeData.classicPages = {};
    for (const [key, item] of Object.entries(classicPages)) {
      this.selectedThemeData.classicPages[key] = item.id;
    }

    this.selectedThemeData.menus = {};
    for (const [key, item] of Object.entries(menus)) {
      this.selectedThemeData.menus[key] = item.id;
    }

    this.selectedThemeData.galleries = {};
    for (const [key, item] of Object.entries(galleries)) {
      this.selectedThemeData.galleries[key] = item.id;
    }

    this.selectedThemeData.files = {};
    for (const [key, item] of Object.entries(files)) {
      this.selectedThemeData.files[key] = item.id;
    }

    this.selectedThemeData.blogCategories = {};
    for (const [key, item] of Object.entries(blogCategories)) {
      this.selectedThemeData.blogCategories[key] = item.id;
    }
  },
  async loadConfigurationStructure() {
    const [styleVariables, defaultConfiguration, configurationStructure] = await Promise.all([
      getThemeStyleVariables(this.selectedTheme.id),
      getThemeDefaultConfiguration(this.selectedTheme.id),
      getThemeConfigurationStructure(this.selectedTheme.id),
    ]);

    this.selectedThemeData.styleVariables = styleVariables;
    this.selectedThemeData.defaultConfiguration = defaultConfiguration;
    this.selectedThemeData.configurationStructure = configurationStructure;
  },
  async selectTheme(theme) {
    this.activeTab = 'details';

    this.selectedTheme = await getTheme(theme.id);
    await Promise.all([this.loadConfigurationStructure(), this.loadLinks()]);

    this.loadThemeConfiguration();
    this.loadThemeVariables();

    this.themeSelected = true;
  },
  async saveConfiguration() {
    dismissToast(this.toastId);
    try {
      await updateThemeConfiguration(this.selectedTheme.id, this.selectedThemeData.configuration);
      this.toastId = positiveToast({
        title: localize({ key: 'design.themes.configuration.success.title' }),
        message: localize({ key: 'design.themes.configuration.success.message' }),
      });
    } catch (e) {
      this.toastId = negativeToast({
        title: localize({ key: 'design.themes.configuration.error.title' }),
        message: localize({ key: 'design.themes.configuration.error.message' }),
      });
    }
  },
  async saveVariables() {
    dismissToast(this.toastId);
    try {
      await updateThemeVariables(this.selectedTheme.id, this.selectedThemeData.scssVariables);
      this.toastId = positiveToast({
        title: localize({ key: 'design.themes.variables.success.title' }),
        message: localize({ key: 'design.themes.variables.success.message' }),
      });
    } catch (e) {
      this.toastId = negativeToast({
        title: localize({ key: 'design.themes.variables.error.title' }),
        message: localize({ key: 'design.themes.variables.error.message' }),
      });
    }
  },
  async saveLinks() {
    dismissToast(this.toastId);
    const { modernPages, classicPages, menus, galleries, files, blogCategories } = this.selectedThemeData;
    const { id } = this.selectedTheme;
    try {
      await Promise.all([
        ...Object.entries(modernPages).map(([name, link]) =>
          updateThemeModernPage(id, name, link === 'null' ? null : link),
        ),
        ...Object.entries(classicPages).map(([name, link]) =>
          updateThemeClassicPage(id, name, link === 'null' ? null : link),
        ),
        ...Object.entries(menus).map(([name, link]) => updateThemeMenu(id, name, link === 'null' ? null : link)),
        ...Object.entries(galleries).map(([name, link]) => updateThemeGallery(id, name, link === 'null' ? null : link)),
        ...Object.entries(files).map(([name, link]) => updateThemeFile(id, name, link === 'null' ? null : link)),
        ...Object.entries(blogCategories).map(([name, link]) =>
          updateThemeBlogCategory(id, name, link === 'null' ? null : link),
        ),
      ]);

      this.toastId = positiveToast({
        title: localize({ key: 'design.themes.links.success.title' }),
        message: localize({ key: 'design.themes.links.success.message' }),
      });
    } catch (e) {
      this.toastId = negativeToast({
        title: localize({ key: 'design.themes.links.error.title' }),
        message: localize({ key: 'design.themes.links.error.message' }),
      });
    }
  },
  async activate() {
    dismissToast(this.toastId);
    const confirmed = await confirm({
      title: localize({ key: 'design.themes.activate.title' }),
      message: localize({
        key: 'design.themes.activate.message',
        values: this.selectedTheme,
      }),
      approveLabel: localize({ key: 'design.themes.activate.approve' }),
      declineLabel: localize({ key: 'design.themes.activate.decline' }),
    });
    if (confirmed) {
      try {
        await activateTheme(this.selectedTheme.id);
        this.toastId = positiveToast({
          title: localize({ key: 'design.themes.activate.success.title' }),
          message: localize({ key: 'design.themes.activate.success.message', values: this.selectedTheme }),
        });
      } catch (e) {
        this.toastId = negativeToast({
          title: localize({ key: 'design.themes.activate.error.title' }),
          message: localize({ key: 'design.themes.activate.error.message', values: this.selectedTheme }),
        });
      }
    }
  },
  async compileAssets() {
    dismissToast(this.toastId);
    const confirmed = await confirm({
      title: localize({ key: 'design.themes.assets.title' }),
      message: localize({
        key: 'design.themes.assets.message',
        values: this.selectedTheme,
      }),
      approveLabel: localize({ key: 'design.themes.assets.approve' }),
      declineLabel: localize({ key: 'design.themes.assets.decline' }),
    });
    if (confirmed) {
      try {
        await compileThemeAssets(this.selectedTheme.id);
        this.toastId = positiveToast({
          title: localize({ key: 'design.themes.assets.success.title' }),
          message: localize({ key: 'design.themes.assets.success.message', values: this.selectedTheme }),
        });
      } catch (e) {
        this.toastId = negativeToast({
          title: localize({ key: 'design.themes.assets.error.title' }),
          message: localize({ key: 'design.themes.assets.error.message', values: this.selectedTheme }),
        });
      }
    }
  },
  async uploadTheme() {
    try {
      const theme = await uploadTheme(this.upload.file);
      this.themes.push(theme);
      await this.selectTheme(theme);
      this.upload.open = false;
    } catch (e) {
      this.upload.error.hasError = true;
      this.upload.error.title = localize({ key: 'design.themes.create.error.title' });
      if (e.status === 409) {
        this.upload.error.message = localize({ key: 'design.themes.create.error.conflict' });
      } else {
        this.upload.error.message = localize({ key: 'design.themes.create.error.generic' });
      }
      this.upload.error.hasError = true;
    }
  },
  async updateTheme() {
    try {
      await updateTheme(this.selectedTheme.id, this.edit.file);
      const idx = this.themes.findIndex((t) => t.id === this.selectedTheme.id);
      this.themes[idx] = await getTheme(this.selectedTheme.id);
      this.edit.open = false;
    } catch (e) {
      this.edit.error.hasError = true;
      this.edit.error.title = localize({ key: 'design.themes.edit.error.title' });
      this.edit.error.message = localize({ key: 'design.themes.edit.error.message' });
      this.edit.error.hasError = true;
    }
  },
  upload: {
    open: false,
    file: null,
    error: {
      title: '',
      message: '',
      hasError: false,
      reset() {
        this.title = '';
        this.message = '';
        this.hasError = false;
      },
    },
    selectFile(files) {
      if (files.length >= 1) {
        this.file = files.item(0);
      }
    },
  },
  edit: {
    open: false,
    title: '',
    error: {
      title: '',
      message: '',
      hasError: false,
      reset() {
        this.title = '';
        this.message = '';
        this.hasError = false;
      },
    },
    selectFile(files) {
      if (files.length >= 1) {
        this.file = files.item(0);
      }
    },
  },
  configuration: {
    message: {
      title: '',
      content: '',
      hasError: false,
      hasMessage: false,
    },
  },
  links: {
    picker: {
      label: localize({ key: 'design.themes.links.file_picker_label' }),
    },
  },
}));
