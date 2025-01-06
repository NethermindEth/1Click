import React, { useState } from 'react';
import { Copy, Check } from 'lucide-react';

const BashCommand = ({ command }) => {
    const [copied, setCopied] = useState(false);

    const copyToClipboard = async () => {
        try {
            await navigator.clipboard.writeText(command);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch (err) {
            console.error('Failed to copy text: ', err);
        }
    };

    return (
        <div className="bash-command-container">
            <div className="bash-command-content">
                <div className="bash-command-code">
                    <pre>
                        <code>{command}</code>
                    </pre>
                </div>
                <div>
                    <button
                        onClick={copyToClipboard}
                        className="bash-command-button"
                        aria-label="Copy to clipboard"
                    >
                        {copied ? (
                            <Check className="bash-command-icon" style={{ color: 'green' }} />
                        ) : (
                            <Copy className="bash-command-icon" style={{ color: 'gray' }} />
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default BashCommand;