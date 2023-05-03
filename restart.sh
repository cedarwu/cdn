#!/bin/bash
set -e

killall cdn || true
export PORT=8082
nohup ./cdn >> `date +%Y-%m-%d`.log 2>&1 &
