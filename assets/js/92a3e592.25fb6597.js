"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[8517],{6737:(e,n,i)=>{i.r(n),i.d(n,{assets:()=>l,contentTitle:()=>o,default:()=>h,frontMatter:()=>t,metadata:()=>d,toc:()=>c});var r=i(4848),s=i(8453);const t={sidebar_position:10,id:"lido-exporter"},o="Lido Exporter",d={id:"quickstart/lido-exporter",title:"Lido Exporter",description:"The Lido Exporter is an independent utility container designed to export metrics from Lido's Community Staking Module (CSM) smart contracts to Prometheus. It is included in Sedge\u2019s Docker Compose stack but can be used in any other stack related to Lido nodes. The exporter is highly flexible and can integrate with external monitoring setups, making it ideal for DevOps pipelines that require insight into Lido\u2019s validator and node operator performance.",source:"@site/docs/quickstart/lido-exporter.mdx",sourceDirName:"quickstart",slug:"/quickstart/lido-exporter",permalink:"/docs/quickstart/lido-exporter",draft:!1,unlisted:!1,editUrl:"https://github.com/NethermindEth/sedge/tree/main/docs/docs/quickstart/lido-exporter.mdx",tags:[],version:"current",sidebarPosition:10,frontMatter:{sidebar_position:10,id:"lido-exporter"},sidebar:"tutorialSidebar",previous:{title:"Staking with Lido using Sedge",permalink:"/docs/quickstart/staking-with-lido"},next:{title:"Running an Optimism Node with Sedge",permalink:"/docs/quickstart/running-optimism-node"}},l={},c=[{value:"Features of the Lido Exporter",id:"features-of-the-lido-exporter",level:2},{value:"Exportable Metrics",id:"exportable-metrics",level:2},{value:"Key Metrics:",id:"key-metrics",level:3},{value:"Accessing Metrics:",id:"accessing-metrics",level:3},{value:"Configuration Options",id:"configuration-options",level:3},{value:"Configuration Settings Table",id:"configuration-settings-table",level:3},{value:"Running the Lido Exporter",id:"running-the-lido-exporter",level:2},{value:"1. <strong>Running with Docker</strong>",id:"1-running-with-docker",level:3},{value:"2. <strong>Running as a CLI Application</strong>",id:"2-running-as-a-cli-application",level:3},{value:"Configuration Precedence",id:"configuration-precedence",level:3},{value:"Using the Lido Exporter in Other DevOps Stacks",id:"using-the-lido-exporter-in-other-devops-stacks",level:2},{value:"Integration Possibilities",id:"integration-possibilities",level:3}];function a(e){const n={code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,s.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(n.header,{children:(0,r.jsx)(n.h1,{id:"lido-exporter",children:"Lido Exporter"})}),"\n",(0,r.jsxs)(n.p,{children:["The ",(0,r.jsx)(n.strong,{children:"Lido Exporter"})," is an independent utility container designed to export metrics from Lido's Community Staking Module (CSM) smart contracts to Prometheus. It is included in Sedge\u2019s Docker Compose stack but can be used in any other stack related to Lido nodes. The exporter is highly flexible and can integrate with external monitoring setups, making it ideal for DevOps pipelines that require insight into Lido\u2019s validator and node operator performance."]}),"\n",(0,r.jsx)(n.h2,{id:"features-of-the-lido-exporter",children:"Features of the Lido Exporter"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Real-Time Smart Contract Event Tracking"}),": Subscribes to Lido CSM smart contract events and converts them into Prometheus metrics. These metrics track critical node operations like penalties, exit requests, and bond status."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Customizable Scraping"}),": Allows users to define how frequently metrics are collected using the ",(0,r.jsx)(n.code,{children:"-scrape-time"})," flag or environment variable (default: 30 seconds)."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Seamless Integration with Prometheus"}),": Provides a ",(0,r.jsx)(n.code,{children:"/metrics"})," endpoint that is compatible with Prometheus for collecting, querying, and alerting based on the exported Lido data."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Versatile Configuration"}),": The Lido Exporter can be configured through environment variables or command-line flags, making it adaptable to different deployment environments."]}),"\n"]}),"\n",(0,r.jsx)(n.h2,{id:"exportable-metrics",children:"Exportable Metrics"}),"\n",(0,r.jsxs)(n.p,{children:["The Lido Exporter collects and exports a comprehensive set of metrics related to Lido CSM operations. These metrics provide valuable insights into validator status, penalties, bonds, exit requests, and rewards. The metrics are labeled by ",(0,r.jsx)(n.strong,{children:"Node Operator ID"})," and ",(0,r.jsx)(n.strong,{children:"network"}),"."]}),"\n",(0,r.jsx)(n.h3,{id:"key-metrics",children:"Key Metrics:"}),"\n",(0,r.jsxs)(n.ol,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Node Operator Metrics"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"nodeOperatorID"}),": The ID of the node operator."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"nodeOperatorManagerAddress"}),": The manager address of the node operator."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"nodeOperatorRewardAddress"}),": The reward address of the node operator."]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Keys Metrics"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"keysStuckValidatorsCount"}),": Number of stuck validators."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"keysRefundedValidatorsCount"}),": Number of refunded validators."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"keysExitedValidatorsCount"}),": Number of validators that exited."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"keysDepositedValidatorsCount"}),": Number of validators deposited."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"keysDepositableValidatorsCount"}),": Number of depositable validators."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"addedKeysCount"}),": Number of keys added."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"withdrawnKeysCount"}),": Number of keys withdrawn."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"vettedKeysCount"}),": Number of vetted keys."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"enqueuedKeysCount"}),": Number of enqueued keys."]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Penalties Metrics"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"penaltiesTotal"}),": Total penalties by type (e.g., ",(0,r.jsx)(n.strong,{children:"EL rewards stealing"}),", ",(0,r.jsx)(n.strong,{children:"initial slashing"}),", ",(0,r.jsx)(n.strong,{children:"withdrawal"}),")."]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Exit Requests Metrics"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"exitRequestsTotal"}),": Total number of exit requests by node operator and network."]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Bond Metrics"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"bondCurrent"}),": The current bond amount for the node operator."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"bondRequired"}),": The required bond amount."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"bondExcess"}),": Excess bond amount."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"bondMissed"}),": Missed bond amount."]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Rewards Metrics"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"nonClaimedRewards"}),": The total amount of rewards that have not yet been claimed."]}),"\n"]}),"\n"]}),"\n"]}),"\n",(0,r.jsx)(n.h3,{id:"accessing-metrics",children:"Accessing Metrics:"}),"\n",(0,r.jsxs)(n.p,{children:["All metrics are exposed via Prometheus at the ",(0,r.jsx)(n.code,{children:"/metrics"})," endpoint, typically available at ",(0,r.jsx)(n.code,{children:"http://localhost:8080/metrics"}),". These metrics can be queried for real-time monitoring or integrated into custom dashboards and alerting solutions."]}),"\n",(0,r.jsx)(n.h3,{id:"configuration-options",children:"Configuration Options"}),"\n",(0,r.jsxs)(n.p,{children:["The ",(0,r.jsx)(n.strong,{children:"Lido Exporter"})," offers flexible configuration through both ",(0,r.jsx)(n.strong,{children:"Environment Variables"})," and ",(0,r.jsx)(n.strong,{children:"Command-Line Flags"}),". These options control key aspects of the exporter, such as which Node Operator to monitor, network selection, and how frequently metrics are collected."]}),"\n",(0,r.jsx)(n.p,{children:"You can configure the Lido Exporter by either:"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["Setting ",(0,r.jsx)(n.strong,{children:"Environment Variables"}),"."]}),"\n",(0,r.jsxs)(n.li,{children:["Using ",(0,r.jsx)(n.strong,{children:"Command-Line Flags"}),"."]}),"\n"]}),"\n",(0,r.jsx)(n.h3,{id:"configuration-settings-table",children:"Configuration Settings Table"}),"\n",(0,r.jsx)(n.p,{children:"The table below outlines the configuration options available for the Lido Exporter, indicating whether they are required or optional, and providing their environment variable and flag equivalents:"}),"\n",(0,r.jsxs)(n.table,{children:[(0,r.jsx)(n.thead,{children:(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.th,{children:"Setting"}),(0,r.jsx)(n.th,{children:"Description"}),(0,r.jsx)(n.th,{children:"Required?"}),(0,r.jsx)(n.th,{children:"Environment Variable"}),(0,r.jsx)(n.th,{children:"Command-Line Flag"})]})}),(0,r.jsxs)(n.tbody,{children:[(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Node Operator ID"})}),(0,r.jsx)(n.td,{children:"The ID of the Lido Node Operator to monitor. This is required for metric collection unless the reward address is provided."}),(0,r.jsxs)(n.td,{children:[(0,r.jsx)(n.strong,{children:"Required"})," (if ",(0,r.jsx)(n.code,{children:"reward-address"})," not provided)"]}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_NODE_OPERATOR_ID"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--node-operator-id"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Reward Address"})}),(0,r.jsx)(n.td,{children:"The reward address of the Node Operator. Used to calculate the Node Operator ID if not explicitly set."}),(0,r.jsxs)(n.td,{children:[(0,r.jsx)(n.strong,{children:"Optional"})," (required if ",(0,r.jsx)(n.code,{children:"node-operator-id"})," is not provided)"]}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_REWARD_ADDRESS"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--reward-address"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Network"})}),(0,r.jsxs)(n.td,{children:["Specifies the target network for monitoring (e.g., ",(0,r.jsx)(n.code,{children:"holesky"}),", ",(0,r.jsx)(n.code,{children:"mainnet"}),")."]}),(0,r.jsxs)(n.td,{children:[(0,r.jsx)(n.strong,{children:"Optional"})," (Default: ",(0,r.jsx)(n.code,{children:"holesky"}),")"]}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_NETWORK"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--network"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"RPC Endpoints"})}),(0,r.jsx)(n.td,{children:"A comma-separated list of Ethereum HTTP RPC endpoints for connecting to the Ethereum network."}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Optional"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_RPC_ENDPOINTS"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--rpc-endpoints"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"WebSocket Endpoints"})}),(0,r.jsx)(n.td,{children:"A comma-separated list of Ethereum WebSocket RPC endpoints for subscribing to events."}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Optional"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_WS_ENDPOINTS"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--ws-endpoints"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Port"})}),(0,r.jsxs)(n.td,{children:["The port on which Prometheus metrics are exposed. Default: ",(0,r.jsx)(n.code,{children:"8080"}),"."]}),(0,r.jsxs)(n.td,{children:[(0,r.jsx)(n.strong,{children:"Optional"})," (Default: ",(0,r.jsx)(n.code,{children:"8080"}),")"]}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_PORT"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--port"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Scrape Time"})}),(0,r.jsxs)(n.td,{children:["The interval at which metrics are collected. Values should be in ",(0,r.jsx)(n.code,{children:"10s"}),", ",(0,r.jsx)(n.code,{children:"1m"}),", ",(0,r.jsx)(n.code,{children:"1h"}),", etc. Default: ",(0,r.jsx)(n.code,{children:"30s"}),"."]}),(0,r.jsxs)(n.td,{children:[(0,r.jsx)(n.strong,{children:"Optional"})," (Default: ",(0,r.jsx)(n.code,{children:"30s"}),")"]}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_SCRAPE_TIME"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--scrape-time"})})]}),(0,r.jsxs)(n.tr,{children:[(0,r.jsx)(n.td,{children:(0,r.jsx)(n.strong,{children:"Log Level"})}),(0,r.jsxs)(n.td,{children:["Sets the verbosity level of logs (",(0,r.jsx)(n.code,{children:"panic"}),", ",(0,r.jsx)(n.code,{children:"fatal"}),", ",(0,r.jsx)(n.code,{children:"error"}),", ",(0,r.jsx)(n.code,{children:"warn"}),", ",(0,r.jsx)(n.code,{children:"info"}),", ",(0,r.jsx)(n.code,{children:"debug"}),", ",(0,r.jsx)(n.code,{children:"trace"}),"). Default: ",(0,r.jsx)(n.code,{children:"info"}),"."]}),(0,r.jsxs)(n.td,{children:[(0,r.jsx)(n.strong,{children:"Optional"})," (Default: ",(0,r.jsx)(n.code,{children:"info"}),")"]}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_LOG_LEVEL"})}),(0,r.jsx)(n.td,{children:(0,r.jsx)(n.code,{children:"--log-level"})})]})]})]}),"\n",(0,r.jsx)(n.h2,{id:"running-the-lido-exporter",children:"Running the Lido Exporter"}),"\n",(0,r.jsx)(n.p,{children:"The Lido Exporter can be run in two main ways:"}),"\n",(0,r.jsxs)(n.h3,{id:"1-running-with-docker",children:["1. ",(0,r.jsx)(n.strong,{children:"Running with Docker"})]}),"\n",(0,r.jsx)(n.p,{children:"You can easily run the Lido Exporter using Docker. There is a published Docker image available, which eliminates the need to build the image yourself."}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsxs)(n.p,{children:[(0,r.jsx)(n.strong,{children:"Step 1"}),": Pull the Docker image:"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"docker pull nethermindeth/lido-exporter:latest\n"})}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsxs)(n.p,{children:[(0,r.jsx)(n.strong,{children:"Step 2"}),": Run the Docker container with the necessary environment variables:"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"docker run -d -p 8080:8080 \\\n  -e LIDO_EXPORTER_NODE_OPERATOR_ID=<your_node_operator_id> \\\n  -e LIDO_EXPORTER_NETWORK=<network_name> \\\n  nethermindeth/lido-exporter:latest\n"})}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["The container listens on port ",(0,r.jsx)(n.code,{children:"8080"})," by default, but you can change this using the ",(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_PORT"})," environment variable."]}),"\n",(0,r.jsxs)(n.li,{children:["The metrics will be available at ",(0,r.jsx)(n.code,{children:"http://localhost:8080/metrics"}),"."]}),"\n"]}),"\n",(0,r.jsx)(n.p,{children:"An example with more optional flags is given below:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"docker run -d -p 9090:9090 \\\n  -e LIDO_EXPORTER_NODE_OPERATOR_ID=12345 \\\n  -e LIDO_EXPORTER_NETWORK=mainnet \\\n  -e LIDO_EXPORTER_RPC_ENDPOINTS=https://mainnet.infura.io/v3/YOUR_INFURA_KEY \\\n  -e LIDO_EXPORTER_WS_ENDPOINTS=wss://mainnet.infura.io/ws/v3/YOUR_INFURA_KEY \\\n  -e LIDO_EXPORTER_PORT=9090 \\\n  -e LIDO_EXPORTER_SCRAPE_TIME=15s \\\n  -e LIDO_EXPORTER_LOG_LEVEL=debug \\\n  nethermindeth/lido-exporter:latest\n"})}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.h3,{id:"2-running-as-a-cli-application",children:["2. ",(0,r.jsx)(n.strong,{children:"Running as a CLI Application"})]}),"\n",(0,r.jsx)(n.p,{children:"You can also run the exporter as a standalone CLI tool:"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsxs)(n.p,{children:[(0,r.jsx)(n.strong,{children:"Step 1"}),": Build the application:"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"cd cmd/lido-exporter\ngo build -o lido-exporter main.go\n"})}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsxs)(n.p,{children:[(0,r.jsx)(n.strong,{children:"Step 2"}),": Run the application with the appropriate flags:"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"./lido-exporter --node-operator-id <your_node_operator_id> --network <network_name>\n"})}),"\n",(0,r.jsx)(n.p,{children:"An example with more optional flags is given below:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"./lido-exporter --node-operator-id 12345 \\\n                --network mainnet \\\n                --rpc-endpoints https://mainnet.infura.io/v3/YOUR_INFURA_KEY \\\n                --ws-endpoints wss://mainnet.infura.io/ws/v3/YOUR_INFURA_KEY \\\n                --port 9090 \\\n                --scrape-time 15s \\\n                --log-level debug\n"})}),"\n"]}),"\n"]}),"\n",(0,r.jsx)(n.h3,{id:"configuration-precedence",children:"Configuration Precedence"}),"\n",(0,r.jsxs)(n.p,{children:["If both ",(0,r.jsx)(n.strong,{children:"Environment Variables"})," and ",(0,r.jsx)(n.strong,{children:"Command-Line Flags"})," are set for the same setting, the ",(0,r.jsx)(n.strong,{children:"Environment Variables"})," will take precedence. For example, if you set ",(0,r.jsx)(n.code,{children:"LIDO_EXPORTER_PORT=9090"})," as an environment variable but also pass ",(0,r.jsx)(n.code,{children:"--port 8080"})," on the command line, the exporter will use ",(0,r.jsx)(n.strong,{children:"9090"}),"."]}),"\n",(0,r.jsx)(n.h2,{id:"using-the-lido-exporter-in-other-devops-stacks",children:"Using the Lido Exporter in Other DevOps Stacks"}),"\n",(0,r.jsxs)(n.p,{children:["The Lido Exporter is highly versatile and can be integrated into other DevOps stacks that involve Lido node operations. Its Prometheus-compatible ",(0,r.jsx)(n.code,{children:"/metrics"})," endpoint allows it to export data to any stack that supports Prometheus or other monitoring tools like Grafana, making it easy to build dashboards or set up alerting for:"]}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsx)(n.li,{children:"Monitoring validator performance in real-time."}),"\n",(0,r.jsx)(n.li,{children:"Tracking penalties, exit requests, or bond statuses."}),"\n",(0,r.jsx)(n.li,{children:"Integrating with other services such as Prometheus Alertmanager for critical notifications."}),"\n"]}),"\n",(0,r.jsx)(n.p,{children:"This allows operators or DevOps teams to leverage Lido Exporter\u2019s metrics in different environments outside of Sedge, providing flexibility for custom tooling, monitoring, and performance analysis."}),"\n",(0,r.jsx)(n.h3,{id:"integration-possibilities",children:"Integration Possibilities"}),"\n",(0,r.jsxs)(n.ol,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Standalone Use in Non-Sedge Stacks"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsx)(n.li,{children:"If you\u2019re managing Lido nodes outside of Sedge, you can deploy the Lido Exporter as a Docker container or CLI tool within your existing stack. Simply point the exporter to your Ethereum RPC/WebSocket endpoints and the desired Lido network (e.g., Holesky, Mainnet)."}),"\n",(0,r.jsxs)(n.li,{children:["The metrics will be exposed in Prometheus format, which can be consumed by external monitoring tools such as ",(0,r.jsx)(n.strong,{children:"Grafana"}),", ",(0,r.jsx)(n.strong,{children:"Prometheus"}),", or even ",(0,r.jsx)(n.strong,{children:"Kubernetes"})," Prometheus Operator."]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Kubernetes Integration"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsx)(n.li,{children:"The exporter can easily be deployed in a Kubernetes environment alongside Prometheus and Grafana. It can act as a source for Prometheus metrics and integrate seamlessly with Grafana dashboards for visualizing Lido node performance."}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"Custom Monitoring and Alerting"}),":","\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["You can also use the Lido Exporter to power custom alerting systems. For example, integrate it with ",(0,r.jsx)(n.strong,{children:"Grafana\u2019s OnCall"})," module or custom alerting pipelines (e.g., using ",(0,r.jsx)(n.strong,{children:"PagerDuty"}),", ",(0,r.jsx)(n.strong,{children:"Slack"}),", or ",(0,r.jsx)(n.strong,{children:"OpsGenie"}),")."]}),"\n"]}),"\n"]}),"\n"]})]})}function h(e={}){const{wrapper:n}={...(0,s.R)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(a,{...e})}):a(e)}},8453:(e,n,i)=>{i.d(n,{R:()=>o,x:()=>d});var r=i(6540);const s={},t=r.createContext(s);function o(e){const n=r.useContext(t);return r.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function d(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:o(e.components),r.createElement(t.Provider,{value:n},e.children)}}}]);