"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[5365],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>m});var a=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function r(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,a,i=function(e,t){if(null==e)return{};var n,a,i={},o=Object.keys(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var u=a.createContext({}),s=function(e){var t=a.useContext(u),n=t;return e&&(n="function"==typeof e?e(t):r(r({},t),e)),n},d=function(e){var t=s(e.components);return a.createElement(u.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},c=a.forwardRef((function(e,t){var n=e.components,i=e.mdxType,o=e.originalType,u=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),c=s(n),m=i,h=c["".concat(u,".").concat(m)]||c[m]||p[m]||o;return n?a.createElement(h,r(r({ref:t},d),{},{components:n})):a.createElement(h,r({ref:t},d))}));function m(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=n.length,r=new Array(o);r[0]=c;var l={};for(var u in t)hasOwnProperty.call(t,u)&&(l[u]=t[u]);l.originalType=e,l.mdxType="string"==typeof e?e:i,r[1]=l;for(var s=2;s<o;s++)r[s]=n[s];return a.createElement.apply(null,r)}return a.createElement.apply(null,n)}c.displayName="MDXCreateElement"},692:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>u,contentTitle:()=>r,default:()=>p,frontMatter:()=>o,metadata:()=>l,toc:()=>s});var a=n(7462),i=(n(7294),n(3905));const o={sidebar_position:6,id:"chiado"},r="Chiado",l={unversionedId:"networks/chiado",id:"networks/chiado",title:"Chiado",description:"Chiado is a Gnosis Chain testnet. Gnosis Chain is an EVM based Ethereum sidechain that is designed to be a platform for",source:"@site/docs/networks/chiado.mdx",sourceDirName:"networks",slug:"/networks/chiado",permalink:"/docs/networks/chiado",draft:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/networks/chiado.mdx",tags:[],version:"current",sidebarPosition:6,frontMatter:{sidebar_position:6,id:"chiado"},sidebar:"tutorialSidebar",previous:{title:"Gnosis",permalink:"/docs/networks/gnosis"},next:{title:"Contribution Guidelines",permalink:"/docs/guidelines"}},u={},s=[{value:"Supported Execution Clients",id:"supported-execution-clients",level:2},{value:"Supported Consensus Clients",id:"supported-consensus-clients",level:2},{value:"Supported Validator Clients",id:"supported-validator-clients",level:2},{value:"Run a Validator or a Full Node",id:"run-a-validator-or-a-full-node",level:2},{value:"Generating a Full Node",id:"generating-a-full-node",level:3},{value:"Generating a Full Node with a Validator",id:"generating-a-full-node-with-a-validator",level:3},{value:"Create keystore for validator",id:"create-keystore-for-validator",level:3},{value:"Running your setup",id:"running-your-setup",level:3},{value:"Consensus Clients Requirements",id:"consensus-clients-requirements",level:2},{value:"Lighthouse Client",id:"lighthouse-client",level:3},{value:"Lighthouse Minimum",id:"lighthouse-minimum",level:4},{value:"Lighthouse Recommended",id:"lighthouse-recommended",level:4},{value:"Teku Client",id:"teku-client",level:3},{value:"Teku Minimum",id:"teku-minimum",level:4},{value:"Teku Recommended",id:"teku-recommended",level:4},{value:"Execution Client Requirements",id:"execution-client-requirements",level:2},{value:"Nethermind Client",id:"nethermind-client",level:3},{value:"Useful data",id:"useful-data",level:2}],d={toc:s};function p(e){let{components:t,...n}=e;return(0,i.kt)("wrapper",(0,a.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"chiado"},"Chiado"),(0,i.kt)("p",null,"Chiado is a Gnosis Chain testnet. Gnosis Chain is an EVM based Ethereum sidechain that is designed to be a platform for\ndecentralized prediction markets."),(0,i.kt)("h2",{id:"supported-execution-clients"},"Supported Execution Clients"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://docs.nethermind.io/"},"Nethermind"))),(0,i.kt)("h2",{id:"supported-consensus-clients"},"Supported Consensus Clients"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://lighthouse-book.sigmaprime.io/"},"Lighthouse")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://chainsafe.github.io/lodestar/"},"Lodestar")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://docs.teku.consensys.net/en/latest/"},"Teku"))),(0,i.kt)("h2",{id:"supported-validator-clients"},"Supported Validator Clients"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://lighthouse-book.sigmaprime.io/"},"Lighthouse")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://chainsafe.github.io/lodestar/"},"Lodestar")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"https://docs.teku.consensys.net/en/latest/"},"Teku"))),(0,i.kt)("h2",{id:"run-a-validator-or-a-full-node"},"Run a Validator or a Full Node"),(0,i.kt)("p",null,"Validators and Full Nodes protect the network by validating transactions and blocks. They are the backbone of the\nnetwork and are responsible for the security of the network. This guide shows you how to use Sedge to set up a full node or a validator on Chiado, a Gnosis testnet."),(0,i.kt)("admonition",{type:"info"},(0,i.kt)("p",{parentName:"admonition"},"Validating the Gnosis Beacon Chain requires 1\n",(0,i.kt)("a",{parentName:"p",href:"https://docs.gnosischain.com/about/tokens/gno"},"GNO")," per validator process.")),(0,i.kt)("h3",{id:"generating-a-full-node"},"Generating a Full Node"),(0,i.kt)("p",null,"To generate a setup of a full node (without a validator node) with random clients, you only need to add the ",(0,i.kt)("inlineCode",{parentName:"p"},"--no-validator")," to ",(0,i.kt)("inlineCode",{parentName:"p"},"sedge generate full-node"),". For example:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"sedge generate full-node --no-validator --network=chiado\n")),(0,i.kt)("h3",{id:"generating-a-full-node-with-a-validator"},"Generating a Full Node with a Validator"),(0,i.kt)("p",null,"To generate a setup of a validator with random clients, you need to omit the ",(0,i.kt)("inlineCode",{parentName:"p"},"--no-validator")," flag. For example:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"sedge generate full-node --network=chiado\n")),(0,i.kt)("h3",{id:"create-keystore-for-validator"},"Create keystore for validator"),(0,i.kt)("p",null,"To create a keystore for a validator, you need to run the following command:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"sedge keys --network chiado\n")),(0,i.kt)("h3",{id:"running-your-setup"},"Running your setup"),(0,i.kt)("p",null,"Once you have generated your setup, you can run it with the following command:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"sedge run\n")),(0,i.kt)("h2",{id:"consensus-clients-requirements"},"Consensus Clients Requirements"),(0,i.kt)("admonition",{type:"tip"},(0,i.kt)("p",{parentName:"admonition"},"The hardware requirements you are about to see are for Gnosis. Chiado is a lighter network, the following requirements can be considered as minimum for a full node in Chiado:"),(0,i.kt)("ul",{parentName:"admonition"},(0,i.kt)("li",{parentName:"ul"},"2 CPU cores"),(0,i.kt)("li",{parentName:"ul"},"4 GB RAM"),(0,i.kt)("li",{parentName:"ul"},"160 GB SDD"))),(0,i.kt)("h3",{id:"lighthouse-client"},"Lighthouse Client"),(0,i.kt)("h4",{id:"lighthouse-minimum"},"Lighthouse Minimum"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Dual-core CPU, 2015 or newer"),(0,i.kt)("li",{parentName:"ul"},"8 GB RAM"),(0,i.kt)("li",{parentName:"ul"},"500 GB SSD"),(0,i.kt)("li",{parentName:"ul"},"10 Mb/s download, 5 Mb/s upload broadband connection")),(0,i.kt)("h4",{id:"lighthouse-recommended"},"Lighthouse Recommended"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Quad-core AMD Ryzen, Intel Broadwell, ARMv8 or newer"),(0,i.kt)("li",{parentName:"ul"},"16 GB RAM"),(0,i.kt)("li",{parentName:"ul"},"1 TB SSD"),(0,i.kt)("li",{parentName:"ul"},"100 Mb/s download, 20 Mb/s upload broadband connection")),(0,i.kt)("h3",{id:"teku-client"},"Teku Client"),(0,i.kt)("h4",{id:"teku-minimum"},"Teku Minimum"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Dual Core CPU, i5-760 or AMD FX-8100 or better"),(0,i.kt)("li",{parentName:"ul"},"8 GB RAM"),(0,i.kt)("li",{parentName:"ul"},"500 GB SSD"),(0,i.kt)("li",{parentName:"ul"},"1 Mb/s broadband connection")),(0,i.kt)("h4",{id:"teku-recommended"},"Teku Recommended"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Quad core CPU, Intel Core i7\u20134770 or AMD FX-8310 or better"),(0,i.kt)("li",{parentName:"ul"},"16 GB RAM"),(0,i.kt)("li",{parentName:"ul"},"1 TB SSD"),(0,i.kt)("li",{parentName:"ul"},"10 Mb/s broadband connection")),(0,i.kt)("admonition",{type:"tip"},(0,i.kt)("p",{parentName:"admonition"},"We recommend using ",(0,i.kt)("inlineCode",{parentName:"p"},"Lighthouse")," or ",(0,i.kt)("inlineCode",{parentName:"p"},"Teku")," while using Sedge for Chiado as they are the most stable and tested clients.")),(0,i.kt)("p",null,"For a more detailed description, you can look at the ",(0,i.kt)("a",{parentName:"p",href:"https://docs.gnosischain.com/node/consensus-layer-validator#beacon-chain-node-requirements"},"Beacon Chain Node Requirements")," in the Gnosis Chain documentation."),(0,i.kt)("h2",{id:"execution-client-requirements"},"Execution Client Requirements"),(0,i.kt)("h3",{id:"nethermind-client"},"Nethermind Client"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"OS: Ubuntu"),(0,i.kt)("li",{parentName:"ul"},"CPU: 2 cores"),(0,i.kt)("li",{parentName:"ul"},"RAM: 8GB"),(0,i.kt)("li",{parentName:"ul"},"Disk: 500gb SSD")),(0,i.kt)("p",null,"For a more detailed description, you can look at the ",(0,i.kt)("a",{parentName:"p",href:"https://docs.gnosischain.com/node/execution-layer-validator"},"Guide for Run a Gnosis Execution Layer Node")),(0,i.kt)("h2",{id:"useful-data"},"Useful data"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},"Native (fee) token: ",(0,i.kt)("inlineCode",{parentName:"p"},"Chiado-xDAI"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},"Staking token: ",(0,i.kt)("inlineCode",{parentName:"p"},"Chiado-GNO"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},"Chain ID: ",(0,i.kt)("inlineCode",{parentName:"p"},"10200"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://docs.gnosischain.com/"},"Gnosis Chain Documentation"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://blockscout.chiadochain.net/"},"Chiado Execution Explorer"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://beacon.chiadochain.net/"},"Chiado Beacon Explorer"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://gnosisfaucet.com/"},"Chiado Faucet"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://rpc.eu-central-2.gateway.fm/v3/gnosis/archival/chiado"},"Execution Layer RPC (Archival)"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://rpc.chiadochain.net"},"Execution Layer RPC"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://ethstats.chiadochain.net/"},"EthStats"))),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"https://forkmon.chiadochain.net/"},"Fork Monitor")))))}p.isMDXComponent=!0}}]);