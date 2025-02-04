function get-field() {
    # $1 - field name
    cat configs/client_images.yaml | yq "$1"
}

function update-field() {
    # $1 - field name
    # $2 - new value
    echo "Updating $1 to $2 in configs/client_images.yaml"
    yq eval -i "$1 = $2" configs/client_images.yaml
}

function update-client() {
    # $1 client name
    # $2 client type
    # $3 client path
    # $4 latest version

    CURRENT_VERSION=$(get-field "$3.version")
    if [[ $CURRENT_VERSION < $4 ]] ; then
        echo "New version of $1 $2 client is available. Current version: $CURRENT_VERSION, new version: $4"
        update-field "$3.version" "\"$4\""
        echo ""
    fi
}

# Geth
GETH_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/ethereum/go-ethereum/releases/latest | jq -r ".tag_name")
update-client "Geth" "execution" ".execution.geth" "$GETH_LATEST_VERSION"

# Besu
BESU_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/hyperledger/besu/releases/latest | jq -r ".tag_name")
update-client "Besu" "execution" ".execution.besu" "$BESU_LATEST_VERSION"

# Netehrmind
NETHERMIND_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/NethermindEth/nethermind/releases/latest | jq -r ".tag_name")
update-client "Nethermind" "execution" ".execution.nethermind" "$NETHERMIND_LATEST_VERSION"

# Erigon
ERIGON_LATEST=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/ledgerwatch/erigon/releases/latest | jq -r ".tag_name")
update-client "Erigon" "execution" ".execution.erigon" "$ERIGON_LATEST"

# Lighthouse
LIGHTHOUSE_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/sigp/lighthouse/releases/latest | jq -r ".tag_name")
update-client "Lighthouse" "consensus" ".consensus.lighthouse" "$LIGHTHOUSE_LATEST_VERSION"
update-client "Lighthouse" "validator" ".validator.lighthouse" "$LIGHTHOUSE_LATEST_VERSION"

# Lodestar
LODESTAR_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/ChainSafe/lodestar/releases/latest | jq -r ".tag_name")
update-client "Lodestar" "consensus" ".consensus.lodestar" "$LODESTAR_LATEST_VERSION"
update-client "Lodestar" "validator" ".validator.lodestar" "$LODESTAR_LATEST_VERSION"

# Teku
TEKU_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/ConsenSys/teku/releases/latest | jq -r ".tag_name")
update-client "Teku" "consensus" ".consensus.teku" "$TEKU_LATEST_VERSION"
update-client "Teku" "validator" ".validator.teku" "$TEKU_LATEST_VERSION"

# Prysm
PRYSM_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/prysmaticlabs/prysm/releases/latest | jq -r ".tag_name")
update-client "Prysm" "consensus" ".consensus.prysm" "$PRYSM_LATEST_VERSION"
update-client "Prysm" "validator" ".validator.prysm" "$PRYSM_LATEST_VERSION"

# Grandine
GRANDINE_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/grandinetech/grandine/releases | jq -r ".tag_name")
update-client "Grandine" "consensus" ".consensus.grandine" "$GRANDINE_LATEST_VERSION"
update-client "Grandine" "validator" ".validator.grandine" "$GRANDINE_LATEST_VERSION"
## Charon
CHARON_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/ObolNetwork/charon/releases/latest | jq -r ".tag_name")
update-client "Charon" "distributed" ".distributed.charon" "$PRYSM_LATEST_VERSION"

## Charon
CHARON_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/ObolNetwork/charon/releases/latest | jq -r ".tag_name")
update-client "Charon" "distributed" ".distributed.charon" "$PRYSM_LATEST_VERSION"
# Grandine
GRANDINE_LATEST_VERSION=$(curl -H "Authorization: Bearer $PAT" -sL https://api.github.com/repos/grandinetech/grandine/releases | jq -r ".tag_name")
update-client "Grandine" "consensus" ".consensus.grandine" "$GRANDINE_LATEST_VERSION"
update-client "Grandine" "validator" ".validator.grandine" "$GRANDINE_LATEST_VERSION"