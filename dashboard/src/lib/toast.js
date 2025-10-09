// Utility functions to trigger toast notifications from anywhere in the app
import { toast } from './toastStore.js';

// Re-export toast functions for easy importing
export const showToast = toast;

// Convenience functions with common use cases
export const showSuccess = (title, message, caption = '') => {
    toast.success(title, message, caption);
};

export const showError = (title, message, caption = '') => {
    toast.error(title, message, caption);
};

export const showWarning = (title, message, caption = '') => {
    toast.warning(title, message, caption);
};

export const showInfo = (title, message, caption = '') => {
    toast.info(title, message, caption);
};

// Common toast patterns
export const showApiError = (error, caption = '') => {
    const title = 'API Error';
    const message = error?.message || 'An unexpected error occurred';
    toast.error(title, message, caption);
};

export const showApiSuccess = (message, caption = '') => {
    const title = 'Success';
    toast.success(title, message, caption);
};

export const showValidationError = (message, caption = '') => {
    const title = 'Validation Error';
    toast.error(title, message, caption);
};

export const showSaveSuccess = (item = 'Item', caption = '') => {
    const title = 'Saved';
    const message = `${item} has been saved successfully`;
    toast.success(title, message, caption);
};

export const showDeleteSuccess = (item = 'Item', caption = '') => {
    const title = 'Deleted';
    const message = `${item} has been deleted successfully`;
    toast.success(title, message, caption);
};
