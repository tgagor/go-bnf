#!/usr/bin/env bash

# set -x

ACCOUNT_ID="3580eb6044584b28b33e61c515b2ce4a"
PROJECT_NAME="blog"
API_TOKEN="YcR3iAvL6PToroo51E3oBM9NPZqvzZfu1arhWogz"
API_URL="https://api.cloudflare.com/client/v4/accounts/$ACCOUNT_ID/pages/projects/$PROJECT_NAME/deployments"

cutoff=$(date -u -d '30 days ago' +"%Y-%m-%dT%H:%M:%SZ")

page=1
per_page=25
total_pages=1

while [ $page -le $total_pages ]; do
  response=$(curl -s -H "Authorization: Bearer $API_TOKEN" "$API_URL?page=$page&per_page=$per_page")

  # Get total_pages from the first response
  if [ $page -eq 1 ]; then
    total_pages=$(echo "$response" | jq '.result_info.total_pages')
    echo "Total pages: $total_pages"
  fi
  echo "Processing page: $page"
  echo "$response" | jq -c '.result[] | {id: .id, created_on: .created_on}' | while read -r dep; do
    dep_id=$(echo "$dep" | jq -r '.id')
    dep_date=$(echo "$dep" | jq -r '.created_on')
    if [[ "$dep_date" < "$cutoff" ]]; then
      echo "Deleting deployment $dep_id from $dep_date"
      curl -s -X DELETE -H "Authorization: Bearer $API_TOKEN" "$API_URL/$dep_id"
    fi
  done
  page=$((page+1))
done
