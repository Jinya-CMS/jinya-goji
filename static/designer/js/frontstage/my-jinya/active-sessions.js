import { Alpine } from '../../../../lib/alpine.js';
import {
  createApiToken,
  deleteApiKey,
  deleteApiToken,
  getApiKeys,
  getApiTokens,
  locateIp,
} from '../../foundation/api/authentication.js';
import UAParser from '../../../lib/uaparser.js';
import localize from '../../foundation/utils/localize.js';
import alert from '../../foundation/ui/alert.js';
import confirm from '../../foundation/ui/confirm.js';

Alpine.data('sessionsData', () => ({
  activeTab: 'logins',
  sessions: [],
  tokens: [],
  currentApiKey: '',
  newToken: {
    name: '',
    expiryDate: null,
    open: false,
  },
  async init() {
    this.sessions = await getApiKeys();
    this.tokens = await getApiTokens();
  },
  getBrowser(userAgent) {
    const parser = new UAParser(userAgent);
    if (parser.getBrowser().name) {
      return localize({
        key: 'my_jinya.devices.browser',
        values: {
          browser: parser.getBrowser().name,
          os: parser.getOS().name,
        },
      });
    }

    return localize({ key: 'my_jinya.devices.unknown_browser' });
  },
  getDevice(userAgent) {
    const parser = new UAParser(userAgent);
    if (parser.getDevice().vendor) {
      return localize({
        key: 'my_jinya.sessions.device',
        values: {
          vendor: parser.getDevice().vendor,
          model: parser.getDevice().model,
        },
      });
    }

    return localize({ key: 'my_jinya.sessions.unknown_device' });
  },
  formatDate(string) {
    const date = new Date(Date.parse(string.toLocaleString()));

    return date.toLocaleString();
  },
  async getLocation(ip) {
    const location = await locateIp(ip);
    if (location?.city && location.city !== 'ZZ') {
      return `${location.city} ${localize({ key: `countries.${location.country}` })}`;
    }

    return localize({ key: 'my_jinya.sessions.unknown' });
  },
  async logout(key) {
    await deleteApiKey(key);
    this.sessions = this.sessions.filter((k) => k.key !== key);
  },
  async deleteToken(token) {
    const tokenToDelete = this.tokens.find((t) => t.token === token);
    if (
      await confirm({
        title: localize({ key: 'my_jinya.sessions.delete_token.title' }),
        message: localize({ key: 'my_jinya.sessions.delete_token.message', values: tokenToDelete }),
        declineLabel: localize({ key: 'my_jinya.sessions.delete_token.decline' }),
        approveLabel: localize({ key: 'my_jinya.sessions.delete_token.approve' }),
        negative: true,
      })
    ) {
      await deleteApiToken(token);
      this.tokens = this.tokens.filter((t) => t.token !== token);
    }
  },
  openNewApiToken() {
    this.newToken.expiryDate = null;
    this.newToken.name = '';
    this.newToken.open = true;
  },
  async createApiToken() {
    const data = {
      name: this.newToken.name,
    };
    if (this.newToken.expiryDate) {
      const date = new Date(Date.parse(this.newToken.expiryDate));
      data.expiryDate = date;
    }

    try {
      const token = await createApiToken(data.name, data.expiryDate);
      this.tokens = await getApiTokens();
      this.newToken.open = false;
      alert({
        title: localize({ key: 'my_jinya.sessions.created.title' }),
        message: localize({ key: 'my_jinya.sessions.created.message', values: token }),
        buttonLabel: localize({ key: 'my_jinya.sessions.created.ok' }),
        positive: true,
      });
    } catch (e) {
      console.error(e);
      alert({
        title: localize({ key: 'my_jinya.sessions.error.title' }),
        message: localize({ key: 'my_jinya.sessions.error.message' }),
        buttonLabel: localize({ key: 'my_jinya.sessions.error.ok' }),
        negative: true,
      });
    }
  },
}));
