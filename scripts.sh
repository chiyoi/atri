#!/bin/zsh
scripts=$0
cd $(dirname $(realpath $scripts)) || return
usage () {
    pwd
    echo "Scripts:"
    echo "$scripts help"
    echo "    Show this help message."
    echo "$scripts run"
    echo "    Run in develop environment."
    echo "$scripts build"
    echo "    Build image."
    echo "$scripts up"
    echo "    Run in docker."
    echo "$scripts logs"
    echo "    Track container log."
    echo "$scripts stop"
    echo "    Stop and clear running container."
    echo "$scripts update"
    echo "    Pull and Build image, Stop running container and Up the new version."
}

ARTIFACT=atri

help () {
    usage
}

run () {
	export ASSISTANT_ID_ATRI="asst_aTI20AjVAwpCli9qAn7uceNs"
	export CATEGORY="1180855047073054750"
    export DATABASE="neko0001"
    export ENDPOINT_COSMOS="https://neko03cosmos.documents.azure.com:443/"
    go run .
}

build () {
    docker build -t chiyoi/$ARTIFACT .
}

up () {
    docker run -d \
    --restart=on-failure:5 \
    --name=$ARTIFACT \
    -e OPENAI_API_KEY=$OPENAI_API_KEY \
    -e TOKEN_DISCORD_APPLICATION=$TOKEN_DISCORD_APPLICATION \
    -e TOKEN_DISCORD_APPLICATION=$TOKEN_DISCORD_APPLICATION \
    -e ASSISTANT_ID_ATRI="asst_aTI20AjVAwpCli9qAn7uceNs" \
	-e CATEGORY="1180508717167419432" \
    -e DATABASE="atri" \
    -e ENDPOINT_COSMOS="https://neko03cosmos.documents.azure.com:443/" \
    chiyoi/$ARTIFACT
}

logs () {
    docker logs -f $ARTIFACT
}

stop () {
    docker stop $ARTIFACT && docker rm $ARTIFACT
}

update () {
    git pull && build || return
    stop 2>/dev/null
    up
}

case "$1" in
""|-h|-help|--help)
usage
exit
;;
help|run|build|up|logs|stop|update) 
$@
;;
*)
usage
exit 1
;;
esac
