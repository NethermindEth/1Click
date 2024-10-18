"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9520],{3620:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>a,contentTitle:()=>r,default:()=>h,frontMatter:()=>s,metadata:()=>d,toc:()=>l});var i=t(4848),o=t(8453);const s={sidebar_position:2,slug:"new-networks"},r="Adding new Networks",d={id:"guidelines/new-network",title:"Adding new Networks",description:"We support Ethereum Mainnet and Gnosis Chain, and over time we will update the different testnets on those networks.",source:"@site/docs/guidelines/new-network.mdx",sourceDirName:"guidelines",slug:"/guidelines/new-networks",permalink:"/docs/guidelines/new-networks",draft:!1,unlisted:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/guidelines/new-network.mdx",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2,slug:"new-networks"},sidebar:"tutorialSidebar",previous:{title:"Contributing",permalink:"/docs/guidelines/contributing"},next:{title:"Adding new Clients",permalink:"/docs/guidelines/new-clients"}},a={},l=[];function c(e){const n={code:"code",h1:"h1",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",...(0,o.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.header,{children:(0,i.jsx)(n.h1,{id:"adding-new-networks",children:"Adding new Networks"})}),"\n",(0,i.jsx)(n.p,{children:"We support Ethereum Mainnet and Gnosis Chain, and over time we will update the different testnets on those networks."}),"\n",(0,i.jsx)(n.p,{children:"We support only active testnets, like Sepolia for Ethereum and Chiado for Gnosis, and in the past we supported\nother networks now deprecated, we will give our best to keep our pipeline updated."}),"\n",(0,i.jsx)(n.p,{children:"If you want to contribute to Sedge adding a new network, you can follow the next steps:"}),"\n",(0,i.jsxs)(n.ol,{children:["\n",(0,i.jsxs)(n.li,{children:["Create a folder with the name of the network under ",(0,i.jsx)(n.code,{children:"templates/envs"}),"."]}),"\n",(0,i.jsxs)(n.li,{children:["Create an envs base that contains network base information at ",(0,i.jsx)(n.code,{children:"templates/envs/[network]/env_base.tmpl"}),", like in the\nabove example for sepolia:"]}),"\n"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{children:'{{/* docker-compose_base.tmpl */}}\n{{ define "env" }}\n# --- Global configuration ---\nNETWORK=sepolia{{if .WithMevBoostClient}}\nRELAY_URLS={{.RelayURLs}}{{end}}{{if .FeeRecipient}}\nFEE_RECIPIENT={{.FeeRecipient}}{{end}}\n{{template "execution" .}}\n{{template "consensus" .}}\n{{template "validator" .}}\n{{ end }}\n'})}),"\n",(0,i.jsxs)(n.ol,{start:"3",children:["\n",(0,i.jsxs)(n.li,{children:["Add configs for each of the clients (consensus, execution, validator) in the respective folders inside\n",(0,i.jsx)(n.code,{children:"templates/envs/[network]/[client]"}),", and fill with the needed environment variables."]}),"\n",(0,i.jsxs)(n.li,{children:["Create an entry on ",(0,i.jsx)(n.code,{children:"configs/init.go"}),", in the method ",(0,i.jsx)(n.code,{children:"InitNetworksConfigs()"})," with the network information, for example:"]}),"\n"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{children:'{\n\tNetworkSepolia: {\n\t\tName:               NetworkSepolia,\n\t\tNetworkService:     "merge",\n\t\tGenesisForkVersion: "0x90000069",\n\t\tSupportsMEVBoost:   true,\n\t\tCheckpointSyncURL:  "https://sepolia.checkpoint-sync.ethpandaops.io",\n\t\tRelayURLs: []string{\n\t\t\t"https://0x845bd072b7cd566f02faeb0a4033ce9399e42839ced64e8b2adcfc859ed1e8e1a5a293336a49feac6d9a5edb779be53a@builder-relay-sepolia.flashbots.net",\n\t\t},\n\t},\n},\n'})}),"\n",(0,i.jsxs)(n.ol,{start:"5",children:["\n",(0,i.jsxs)(n.li,{children:["Update documentation, including all the references on ",(0,i.jsx)(n.code,{children:"docs/docs"})," folder, that are going to be displayed on Sedge\nofficial documentation, and on the ",(0,i.jsx)(n.code,{children:"Readme.md"})]}),"\n",(0,i.jsxs)(n.li,{children:["Add entry on the ",(0,i.jsx)(n.code,{children:"CHANGELOG.md"})]}),"\n",(0,i.jsx)(n.li,{children:"Run tests, and add new tests if needed."}),"\n",(0,i.jsx)(n.li,{children:"Create a PR with your changes, and we will review it as soon as possible."}),"\n"]})]})}function h(e={}){const{wrapper:n}={...(0,o.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(c,{...e})}):c(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>r,x:()=>d});var i=t(6540);const o={},s=i.createContext(o);function r(e){const n=i.useContext(s);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function d(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:r(e.components),i.createElement(s.Provider,{value:n},e.children)}}}]);