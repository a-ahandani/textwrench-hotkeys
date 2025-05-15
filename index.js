const os = require('os');
const path = require('path');

function getBinaryPath() {
    const platform = os.platform();
    const isDev = process.env.NODE_ENV === 'development' || process.defaultApp || /[\\/]electron[\\/]/.test(process.execPath);

    if (false) {
        if (platform === 'win32') {
            return path.join(__dirname, 'textwrench-hotkeys.exe');
        } else if (platform === 'darwin') {
            return path.join(__dirname, 'textwrench-hotkeys-macos');
        }
    } else {
        if (platform === 'darwin') {
            // HARDCODE full absolute path for testing
            return '/Users/ahmad/Dev/electron/textwrench-app/dist/mac-arm64/Textwrench.app/Contents/Resources/textwrench-hotkeys-macos';
        }

        if (platform === 'win32') {
            return 'C:\\Path\\To\\Binary\\textwrench-hotkeys.exe'; // adjust as needed
        }
    }

    throw new Error(`Unsupported platform: ${platform}`);
}

module.exports = { getBinaryPath };

