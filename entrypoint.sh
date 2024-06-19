#!/bin/sh

echo "Generate config from env."

# Convert "true" to true, "false" to false, digital string to number
configKVPair=$(jq -n 'env|to_entries[]' | jq '{
    (.key): (
        .value|(
            if . == "true" then true
            elif . == "false" then false
            else (tonumber? // .)
            end
        )
    )
}')

(echo $configKVPair | jq -s add) > config.json

exec ./PT-TrackerProxy
