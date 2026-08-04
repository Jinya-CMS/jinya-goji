import { Alpine } from '../lib/alpine.js';

Alpine.data('data', () => ({
  mailing: {
    host: '',
    port: 587,
    username: '',
    password: '',
    encryption: 'tls',
    from: '',
    authRequired: true,
  },
  firstAdmin: {
    artistName: '',
    email: '',
    password: '',
  },
  error: {
    message: '',
    hasError: false,
  },
  state: 'configuration',
  loading: false,
  async saveConfiguration() {
    this.loading = true;
    const data = new FormData();
    for (const key of Object.keys(this.mailing)) {
      data.append(key, this.mailing[key]);
    }
    const result = await fetch('/installer/mailer', {
      method: 'POST',
      body: data,
    });
    if (!result.ok || result.status !== 204) {
      this.error.hasError = true;
      this.error.message = await result.text();
    } else {
      this.state = 'first-admin';
    }
    this.loading = false;
  },
  async createFirstAdmin() {
    this.loading = true;
    const data = new FormData();
    for (const key of Object.keys(this.firstAdmin)) {
      data.append(key, this.firstAdmin[key]);
    }
    const result = await fetch('/installer/admin', {
      method: 'POST',
      body: data,
    });
    if (!result.ok || result.status !== 204) {
      this.error.hasError = true;
      this.error.message = await result.text();
    } else {
      this.state = 'finish';
    }
    this.loading = false;
  },
}));

Alpine.start();
