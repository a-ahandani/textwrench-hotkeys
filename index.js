const path = require('path');
const os = require('os');

module.exports = {
    getBinaryPath: () => {
        const platform = os.platform();
        if (platform === 'win32') {
            return path.join(__dirname, 'textwrench-hotkeys-windows.exe');
        } else if (platform === 'darwin') {
            return path.join(__dirname, 'textwrench-hotkeys-macos');
        } else {
            throw new Error(`Unsupported platform: ${platform}`);
        }
    }
};
