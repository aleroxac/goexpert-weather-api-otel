#!/usr/bin/env bash


TARGET_API=$1
[ -z "${TARGET_API}" ] && echo "Please inform the target-api: input, orchestrator" && exit 0

case "${TARGET_API}" in
    input) 
        /app/main
    ;;
    orchestrator)
        /app/orchestrator-api
    ;;
    *) 
        echo "Invalid API ${TARGET_API}"
        echo "Please informe the target-api: input, orchestrator"
        exit 0
    ;;
esac
