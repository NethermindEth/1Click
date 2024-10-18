"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[6579],{1628:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>a,contentTitle:()=>o,default:()=>l,frontMatter:()=>i,metadata:()=>r,toc:()=>h});var s=t(4848),c=t(8453);const i={id:"setting-checkpoint-sync",sidebar_position:8},o="Setting Checkpoint Sync",r={id:"quickstart/samples/setting-checkpoint-sync",title:"Setting Checkpoint Sync",description:"Overview",source:"@site/docs/quickstart/samples/setting-check-sync.mdx",sourceDirName:"quickstart/samples",slug:"/quickstart/samples/setting-checkpoint-sync",permalink:"/docs/quickstart/samples/setting-checkpoint-sync",draft:!1,unlisted:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/quickstart/samples/setting-check-sync.mdx",tags:[],version:"current",sidebarPosition:8,frontMatter:{id:"setting-checkpoint-sync",sidebar_position:8},sidebar:"tutorialSidebar",previous:{title:"Using custom MEV Relays",permalink:"/docs/quickstart/samples/using-relays"},next:{title:"Exposing APIs",permalink:"/docs/quickstart/samples/exposing-apis"}},a={},h=[{value:"Overview",id:"overview",level:2},{value:"Setting Checkpoint Sync",id:"setting-checkpoint-sync-1",level:2}];function d(e){const n={code:"code",h1:"h1",h2:"h2",header:"header",p:"p",pre:"pre",...(0,c.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.header,{children:(0,s.jsx)(n.h1,{id:"setting-checkpoint-sync",children:"Setting Checkpoint Sync"})}),"\n",(0,s.jsx)(n.h2,{id:"overview",children:"Overview"}),"\n",(0,s.jsx)(n.p,{children:'Checkpoint sync, also known as "synchronization with checkpoints" or "checkpoint-based syncing," is a method of\nsynchronizing the blockchain by using predefined checkpoints in the blockchain. This approach is designed to\nspeed up the synchronization process and reduce the amount of data that needs to be downloaded, verified, and stored by\nnodes.'}),"\n",(0,s.jsx)(n.p,{children:"In a traditional full sync, a node downloads and verifies every block and its associated transactions from the genesis\nblock (the first block in the blockchain) to the most recent block. This can be a slow and resource-intensive process,\nespecially as the blockchain grows in size over time."}),"\n",(0,s.jsx)(n.p,{children:"In a checkpoint sync, nodes use predetermined checkpoints in the blockchain history. These checkpoints are block headers\nthat have been agreed upon by the Ethereum community as trustworthy and have been hardcoded into the client software.\nWhen a node performs a checkpoint sync, it downloads and verifies the data starting from the most recent checkpoint,\nrather than the genesis block. This significantly reduces the amount of data the node has to download and process,\nmaking synchronization faster and more efficient."}),"\n",(0,s.jsx)(n.p,{children:"It's important to note that checkpoint syncs rely on trust in the checkpoint data, which could be a potential point of\ncentralization. However, the checkpoints are generally chosen from well-known and verified block headers, reducing the\nrisk of a malicious checkpoint being introduced."}),"\n",(0,s.jsx)(n.p,{children:"Checkpoint sync is one of several synchronization methods that Ethereum clients can use, including full sync, fast sync,\nand snap sync. Each method offers a trade-off between speed, security, and resource usage."}),"\n",(0,s.jsx)(n.h2,{id:"setting-checkpoint-sync-1",children:"Setting Checkpoint Sync"}),"\n",(0,s.jsxs)(n.p,{children:["Each network has set a default checkpoint sync url on Sedge. You can set the checkpoint sync url to a custom url using\nthe ",(0,s.jsx)(n.code,{children:"--checkpoint-sync-url"})," flag. For example, to set the checkpoint sync url to ",(0,s.jsx)(n.code,{children:"https://example.com/checkpoint.json"}),", run:"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{children:"sedge generate full-node --checkpoint-sync-url https://example.com/checkpoint.json\n"})})]})}function l(e={}){const{wrapper:n}={...(0,c.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>o,x:()=>r});var s=t(6540);const c={},i=s.createContext(c);function o(e){const n=s.useContext(i);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function r(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(c):e.components||c:o(e.components),s.createElement(i.Provider,{value:n},e.children)}}}]);