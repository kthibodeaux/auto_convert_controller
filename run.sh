#!/bin/bash

git pull
docker compose down
docker compose build --pull
docker compose up
