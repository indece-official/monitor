export const Environment = {
    server: {
        version:    (window as any).CONFIG && (window as any).CONFIG.SERVER_VERSION ? (window as any).CONFIG.SERVER_VERSION : '?'
    },
    setup: {
        enabled:    (window as any).CONFIG && (window as any).CONFIG.SETUP_ENABLED ? (window as any).CONFIG.SETUP_ENABLED : false,
        token:      (window as any).CONFIG && (window as any).CONFIG.SETUP_TOKEN ? (window as any).CONFIG.SETUP_TOKEN : null
    }
};

