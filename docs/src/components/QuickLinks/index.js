import React from 'react';
import Link from '@docusaurus/Link';
import styled from '@emotion/styled';
import { ArrowUpRight as LinkIcon } from 'react-feather';

const QuickLinksSection = styled.div`
  display: flex;
  flex-direction: column;
  max-width: 1200px;
  margin: 4rem auto;
  padding: 0 1rem;
`;

const ColumnWrapper = styled.div`
  display: flex;
  flex-direction: row;
  gap: 48px;

  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

const Column = styled.div`
  flex: 1;
`;

const Title = styled.h2`
  font-size: 1.5rem;
  margin-bottom: 1rem;
`;

const Description = styled.p`
  font-size: 1rem;
  color: var(--ifm-color-emphasis-700);
  margin-bottom: 2rem;
`;

const LinkCard = styled(Link)`
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    margin-bottom: 0.5rem;
    border-radius: 8px;
    background-color: var(--ifm-background-surface-color);
    text-decoration: none;
    color: var(--ifm-font-color-base);
    transition: background-color 0.2s;

    &:hover {
        background-color: var(--ifm-color-emphasis-100);
        text-decoration: none; /* Prevents underline on hover */
    }
`;

const LinkTitle = styled.h3`
  font-size: 1rem;
  margin: 0;
`;

const LinkDescription = styled.p`
  font-size: 0.875rem;
  color: var(--ifm-color-emphasis-600);
  margin: 0.25rem 0 0;
`;

const dAppLinks = [
    {
        title: 'Get started',
        description: 'Step-by-step guide to install and configure the necessary tools',
        to: '/docs/quickstart/install-guide',
    },
    {
        title: 'Complete Guide',
        description: 'Follow a comprehensive guide to set up your environment',
        to: '/docs/quickstart/complete-guide',
    },
    {
        title: 'Check Dependencies',
        description: 'Check the dependencies required for your setup',
        to: '/docs/quickstart/dependencies',
    },

    {
        title: 'Keys Management',
        description: 'Manage your validator keys efficiently',
        to: '/docs/quickstart/keys-management',
    },
    {
        title: 'Commands',
        description: 'Learn about advanced commands to manage or set up your node',
        to: '/docs/commands',
    },
];


const smartContractLinks = [
    {
        title: 'Ethereum Mainnet',
        description: 'Deploy and manage your nodes on Ethereum Mainnet',
        to: '/docs/networks/mainnet',
    },
    {
        title: 'Gnosis Network',
        description: 'Get started with the Gnosis Mainnet and Testnet',
        to: '/docs/networks/gnosis',
    },
    {
        title: 'Optimism Network',
        description: 'Deploy on Optimism Mainnet and its test networks',
        to: '/docs/quickstart/running-optimism-node',
    },
    {
        title: 'Exposing API',
        description: 'Guide on exposing APIs for your validator setup',
        to: '/docs/quickstart/samples/exposing-apis',
    },
    {
        title: 'Using Relays',
        description: 'Learn how to utilize relays in your setup',
        to: '/docs/quickstart/samples/using-relays',
    },

];


export function QuickLinks() {
    return (
        <QuickLinksSection>
            <ColumnWrapper>
                <Column>
                    <Title>Integrate with Sedge</Title>
                    <Description>
                        Explore these guided tutorials to get started with Sedge for your Ethereum staking needs.
                    </Description>
                    {dAppLinks.map((link) => (
                        <LinkCard key={link.title} to={link.to}>
                            <div>
                                <LinkTitle>{link.title}</LinkTitle>
                                <LinkDescription>{link.description}</LinkDescription>
                            </div>
                            <LinkIcon size={20} />
                        </LinkCard>
                    ))}
                </Column>
                <Column>
                    <Title>Advanced Sedge Features</Title>
                    <Description>
                        Dive deeper into Sedge's capabilities and optimize your Ethereum staking experience.
                    </Description>
                    {smartContractLinks.map((link) => (
                        <LinkCard key={link.title} to={link.to}>
                            <div>
                                <LinkTitle>{link.title}</LinkTitle>
                                <LinkDescription>{link.description}</LinkDescription>
                            </div>
                            <LinkIcon size={20} />
                        </LinkCard>
                    ))}
                </Column>
            </ColumnWrapper>
        </QuickLinksSection>
    );
}

