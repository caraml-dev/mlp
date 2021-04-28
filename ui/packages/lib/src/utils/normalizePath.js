export const normalizePath = key => key.replace(/\[([^}\]]+)]/g, ".$1");
