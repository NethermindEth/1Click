import React from 'react';
import styled from '@emotion/styled';
import ThemedImage from '@theme/ThemedImage';
import useBaseUrl from '@docusaurus/useBaseUrl';

const SectionTitle = styled.h2`
    text-align: left;
    margin-left: 1rem;
    margin-bottom: 2rem;
    font-size: 1.5rem;
`;

const GridContainer = styled.div`
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
`;

const NetworkCard = styled.div`
    display: flex;
    flex-direction: column;
    padding-right: 4rem;
    padding-left: 4rem;
    padding-top: 2rem;
    padding-bottom: 1rem;
    border: 1px solid var(--ifm-color-emphasis-300);
    border-radius: 8px;
    text-align: center;
    background-color: var(--ifm-background-surface-color);
`;

const NetworkLogo = styled(ThemedImage)`
    width: 80px;
    height: 80px;
    object-fit: contain;
    margin: 0 auto 0.5rem;
`;

const NetworkName = styled.h3`
    margin-bottom: 0.5rem;
    font-size: 1.1rem;
    color: var(--ifm-color-emphasis-900);
`;

const NetworkButton = styled.a`
    display: inline-block;
    margin: 0.25rem;
    padding: 0.5rem 1rem;
    background-color: var(--ifm-color-emphasis-200);
    color: var(--ifm-color-emphasis-700);
    border-radius: 4px;
    text-decoration: none;
    font-size: 0.9rem;
`;

const TestnetsList = styled.div`
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 0.5rem;
`;

const TestnetButton = styled(NetworkButton)`
    display: block;
    margin: 0;
`;

export const SupportedNetworks = () => {
    const networks = [
        {
            name: 'Base',
            logo: {
                light: '/img/chains/base-logo.png',
                dark: '/img/chains/base-logo.png',
            },
            networks: [
                { name: 'Mainnet', link: '/docs/quickstart/running-optimism-node#base-support' },
            ]
        },
        {
            name: 'Ethereum',
            logo: {
                light: '/img/chains/eth-logo.svg',
                dark: '/img/chains/eth-logo.svg'
            },
            networks: [
                { name: 'Mainnet', link: '/docs/networks/mainnet' },
            ]
        },
        {
            name: 'Gnosis',
            logo: {
                light: '/img/chains/gno-logo.png',
                dark: '/img/chains/gno-logo.png'
            },
            networks: [
                { name: 'Mainnet', link: '/docs/networks/gnosis' },
            ]
        },
        {
            name: 'Optimism',
            logo: {
                light: '/img/chains/op-logo.png',
                dark: '/img/chains/op-logo.png'
            },
            networks: [
                { name: 'Mainnet', link: '/docs/quickstart/running-optimism-node' }
            ]
        },
        {
            name: 'Testnets',
            networks: [
                { name: 'Holesky', link: '/docs/networks/holesky' },
                { name: 'Sepolia', link: '/docs/networks/sepolia' },
                { name: 'Chiado', link: '/docs/networks/chiado' }
            ]
        },
        {
            name: 'Lido',
            logo: {
                light: '/img/chains/lido-logo.png',
                dark: '/img/chains/lido-logo.png',
            },
            networks: [
                { name: 'CSM', link: '/docs/quickstart/staking-with-lido' }
            ]
        },
    ];

    return (
        <>
            <SectionTitle>Supported Networks</SectionTitle>
            <GridContainer>
                {networks.map((network) => (
                    <NetworkCard key={network.name}>
                        {network.logo && (
                            <NetworkLogo
                                sources={{
                                    light: useBaseUrl(network.logo.light),
                                    dark: useBaseUrl(network.logo.dark),
                                }}
                                alt={`${network.name} logo`}
                            />
                        )}
                        <NetworkName>{network.name}</NetworkName>
                        {network.name === 'Testnets' ? (
                            <TestnetsList>
                                {network.networks.map((net) => (
                                    <TestnetButton key={net.name} href={net.link}>
                                        {net.name}
                                    </TestnetButton>
                                ))}
                            </TestnetsList>
                        ) : (
                            network.networks.map((net) => (
                                <NetworkButton key={net.name} href={net.link}>
                                    {net.name}
                                </NetworkButton>
                            ))
                        )}
                    </NetworkCard>
                ))}
            </GridContainer>
        </>
    );
};

export default SupportedNetworks;