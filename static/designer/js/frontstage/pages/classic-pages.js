import { Alpine } from '../../../../lib/alpine.js';
import {
  createClassicPage,
  deleteClassicPage,
  getClassicPage,
  getClassicPages,
  updateClassicPage,
} from '../../foundation/api/classic-pages.js';
import localize from '../../foundation/utils/localize.js';
import confirm from '../../foundation/ui/confirm.js';
import { getClassicPageDatabase } from '../../foundation/database/classic-page.js';
import { filterArray } from '../../foundation/utils/search.js';
import { dismissToast, negativeToast, positiveToast } from '../../foundation/ui/toast.js';

import '../../foundation/ui/components/toolbar-editor.js';
import '../../foundation/ui/components/diagrams/sparkline.js';

const classicPageDatabase = getClassicPageDatabase();

Alpine.data('classicPagesData', () => ({
  pages: [],
  selectedPage: null,
  content: null,
  searchTerm: '',
  toastId: null,
  get filteredPages() {
    return filterArray(this.pages, this.searchTerm, ['title', 'content']);
  },
  get formAsJson() {
    return {
      id: this.selectedPage?.id,
      title: this.title,
      content: this.content,
    };
  },
  get title() {
    return `#${this.selectedPage.id} ${this.selectedPage.title}`;
  },
  async init() {
    this.pages = await getClassicPages();
    if (this.pages.length > 0) {
      await this.selectPage(this.pages[0]);
    }

    this.$watch('content', async (value) => await this.savePageContent(value));
    this.$watch('selectedPage', () => dismissToast(this.toastId));
  },
  destroy() {
    dismissToast(this.toastId);
  },
  openCreateDialog() {
    this.create.error.reset();
    this.create.title = '';
    this.create.open = true;
  },
  openEditDialog() {
    this.edit.error.reset();
    this.edit.title = this.selectedPage.title;
    this.edit.open = true;
  },
  async selectPage(page) {
    const savedPage = await this.getSavedPage(page.id);
    this.selectedPage = page;
    this.content = page.content;
    const savedPageUpdatedAt = Date.parse(savedPage?.updated?.at) ?? 0;
    const pageUpdatedAt = Date.parse(page.updated.at);

    if (savedPageUpdatedAt < pageUpdatedAt) {
      await this.deleteSavedPage(page.id);
    } else if (savedPage && savedPage.content !== page.content && savedPage.content) {
      const confirmed = await confirm({
        title: localize({ key: 'pages.classic.load.title' }),
        message: localize({ key: 'pages.classic.load.message' }),
        declineLabel: localize({ key: 'pages.classic.load.decline' }),
        approveLabel: localize({ key: 'pages.classic.load.approve' }),
      });
      if (confirmed) {
        this.content = savedPage.content;
      } else {
        await this.deleteSavedPage(page.id);
      }
    }
  },
  getSavedPage(id) {
    return classicPageDatabase.getChangedPage(id);
  },
  deleteSavedPage(id) {
    return classicPageDatabase.deleteChangedPage(id);
  },
  updateContent(event) {
    this.content = event.target.value;
  },
  async savePageContent() {
    await classicPageDatabase.saveChangedPage(Alpine.raw(this.selectedPage.id), Alpine.raw(this.content));
  },
  async savePage() {
    dismissToast(this.toastId);
    try {
      await updateClassicPage(this.selectedPage.id, this.selectedPage.title, this.content);
      this.toastId = positiveToast({
        title: localize({ key: 'pages.classic.edit.saved.title' }),
        message: localize({ key: 'pages.classic.edit.saved.content' }),
      });
    } catch (e) {
      this.toastId = negativeToast({
        title: localize({ key: 'pages.classic.edit.error.title' }),
        message: localize({ key: 'pages.classic.edit.error.content' }),
      });
    }
  },
  async updatePage() {
    try {
      await updateClassicPage(this.selectedPage.id, this.edit.title, this.selectedPage.content);
      const savedPage = await getClassicPage(this.selectedPage.id);
      this.edit.open = false;
      this.pages[this.pages.indexOf(this.selectedPage)].title = savedPage.title;
      this.pages[this.pages.indexOf(this.selectedPage)].created = savedPage.created;
      this.pages[this.pages.indexOf(this.selectedPage)].updated = savedPage.updated;
      await this.selectPage(savedPage);
    } catch (e) {
      this.edit.error.hasError = true;
      this.edit.error.title = localize({ key: 'pages.classic.edit.error.title' });
      if (e.status === 409) {
        this.edit.error.message = localize({ key: 'pages.classic.edit.error.conflict' });
      } else {
        this.edit.error.message = localize({ key: 'pages.classic.edit.error.generic' });
      }
    }
  },
  async createPage() {
    try {
      const savedPage = await createClassicPage(this.create.title, '');
      this.create.open = false;
      this.pages.push(savedPage);
      await this.selectPage(savedPage);
    } catch (e) {
      this.create.error.hasError = true;
      this.create.error.title = localize({ key: 'pages.classic.create.error.title' });
      if (e.status === 409) {
        this.create.error.message = localize({ key: 'pages.classic.create.error.conflict' });
      } else {
        this.create.error.message = localize({ key: 'pages.classic.create.error.generic' });
      }
    }
  },
  async deletePage() {
    dismissToast(this.toastId);
    const confirmation = await confirm({
      title: localize({ key: 'pages.classic.delete.title' }),
      message: localize({
        key: 'pages.classic.delete.message',
        values: this.selectedPage,
      }),
      declineLabel: localize({ key: 'pages.classic.delete.keep' }),
      approveLabel: localize({ key: 'pages.classic.delete.delete' }),
      negative: true,
    });
    if (confirmation) {
      try {
        await deleteClassicPage(this.selectedPage.id);
        this.pages = this.pages.filter((page) => page.id !== this.selectedPage.id);
        if (this.pages.length > 0) {
          await this.selectPage(this.pages[0]);
        } else {
          this.selectedPage = null;
        }
      } catch (e) {
        let message = '';
        if (e.status === 409) {
          message = localize({ key: 'pages.classic.delete.error.conflict' });
        } else {
          message = localize({ key: 'pages.classic.delete.error.generic' });
        }
        this.toastId = negativeToast({
          title: localize({ key: 'pages.classic.delete.error.title' }),
          message,
        });
      }
    }
  },
  create: {
    open: false,
    title: '',
    error: {
      title: '',
      message: '',
      tagError: '',
      hasError: false,
      reset() {
        this.title = '';
        this.message = '';
        this.hasError = false;
      },
    },
  },
  edit: {
    open: false,
    title: '',
    error: {
      title: '',
      message: '',
      tagError: '',
      hasError: false,
      reset() {
        this.title = '';
        this.message = '';
        this.hasError = false;
      },
    },
  },
}));
