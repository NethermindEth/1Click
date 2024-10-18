"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[4720],{9459:(e,n,o)=>{o.r(n),o.d(n,{assets:()=>d,contentTitle:()=>c,default:()=>u,frontMatter:()=>t,metadata:()=>r,toc:()=>l});var i=o(4848),s=o(8453);const t={sidebar_position:10,id:"run"},c="Run",r={id:"commands/run",title:"Run",description:"Running sedge run will run all the services in the docker-compose.yml file under",source:"@site/docs/commands/run.mdx",sourceDirName:"commands",slug:"/commands/run",permalink:"/docs/commands/run",draft:!1,unlisted:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/commands/run.mdx",tags:[],version:"current",sidebarPosition:10,frontMatter:{sidebar_position:10,id:"run"},sidebar:"tutorialSidebar",previous:{title:"Networks",permalink:"/docs/commands/networks"},next:{title:"Slashing Export",permalink:"/docs/commands/slashing-export"}},d={},l=[{value:"Help",id:"help",level:2},{value:"Execution Example",id:"execution-example",level:2},{value:"Running a set of services",id:"running-a-set-of-services",level:3}];function a(e){const n={admonition:"admonition",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",mdxAdmonitionTitle:"mdxAdmonitionTitle",p:"p",pre:"pre",ul:"ul",...(0,s.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.header,{children:(0,i.jsx)(n.h1,{id:"run",children:"Run"})}),"\n",(0,i.jsxs)(n.p,{children:["Running ",(0,i.jsx)(n.code,{children:"sedge run"})," will run all the services in the docker-compose.yml file under\nthe generation folder. To run the services this sequence of actions are made:"]}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Build the necessary images (",(0,i.jsx)(n.code,{children:"docker compose build"}),")"]}),"\n",(0,i.jsxs)(n.li,{children:["Pull the necessary images (",(0,i.jsx)(n.code,{children:"docker compose pull"}),")"]}),"\n",(0,i.jsxs)(n.li,{children:["Create containers (",(0,i.jsx)(n.code,{children:"docker compose create"}),")"]}),"\n",(0,i.jsxs)(n.li,{children:["Start containers (",(0,i.jsx)(n.code,{children:"docker compose up"}),")"]}),"\n"]}),"\n",(0,i.jsxs)(n.admonition,{type:"info",children:[(0,i.jsxs)(n.mdxAdmonitionTitle,{children:["Skip ",(0,i.jsx)(n.code,{children:"docker compose pull"})," action"]}),(0,i.jsxs)(n.p,{children:["To skip the ",(0,i.jsx)(n.code,{children:"docker compose pull"})," command, the ",(0,i.jsx)(n.code,{children:"--skip-pull"})," flag could be used. This is useful\nwhen the images are already pulled or the user wants to use a local custom image."]})]}),"\n",(0,i.jsx)(n.h2,{id:"help",children:"Help"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-shell",children:'$ sedge run -h\nRun all the generated services\n\nUsage:\n  sedge run [flags]\n\nFlags:\n  -h, --help               help for run\n  -p, --path string        generation path for sedge data (default "/path/to/sedge-data")\n      --services strings   List of services to run. If this flag is not provided, all services will run.\n      --skip-pull          Avoid pulling images before running containers. If the images are not available locally, this flag could cause an error.\n\nGlobal Flags:\n      --log-level string   Set Log Level, e.g panic, fatal, error, warn, warning, info, debug, trace (default "info")\n'})}),"\n",(0,i.jsx)(n.h2,{id:"execution-example",children:"Execution Example"}),"\n",(0,i.jsx)(n.admonition,{title:"Docker dependency",type:"info",children:(0,i.jsxs)(n.p,{children:["The ",(0,i.jsx)(n.code,{children:"run"})," command will first check if Docker is installed, if it is not installed\nthen Sedge tries to install it."]})}),"\n",(0,i.jsxs)(n.p,{children:["Once the generated data is ready, services inside the docker-compose.yml file could\nbe run using the ",(0,i.jsx)(n.code,{children:"sedge run"})," command."]}),"\n",(0,i.jsx)(n.admonition,{title:"Generate files first",type:"caution",children:(0,i.jsxs)(n.p,{children:["The ",(0,i.jsx)(n.code,{children:"run"})," command assumes that the generated files are ready to run."]})}),"\n",(0,i.jsxs)(n.p,{children:["First, open a Terminal with access to the sedge binary to get started. Then run\nthe following command to start all services in the default generation path\n",(0,i.jsx)(n.code,{children:"./docker-compose-scripts"}),"."]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-shell",children:"$ sedge run\nUsing config file: /root/.sedge.yaml\n2022-12-29 19:55:55 -- [INFO] [Logger Init] Logging configuration: {Level:debug}\n2022-12-29 19:55:55 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n2022-12-29 20:40:24 -- [INFO] Checking dependencies: docker\n2022-12-29 20:40:24 -- [INFO] All dependencies are installed on host machine\n2022-12-29 19:55:55 -- [INFO] Setting up containers\n# highlight-next-line\n2022-12-29 19:55:55 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml build\n2022-12-29 19:55:55 -- [DEBU] Running command with sudo.\n# highlight-next-line\n2022-12-29 19:55:55 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml pull\n2022-12-29 19:55:55 -- [DEBU] Running command with sudo.\n[+] Running 5/5\n \u283f validator-import Pulled  0.3s\n \u283f validator Pulled         0.2s\n \u283f execution Pulled         0.3s\n \u283f consensus Pulled         0.2s\n \u283f validator-blocker Pulled 0.2s\n# highlight-next-line\n2022-12-29 19:55:56 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml create\n2022-12-29 19:55:56 -- [DEBU] Running command with sudo.\n[+] Running 7/7\n \u283f Network sedge_network                                 Created  0.1s\n \u283f Network docker-compose-scripts_default                Created  0.1s\n \u283f Container execution-client                            Created  0.1s\n \u283f Container docker-compose-scripts-validator-blocker-1  Created  0.0s\n \u283f Container consensus-client                            Created  0.0s\n \u283f Container validator-import-client                     Created  0.0s\n \u283f Container validator-client                            Created  0.0s\n# highlight-next-line\n2022-12-29 19:55:56 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml up -d\n2022-12-29 19:55:56 -- [DEBU] Running command with sudo.\n[+] Running 2/4\n \u283f Container validator-import-client                     Exited   1.2s\n \u283f Container consensus-client                            Waiting  14.4s\n \u283f Container execution-client                            Started  0.9s\n \u283f Container docker-compose-scripts-validator-blocker-1  Waiting\n"})}),"\n",(0,i.jsx)(n.p,{children:"In the logs above, the commands logs are highlighted to better understand the sequence\nof actions performed."}),"\n",(0,i.jsx)(n.h3,{id:"running-a-set-of-services",children:"Running a set of services"}),"\n",(0,i.jsxs)(n.p,{children:["If it is no necessary to run all the services, then a set of services could be specified\nusing the ",(0,i.jsx)(n.code,{children:"--services"})," flag. A good example will be running the execution and consensus\nclient initially without the validator to sync the nodes and after the synchronization\nstart the validator."]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-shell",children:"$ sedge run --services execution,consensus\nUsing config file: /root/.sedge.yaml\n2022-12-29 21:03:29 -- [INFO] [Logger Init] Logging configuration: {Level:debug}\n2022-12-29 21:03:29 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n2022-12-29 21:03:29 -- [INFO] Checking dependencies: docker\n2022-12-29 21:03:29 -- [INFO] All dependencies are installed on host machine\n2022-12-29 21:03:29 -- [INFO] Setting up containers\n# highlight-next-line\n2022-12-29 21:03:29 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml build execution consensus\n2022-12-29 21:03:29 -- [DEBU] Running command with sudo.\n# highlight-next-line\n2022-12-29 21:03:29 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml pull execution consensus\n2022-12-29 21:03:29 -- [DEBU] Running command with sudo.\n[+] Running 2/2\n \u283f execution Pulled 0.2s\n \u283f consensus Pulled 0.4s\n# highlight-next-line\n2022-12-29 21:03:30 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml create execution consensus\n2022-12-29 21:03:30 -- [DEBU] Running command with sudo.\n[+] Running 4/2\n \u283f Network docker-compose-scripts_default  Created  0.1s\n \u283f Network sedge_network                   Created  0.1s\n \u283f Container consensus-client              Created  0.0s\n \u283f Container execution-client              Created  0.1s\n# highlight-next-line\n2022-12-29 21:03:30 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml up -d execution consensus\n2022-12-29 21:03:30 -- [DEBU] Running command with sudo.\n[+] Running 2/2\n \u283f Container consensus-client  Started  0.6s\n \u283f Container execution-client  Started  0.6s\n"})}),"\n",(0,i.jsxs)(n.p,{children:["Now, if you go to the folder with the generated docker compose and run\n",(0,i.jsx)(n.code,{children:"docker compose ps"})," you can check the status of the execution and consensus services:"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-logs",children:'NAME                COMMAND                  SERVICE             STATUS               PORTS\nconsensus-client    "/app/cmd/beacon-cha\u2026"   consensus           running (starting)   0.0.0.0:5054->5054/tcp, :::5054->5054/tcp, 4000-4001/tcp, 0.0.0.0:9000->9000/tcp, :::9000->9000/tcp, 0.0.0.0:9000->9000/udp, :::9000->9000/udp\nexecution-client    "./Nethermind.Runner\u2026"   execution           running              8545/tcp, 0.0.0.0:8008->8008/tcp, :::8008->8008/tcp, 0.0.0.0:30303->30303/tcp, :::30303->30303/tcp, 8551/tcp, 0.0.0.0:30303->30303/udp, :::30303->30303/udp\n'})}),"\n",(0,i.jsx)(n.p,{children:"Notice in this case consensus service is starting, that means is still syncing.\nNow you can start the validator client like follow:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-shell",children:"$ sedge run --services validator\nUsing config file: /root/.sedge.yaml\n2022-12-29 21:24:12 -- [INFO] [Logger Init] Logging configuration: {Level:debug}\n2022-12-29 21:24:12 -- [INFO] You are running the latest version of sedge. Version:  v1.6.0\n2022-12-29 21:24:12 -- [INFO] Checking dependencies: docker\n2022-12-29 21:24:12 -- [INFO] All dependencies are installed on host machine\n2022-12-29 21:24:12 -- [INFO] Setting up containers\n# highlight-next-line\n2022-12-29 21:24:12 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml build validator\n2022-12-29 21:24:12 -- [DEBU] Running command with sudo.\n# highlight-next-line\n2022-12-29 21:24:12 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml pull validator\n2022-12-29 21:24:12 -- [DEBU] Running command with sudo.\n[+] Running 1/1\n \u283f validator Pulled 0.2s\n# highlight-next-line\n2022-12-29 21:24:12 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml create validator\n2022-12-29 21:24:12 -- [DEBU] Running command with sudo.\n[+] Running 4/3\n \u283f Container docker-compose-scripts-validator-blocker-1  Created  0.1s\n \u283f Container validator-import-client                     Created  0.1s\n \u283f Container consensus-client                            Running  0.0s\n \u283f Container validator-client                            Created  0.0s\n# highlight-next-line\n2022-12-29 21:24:13 -- [INFO] Running command: docker compose -f docker-compose-scripts/docker-compose.yml up -d validator\n2022-12-29 21:24:13 -- [DEBU] Running command with sudo.\n[+] Running 1/3\n \u283f Container consensus-client                            Waiting  6.5s\n \u283f Container validator-import-client                     Exited   1.1s\n \u283f Container docker-compose-scripts-validator-blocker-1  Waiting  6.5s\n"})}),"\n",(0,i.jsx)(n.p,{children:"After this, the validator service will wait until consensus syncs to start."})]})}function u(e={}){const{wrapper:n}={...(0,s.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(a,{...e})}):a(e)}},8453:(e,n,o)=>{o.d(n,{R:()=>c,x:()=>r});var i=o(6540);const s={},t=i.createContext(s);function c(e){const n=i.useContext(t);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function r(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:c(e.components),i.createElement(t.Provider,{value:n},e.children)}}}]);