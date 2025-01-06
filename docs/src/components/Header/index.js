import React from "react";
import styled from "@emotion/styled";
import ThemedImage from "@theme/ThemedImage";
import useBaseUrl from "@docusaurus/useBaseUrl";
import SearchBar from "@theme/SearchBar";
import InstallCommand from "../InstallCommand";

export function Header() {
    const HeaderWrapper = styled.div`
        width: 100%;
        //background-image: url('/img/background-image.jpg');
        background-size: cover;
        background-position: center;
        position: relative;

        //&::before {
        //    content: '';
        //    position: absolute;
        //    top: 0;
        //    left: 0;
        //    right: 0;
        //    bottom: 0;
        //    background: rgba(0, 0, 0, 0.7); // Darker overall overlay
        //}

    `;


    const HeaderContent = styled.header`
        position: relative;
        padding: 4rem 0;
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        max-width: var(--ifm-container-width);
        margin: 0 auto;
        color: var(--ifm-font-color-base);
    `;

    const Title = styled.h1`
        font-weight: 600;
        font-size: 3rem;
        margin-top: 1rem;
        margin-bottom: 1rem;
    `;

    const SearchWrapper = styled.div`
        margin-top: 2rem;
        width: 100%;
        max-width: 600px;
        display: flex;
        justify-content: center;
    `;

    return (
        <HeaderWrapper>
            <HeaderContent>
                <ThemedImage
                    sources={{
                        dark: useBaseUrl('/img/Sedge_Horizontal_Light.svg'),
                        light: useBaseUrl('/img/Sedge_Horizontal_Dark.svg'),
                    }}
                    alt="Sedge Logo"
                />
                <Title>Easy node setup and deployment tool</Title>
                <InstallCommand/>
                <SearchWrapper>
                    <SearchBar />
                </SearchWrapper>
            </HeaderContent>
        </HeaderWrapper>
    );
}