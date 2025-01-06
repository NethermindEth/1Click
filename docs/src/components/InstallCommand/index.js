import React from 'react';
import BashCommand from '../BashScript';

const InstallCommand = () => {
    const [os, setOs] = React.useState('linux-macos');

    React.useEffect(() => {
        const userAgent = navigator.userAgent || navigator.vendor || window.opera;

        if (/windows phone/i.test(userAgent) || /win/i.test(userAgent)) {
            setOs('windows');
        } else if (/android/i.test(userAgent)) {
            setOs('linux-macos');
        } else if (/iPad|iPhone|iPod/.test(userAgent) || /mac/i.test(userAgent)) {
            setOs('linux-macos');
        } else if (/linux/i.test(userAgent)) {
            setOs('linux-macos');
        }
    }, []);

    const commands = {
        'linux-macos': "bash <(curl -fsSL https://github.com/NethermindEth/sedge/raw/main/scripts/install.sh)",
        windows: "Set-ExecutionPolicy Bypass -Scope Process -Force; Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://github.com/NethermindEth/sedge/raw/main/scripts/install.ps1'))",
    };

    return (
        <div>
            <p>To install <strong>Sedge</strong>, run the following command in your terminal:</p>
            <BashCommand command={commands[os]} />
        </div>
    );
};

export default InstallCommand;