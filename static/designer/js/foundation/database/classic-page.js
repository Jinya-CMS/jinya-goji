import { Dexie } from '../../../lib/dexie.js';
import { getTimestamp } from './utils.js';

class ClassicPageDatabase {
  #database;

  constructor() {
    this.#database = new Dexie('classicPages');
    this.#database.version(1).stores({
      changes: `++id`,
    });
  }

  getChangedPage(id) {
    return this.#database.changes.get(id);
  }

  async deleteChangedPage(id) {
    await this.#database.changes.delete(id);
  }

  async saveChangedPage(id, content) {
    await this.#database.changes.put({
      id,
      content,
      updated: {
        at: getTimestamp(),
      },
    });
  }
}

const classicPageDatabase = new ClassicPageDatabase();

export function getClassicPageDatabase() {
  return classicPageDatabase;
}
