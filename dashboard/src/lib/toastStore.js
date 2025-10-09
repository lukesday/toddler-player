import { writable } from 'svelte/store';

// Toast store to manage toast notifications
export const toastStore = writable([]);

// Toast actions
export const toast = {
    success: (title, message, caption = '') => {
        addToast('success', title, message, caption);
    },
    error: (title, message, caption = '') => {
        addToast('error', title, message, caption);
    },
    warning: (title, message, caption = '') => {
        addToast('warning', title, message, caption);
    },
    info: (title, message, caption = '') => {
        addToast('info', title, message, caption);
    }
};

function addToast(type, title, message, caption) {
    const id = Math.random().toString(36).substr(2, 9);
    const toastData = {
        id,
        type,
        title,
        message,
        caption,
        timestamp: Date.now()
    };

    toastStore.update(toasts => [...toasts, toastData]);

    // Auto-remove toast after 5 seconds
    setTimeout(() => {
        removeToast(id);
    }, 5000);
}

export function removeToast(id) {
    toastStore.update(toasts => toasts.filter(toast => toast.id !== id));
}
