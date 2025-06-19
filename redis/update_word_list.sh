#!/bin/bash

WORDLIST_PATH="word_list.txt"
WORDLIST_KEY="p-wordlist"

REDIS_URI="redis://default:${REDIS_PASSWORD}@redis-11110.c257.us-east-1-3.ec2.redns.redis-cloud.com:11110"

redis-cli -u "$REDIS_URI" DEL "$WORDLIST_KEY"

while IFS= read -r word; do
    redis-cli -u "$REDIS_URI" RPUSH "$WORDLIST_KEY" "$word"
done < "$WORDLIST_PATH"