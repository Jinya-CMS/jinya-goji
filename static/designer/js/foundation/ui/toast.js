export function dismissToast(id) {
  document.querySelector('cms-toast-container')?.closeToast(id);
}

export function positiveToast({ title, message, autoClose = true }) {
  return document.querySelector('cms-toast-container')?.openToast(title, message, 'positive', autoClose);
}

export function negativeToast({ title, message, autoClose = false }) {
  return document.querySelector('cms-toast-container')?.openToast(title, message, 'negative', autoClose);
}
