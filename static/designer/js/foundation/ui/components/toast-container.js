class ToastContainer extends HTMLElement {
  constructor() {
    super();

    this.root = this.attachShadow({
      mode: 'closed',
    });
  }

  connectedCallback() {
    this.classList.add('jinya-toast__container');
    this.root.innerHTML = `
        <link href="/static/designer/css/toasts.css" rel="stylesheet">
    `;
  }

  openToast(title, content, type, autoClose = true) {
    const uniqueId = crypto.randomUUID();

    const toast = document.createElement('div');
    toast.id = uniqueId;
    toast.classList.add('jinya-toast');
    toast.classList.add(`is--${type}`);
    toast.innerHTML = `
      <button class="jinya-toast__close">
        <svg viewBox="0 0 24 24" class="jinya-toast__icon is--close">
          <path d="M18 6 6 18"/>
          <path d="m6 6 12 12"/>
        </svg>
      </button>
      <h3 class="jinya-toast__title">${title}</h3>
      <div class="jinya-toast__content">${content}</div>
    `;
    toast.querySelector('.jinya-toast__close').addEventListener('click', () => {
      toast.classList.remove('is--open');
      setTimeout(() => {
        toast.remove();
      }, 500);
    });
    this.root.appendChild(toast);
    setTimeout(() => toast.classList.add('is--open'), 1);
    if (autoClose) {
      setTimeout(() => {
        this.closeToast(uniqueId);
      });
    }

    return uniqueId;
  }

  closeToast(id) {
    const toast = this.root.getElementById(id);
    if (!toast) {
      return;
    }

    toast.classList.remove('is--open');
    setTimeout(() => {
      toast.remove();
    }, 500);
  }
}

if (!customElements.get('cms-toast-container')) {
  customElements.define('cms-toast-container', ToastContainer);
}
