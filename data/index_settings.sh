#!/bin/sh
curl -v -X PATCH \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${MEILI_MASTER_KEY}" \
    -d '{ "sortableAttributes": ["popularity"], "filterableAttributes": ["popularity"] }' \
    http://meilisearch:7700/indexes/actors/settings