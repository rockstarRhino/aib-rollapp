#!/bin/bash

# Function to set and log environment variable
set_env_var() {
    local var_name="$1"
    local value="$2"

    export "$var_name"="$value"
    echo "Exported $var_name: $value"
}

# Set and log environment variables
set_env_var ROLLAPP_CHAIN_ID "aib_100-1"
set_env_var KEY_NAME_ROLLAPP "aib-user"
set_env_var DENOM "aaib"
set_env_var MONIKER "$ROLLAPP_CHAIN_ID-sequencer"
