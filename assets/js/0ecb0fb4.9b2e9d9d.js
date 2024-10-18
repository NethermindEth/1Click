"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9285],{6096:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>d,contentTitle:()=>r,default:()=>h,frontMatter:()=>i,metadata:()=>a,toc:()=>l});var o=n(4848),s=n(8453);const i={sidebar_position:6,id:"import-key"},r="Import key",a={id:"commands/import-key",title:"Import key",description:"Running sedge import-key will import validator keys from a directory following",source:"@site/docs/commands/importKey.mdx",sourceDirName:"commands",slug:"/commands/import-key",permalink:"/docs/commands/import-key",draft:!1,unlisted:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/commands/importKey.mdx",tags:[],version:"current",sidebarPosition:6,frontMatter:{sidebar_position:6,id:"import-key"},sidebar:"tutorialSidebar",previous:{title:"Generate",permalink:"/docs/commands/generate"},next:{title:"Keys",permalink:"/docs/commands/keys"}},d={},l=[{value:"Help",id:"help",level:2},{value:"Execution examples",id:"execution-examples",level:2}];function c(e){const t={admonition:"admonition",code:"code",h1:"h1",h2:"h2",header:"header",p:"p",pre:"pre",...(0,s.R)(),...e.components};return(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)(t.header,{children:(0,o.jsx)(t.h1,{id:"import-key",children:"Import key"})}),"\n",(0,o.jsxs)(t.p,{children:["Running ",(0,o.jsx)(t.code,{children:"sedge import-key"})," will import validator keys from a directory following\nthe EIP-2335: BLS12-381 Keystore standard. This command needs to be run on an already\ninitialized sedge setup containing a validator client."]}),"\n",(0,o.jsx)(t.h2,{id:"help",children:"Help"}),"\n",(0,o.jsxs)(t.p,{children:["To know more about the command options, run ",(0,o.jsx)(t.code,{children:"sedge import-key --help"}),":"]}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:'$ sedge import-key --help\n\nImport validator client keys, use the \'from\' flag to specify the keys location,\nand make sure that follows the EIP-2335: BLS12-381 Keystore standard. This command\nassumes that the validator client container exists, stopped or not.\n\nThis command stops the validator client during the importing process due to the\nvalidator database being locked while it\'s running but leaves the validator client\nin the same state in which it was found. That means if the validator is running/stopped\nbefore the import, then the validator will be running/stopped after the command\nis executed, regardless of whether the export fails or not. To force a different\nbehavior use --start-validator and --stop-validator flags.\n\nThe [validator] is a required argument used to specify which validator client, from\nall supported by Sedge (lighthouse, lodestar, prysm or teku), is used to import the\nvalidator keys. This is necessary because each client has its own way to achieve\nthe importation.\n\nUsage:\n  sedge import-key [flags] [validator]\n\nFlags:\n      --container-tag string         Container tag to use. If defined, sedge will add to each container and the network, a suffix with the tag. e.g. sedge-validator-client -> sedge-validator-client-<tag>.\n      --custom-config string         file path or url to use as custom network config.\n      --custom-deploy-block string   custom network deploy block.\n      --custom-genesis string        file path or url to use as custom network genesis.\n      --from string                  path to the validator keys, must follow the EIP-2335: BLS12-381 Keystore standard (default "[WORK_DIR]/sedge-data/keystore")\n  -h, --help                         help for import-key\n  -n, --network string               network (default "mainnet")\n  -p, --path string                  path to the generation directory (default "[WORK_DIR]/sedge-data")\n      --start-validator              starts the validator client after import, regardless of the state the validator was in before\n      --stop-validator               stops the validator client after import, regardless of the state the validator was in before\n\nGlobal Flags:\n      --log-level string   Set Log Level, e.g panic, fatal, error, warn, warning, info, debug, trace (default "info")\n'})}),"\n",(0,o.jsx)(t.h2,{id:"execution-examples",children:"Execution examples"}),"\n",(0,o.jsx)(t.p,{children:"In this example we will import validator keys from a non default directory into\na sedge setup with a stopped validator client (Prysm in this case). This is the\nfolder structure:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:".\n\u251c\u2500\u2500 keystore\n\u2502   \u251c\u2500\u2500 deposit_data.json\n\u2502   \u251c\u2500\u2500 keystore_password.txt\n\u2502   \u2514\u2500\u2500 validator_keys\n\u2502       \u2514\u2500\u2500 keystore-m_12381_3600_0_0_0.json\n\u2514\u2500\u2500 sedge-data\n    \u251c\u2500\u2500 docker-compose.yml\n    \u251c\u2500\u2500 .env\n    \u2514\u2500\u2500 jwtsecret\n"})}),"\n",(0,o.jsxs)(t.p,{children:["The ",(0,o.jsx)(t.code,{children:"keystore"})," folder contains the validator keys, the ",(0,o.jsx)(t.code,{children:"sedge-data"})," folder contains\nthe sedge setup. The ",(0,o.jsx)(t.code,{children:"keystore_password.txt"})," file contains the password to unlock\nthe validator keys."]}),"\n",(0,o.jsx)(t.p,{children:"To import the validator keys, and start the validator client after the import, run:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-shell",children:"$ sedge import-key --from keystore -n sepolia --start-validator prysm\n2023-01-26 11:59:34 -- [INFO] [Logger Init] Log level: info\n2023-01-26 11:59:34 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n# highlight-next-line\n2023-01-26 11:59:34 -- [WARN] The keys path is not the default one, copying the keys to the default path /root/sedge/example/sedge-data/keystore\n2023-01-26 11:59:34 -- [INFO] Importing validator keys\n2023-01-26 11:59:34 -- [INFO] The keys import container is starting\n# highlight-next-line\n2023-01-26 11:59:35 -- [INFO] The validator container is being restarted\n# highlight-next-line\n2023-01-26 11:59:36 -- [INFO] Validator keys imported successfully\n"})}),"\n",(0,o.jsxs)(t.p,{children:["Notice the warning message, this is because the ",(0,o.jsx)(t.code,{children:"--from"})," flag is not the default\npath for the validator keys. Sedge will copy the keys to the default path, and\nthen import them."]}),"\n",(0,o.jsxs)(t.p,{children:["Notice also that the validator client is restarted after the import, this is\nbecause the ",(0,o.jsx)(t.code,{children:"--start-validator"})," flag was used."]}),"\n",(0,o.jsx)(t.p,{children:"The resulted folder structure is:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:".\n\u251c\u2500\u2500 keystore\n\u2502   \u251c\u2500\u2500 deposit_data.json\n\u2502   \u251c\u2500\u2500 keystore_password.txt\n\u2502   \u2514\u2500\u2500 validator_keys\n\u2502       \u2514\u2500\u2500 keystore-m_12381_3600_0_0_0.json\n\u2514\u2500\u2500 sedge-data\n    \u251c\u2500\u2500 docker-compose.yml\n    \u251c\u2500\u2500 .env\n    \u251c\u2500\u2500 jwtsecret\n    \u251c\u2500\u2500 keystore/\n    \u2514\u2500\u2500 validator-data/\n"})}),"\n",(0,o.jsxs)(t.admonition,{type:"note",children:[(0,o.jsx)(t.p,{children:"If you have to import your keys to nimbus, you will need to manually type the password when loading the keys, since the nimbus client does not support non interactive password loading."}),(0,o.jsx)(t.p,{children:"Importing keys on Nimbus is done using the consensus client container instead of the validator client container."})]}),"\n",(0,o.jsxs)(t.p,{children:["Notice the new folder ",(0,o.jsx)(t.code,{children:"keystore"})," inside the ",(0,o.jsx)(t.code,{children:"sedge-data"})," folder, this is where\nthe validator keys are copied to. Also notice the new folder ",(0,o.jsx)(t.code,{children:"validator-data"}),",\nthis is where the validator client stores its data."]})]})}function h(e={}){const{wrapper:t}={...(0,s.R)(),...e.components};return t?(0,o.jsx)(t,{...e,children:(0,o.jsx)(c,{...e})}):c(e)}},8453:(e,t,n)=>{n.d(t,{R:()=>r,x:()=>a});var o=n(6540);const s={},i=o.createContext(s);function r(e){const t=o.useContext(i);return o.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:r(e.components),o.createElement(i.Provider,{value:t},e.children)}}}]);