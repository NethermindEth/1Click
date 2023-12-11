"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7782],{3905:(e,n,t)=>{t.d(n,{Zo:()=>c,kt:()=>m});var r=t(7294);function s(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function o(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){s(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,r,s=function(e,n){if(null==e)return{};var t,r,s={},i=Object.keys(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||(s[t]=e[t]);return s}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(s[t]=e[t])}return s}var a=r.createContext({}),u=function(e){var n=r.useContext(a),t=n;return e&&(t="function"==typeof e?e(n):o(o({},n),e)),t},c=function(e){var n=u(e.components);return r.createElement(a.Provider,{value:n},e.children)},p={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},d=r.forwardRef((function(e,n){var t=e.components,s=e.mdxType,i=e.originalType,a=e.parentName,c=l(e,["components","mdxType","originalType","parentName"]),d=u(t),m=s,g=d["".concat(a,".").concat(m)]||d[m]||p[m]||i;return t?r.createElement(g,o(o({ref:n},c),{},{components:t})):r.createElement(g,o({ref:n},c))}));function m(e,n){var t=arguments,s=n&&n.mdxType;if("string"==typeof e||s){var i=t.length,o=new Array(i);o[0]=d;var l={};for(var a in n)hasOwnProperty.call(n,a)&&(l[a]=n[a]);l.originalType=e,l.mdxType="string"==typeof e?e:s,o[1]=l;for(var u=2;u<i;u++)o[u]=t[u];return r.createElement.apply(null,o)}return r.createElement.apply(null,t)}d.displayName="MDXCreateElement"},1588:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>a,contentTitle:()=>o,default:()=>p,frontMatter:()=>i,metadata:()=>l,toc:()=>u});var r=t(7462),s=(t(7294),t(3905));const i={sidebar_position:2,id:"clients"},o="Clients",l={unversionedId:"commands/clients",id:"commands/clients",title:"Clients",description:"Running sedge clients will give you a list of supported clients for consensus, execution and validators.",source:"@site/docs/commands/clients.mdx",sourceDirName:"commands",slug:"/commands/clients",permalink:"/docs/commands/clients",draft:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/commands/clients.mdx",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2,id:"clients"},sidebar:"tutorialSidebar",previous:{title:"Cli",permalink:"/docs/commands/cli"},next:{title:"Deps",permalink:"/docs/commands/deps"}},a={},u=[{value:"Help",id:"help",level:2},{value:"Execution Example",id:"execution-example",level:2}],c={toc:u};function p(e){let{components:n,...t}=e;return(0,s.kt)("wrapper",(0,r.Z)({},c,t,{components:n,mdxType:"MDXLayout"}),(0,s.kt)("h1",{id:"clients"},"Clients"),(0,s.kt)("p",null,"Running ",(0,s.kt)("inlineCode",{parentName:"p"},"sedge clients")," will give you a list of supported clients for consensus, execution and validators."),(0,s.kt)("h2",{id:"help"},"Help"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre"},'$ sedge clients -h\nList supported clients for execution and consensus engines\n\nUsage:\n  sedge clients [flags]\n\nFlags:\n  -h, --help   help for clients\n\nGlobal Flags:\n      --logLevel string   Set Log Level (default "info")\n\n')),(0,s.kt)("h2",{id:"execution-example"},"Execution Example"),(0,s.kt)("p",null,"The execution of ",(0,s.kt)("inlineCode",{parentName:"p"},"sedge clients")," will result in an output like this, that will show an accurate list of supported clients:"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre"},"$ sedge clients\n2023-10-13 14:13:44 -- [INFO] [Logger Init] Log level: info\n2023-10-13 14:13:45 -- [INFO] You are running the latest version of sedge. Version:  v1.3.0\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network chiado\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   -                   teku                teku\n 3   -                   lodestar            lodestar\n\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network custom\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   -                   prysm               prysm\n 3   -                   teku                teku\n 4   -                   lodestar            lodestar\n\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network gnosis\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   erigon              teku                teku\n 3   -                   lodestar            lodestar\n\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network goerli\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   geth                prysm               prysm\n 3   erigon              teku                teku\n 4   besu                lodestar            lodestar\n\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network holesky\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   geth                teku                teku\n 3   erigon              lodestar            lodestar\n 4   besu                prysm               prysm\n\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network mainnet\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   geth                prysm               prysm\n 3   erigon              teku                teku\n 4   besu                lodestar            lodestar\n\n2023-10-13 14:13:45 -- [INFO] Listing supported clients for network sepolia\n\n #   Execution Clients   Consensus Clients   Validator Clients\n=== =================== =================== ===================\n 1   nethermind          lighthouse          lighthouse\n 2   geth                prysm               prysm\n 3   erigon              teku                teku\n 4   besu                lodestar            lodestar\n")))}p.isMDXComponent=!0}}]);