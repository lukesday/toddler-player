<script>
    import { toastStore, removeToast } from './toastStore.js';
    import { ToastNotification, Button } from 'carbon-components-svelte';

    // Subscribe to toast store
    $: toasts = $toastStore;

    function handleClose(toastId) {
        removeToast(toastId);
    }

    function getToastKind(type) {
        switch (type) {
            case 'success':
                return 'success';
            case 'error':
                return 'error';
            case 'warning':
                return 'warning';
            case 'info':
            default:
                return 'info';
        }
    }
</script>

<!-- Toast container positioned at top-right -->
<div class="toast-container">
    {#each toasts as toast (toast.id)}
        <div class="toast-item">
            <ToastNotification
                kind={getToastKind(toast.type)}
                title={toast.title}
                subtitle={toast.message}
                caption={toast.caption}
                timeout={0}
                lowContrast={false}
            />
        </div>
    {/each}
</div>

<style>
    .toast-container {
        position: fixed;
        top: 1rem;
        right: 1rem;
        z-index: 9999;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        max-width: 400px;
    }

    .toast-item {
        animation: slideIn 0.3s ease-out;
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
</style>
