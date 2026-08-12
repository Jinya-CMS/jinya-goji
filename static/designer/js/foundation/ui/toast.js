export function dismissToast(id) {
  document.querySelector('cms-toast-container')?.closeToast(id);
}

export function positiveToast(title, content, autoClose = true) {
  return document.querySelector('cms-toast-container')?.openToast(title, content, 'positive', autoClose);
}

export function negativeToast(title, content, autoClose = false) {
  return document.querySelector('cms-toast-container')?.openToast(title, content, 'negative', autoClose);
}
