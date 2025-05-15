const path = require('path');
const os = require('os');

function getBinaryPath() {
    const platform = os.platform();
    const isDev = process.env.NODE_ENV === 'development' || process.defaultApp || /[\\/]electron[\\/]/.test(process.execPath);

    if (isDev) {
        // In development: use local path
        if (platform === 'win32') {
            return path.join(__dirname, 'textwrench-hotkeys.exe');
        } else if (platform === 'darwin') {
            return path.join(__dirname, 'textwrench-hotkeys-macos');
        }
    } else {
        // In production: ensure binary is outside app.asar
        const resourcesPath = process.resourcesPath;

        if (platform === 'win32') {
            return path.join(resourcesPath, 'textwrench-hotkeys.exe');
        } else if (platform === 'darwin') {
            return path.join(resourcesPath, 'textwrench-hotkeys-macos');
        }
    }

    throw new Error(`Unsupported platform: ${platform}`);
}

module.exports = { getBinaryPath };
