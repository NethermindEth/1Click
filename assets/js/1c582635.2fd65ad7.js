"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9500],{1376:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>o,default:()=>g,frontMatter:()=>a,metadata:()=>r,toc:()=>h});var i=n(4848),s=n(8453);const a={sidebar_position:12,id:"slashing-import"},o="Slashing Import",r={id:"commands/slashing-import",title:"Slashing Import",description:"Running sedge slashing-import will import slashing protection data to the current validator client. The slashing protection data is a JSON file that meets with the EIP-3076 specification.",source:"@site/docs/commands/slashingImport.mdx",sourceDirName:"commands",slug:"/commands/slashing-import",permalink:"/docs/commands/slashing-import",draft:!1,unlisted:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/commands/slashingImport.mdx",tags:[],version:"current",sidebarPosition:12,frontMatter:{sidebar_position:12,id:"slashing-import"},sidebar:"tutorialSidebar",previous:{title:"Slashing Export",permalink:"/docs/commands/slashing-export"},next:{title:"Version",permalink:"/docs/commands/version"}},l={},h=[{value:"Help",id:"help",level:2},{value:"Execution Example",id:"execution-example",level:2}];function d(e){const t={a:"a",code:"code",h1:"h1",h2:"h2",header:"header",p:"p",pre:"pre",...(0,s.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(t.header,{children:(0,i.jsx)(t.h1,{id:"slashing-import",children:"Slashing Import"})}),"\n",(0,i.jsxs)(t.p,{children:["Running ",(0,i.jsx)(t.code,{children:"sedge slashing-import"})," will import slashing protection data to the current validator client. The slashing protection data is a JSON file that meets with the ",(0,i.jsx)(t.a,{href:"https://eips.ethereum.org/EIPS/eip-3076",children:"EIP-3076"})," specification."]}),"\n",(0,i.jsx)(t.h2,{id:"help",children:"Help"}),"\n",(0,i.jsx)(t.pre,{children:(0,i.jsx)(t.code,{className:"language-shell",children:'$ sedge slashing-import --help\n\nImport Slashing Protection Interchange Format (EIP-3076) data. This command assumes\nthat the validator client container exists, stopped or not and that its database\nis already initialized. The validator database is initialized if the validator is\nrunning or has already run but is stopped, and also after importing the validator keys.\nThis command stops the validator client during the importing process due to the\nvalidator database being locked while it\'s running but leaves the validator client\nin the same state in which it was found. That means if the validator is running/stopped\nbefore the import, then the validator will be running/stopped after the command\nis executed, regardless of whether the export fails or not. To force a different\nbehavior use --start-validator and --stop-validator flags.\nThe [validator] is a required argument used to specify which validator client, from\nall supported by Sedge (lighthouse, lodestar, prysm or teku), is used to import the\nSlashing Protection data. This is necessary because each client has its own way to\nachieve the importation.\n\nUsage:\n  sedge slashing-import [flags] [validator]\n\nExamples:\n\nsedge slashing-import --from slashing-data.json prysm\nsedge slashing-import --from slashing-data.json --stop-validator lodestar\nsedge slashing-import --from slashing-data.json --start-validator lighthouse\n\nFlags:\n      --container-tag string   Container tag to use. If defined, sedge will add to each container and the network, a suffix with the tag. e.g. sedge-validator-client -> sedge-validator-client-<tag>.\n  -f, --from string            path to the JSON file in the EIP-3076 format with the slashing protection data to import (default: <generation-dir>/slashing_protection.json)\n  -h, --help                   help for slashing-import\n  -n, --network string         network (default "mainnet")\n  -p, --path string            path to the generation directory (default "/path/to/sedge-data")\n      --start-validator        starts the validator client after import, regardless of the state the validator was in before\n      --stop-validator         stops the validator client after import, regardless of the state the validator was in before\n\nGlobal Flags:\n      --log-level string   Set Log Level, e.g panic, fatal, error, warn, warning, info, debug, trace (default "info")\n'})}),"\n",(0,i.jsx)(t.h2,{id:"execution-example",children:"Execution Example"}),"\n",(0,i.jsx)(t.p,{children:"This is an example of importing slashing protection data to a setup using sepolia network and prysm as validator client that is already stopped at the moment of the import."}),"\n",(0,i.jsx)(t.pre,{children:(0,i.jsx)(t.code,{className:"language-shell",children:"$ sedge slashing-import prysm -f slashing-export.json -n sepolia\n2023-01-06 14:59:11 -- [INFO] [Logger Init] Log level: info\n2023-01-06 14:59:11 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n# highlight-next-line\n2023-01-06 14:59:11 -- [INFO] Importing slashing data to client prysm from slashing-export.json\n# highlight-next-line\n2023-01-06 14:59:11 -- [INFO] The slashing protection container is starting...\n# highlight-next-line\n2023-01-06 14:59:12 -- [INFO] The slashing container ends successfully.\n"})}),"\n",(0,i.jsxs)(t.p,{children:["Notice in this case the validator client remains stopped because it has been stopped since before the import. If necessary the validator client could be started after the import using the ",(0,i.jsx)(t.code,{children:"--start-validator"})," flag, for example:"]}),"\n",(0,i.jsx)(t.pre,{children:(0,i.jsx)(t.code,{className:"language-shell",children:"$ sedge slashing-import prysm -f slashing-export.json -n sepolia --start-validator\n2023-01-06 15:08:05 -- [INFO] [Logger Init] Log level: info\n2023-01-06 15:08:06 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n2023-01-06 15:08:06 -- [INFO] Importing slashing data to client prysm from slashing-export.json\n2023-01-06 15:08:06 -- [INFO] The slashing protection container is starting...\n2023-01-06 15:08:06 -- [INFO] The slashing container ends successfully.\n# highlight-next-line\n2023-01-06 15:08:06 -- [INFO] The validator container is being restarted\n# highlight-next-line\n2023-01-06 15:08:06 -- [INFO] Validator started.\n"})}),"\n",(0,i.jsx)(t.p,{children:"Another case may be importing the slashing data protection when the validator is currently running, for example:"}),"\n",(0,i.jsx)(t.pre,{children:(0,i.jsx)(t.code,{className:"language-shell",children:"$ sedge slashing-import prysm -f slashing-export.json -n sepolia \n2023-01-06 15:10:27 -- [INFO] [Logger Init] Log level: info\n2023-01-06 15:10:27 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n# highlight-next-line\n2023-01-06 15:10:27 -- [INFO] Stopping validator client...\n# highlight-next-line\n2023-01-06 15:10:27 -- [INFO] stopping service: validator-client, currently on running status\n# highlight-next-line\n2023-01-06 15:10:28 -- [INFO] Validator client stopped.\n2023-01-06 15:10:28 -- [INFO] Importing slashing data to client prysm from slashing-export.json\n2023-01-06 15:10:28 -- [INFO] The slashing protection container is starting...\n2023-01-06 15:10:28 -- [INFO] The slashing container ends successfully.\n# highlight-next-line\n2023-01-06 15:10:28 -- [INFO] The validator container is being restarted\n# highlight-next-line\n2023-01-06 15:10:28 -- [INFO] Validator started.\n"})}),"\n",(0,i.jsxs)(t.p,{children:["In this case, the validator client is stopped before the import and started again afterward. If necessary validator could be stopped after the import using the ",(0,i.jsx)(t.code,{children:"--stop-validator"})," flag. for example:"]}),"\n",(0,i.jsx)(t.pre,{children:(0,i.jsx)(t.code,{className:"language-shell",children:"$ sedge slashing-import prysm -f slashing-export.json -n sepolia --stop-validator\n2023-01-06 15:12:22 -- [INFO] [Logger Init] Log level: info\n2023-01-06 15:12:22 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n# highlight-next-line\n2023-01-06 15:12:22 -- [INFO] Stopping validator client...\n# highlight-next-line\n2023-01-06 15:12:22 -- [INFO] stopping service: validator-client, currently on running status\n# highlight-next-line\n2023-01-06 15:12:22 -- [INFO] Validator client stopped.\n2023-01-06 15:12:22 -- [INFO] Importing slashing data to client prysm from slashing-export.json\n2023-01-06 15:12:22 -- [INFO] The slashing protection container is starting...\n2023-01-06 15:12:23 -- [INFO] The slashing container ends successfully.\n"})}),"\n",(0,i.jsx)(t.p,{children:"In this case, the validator client is stopped before the import but is not started again afterward."})]})}function g(e={}){const{wrapper:t}={...(0,s.R)(),...e.components};return t?(0,i.jsx)(t,{...e,children:(0,i.jsx)(d,{...e})}):d(e)}},8453:(e,t,n)=>{n.d(t,{R:()=>o,x:()=>r});var i=n(6540);const s={},a=i.createContext(s);function o(e){const t=i.useContext(a);return i.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function r(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:o(e.components),i.createElement(a.Provider,{value:t},e.children)}}}]);