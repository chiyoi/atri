#!/bin/zsh
scripts=$0
cd $(dirname $(realpath $scripts)) || return
usage () {
    pwd
    echo "Scripts:"
    echo "$scripts build"
    echo "    Build docker image."
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

build () {
    docker build \
    --build-arg OPENAI_API_KEY=$OPENAI_API_KEY \
    --build-arg TOKEN_DISCORD_APPLICATION=$TOKEN_DISCORD_APPLICATION \
    -t chiyoi/$ARTIFACT .
}

up () {
    docker run -d \
    --restart=on-failure:5 \
    --name=$ARTIFACT \
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
help|build|up|logs|stop|update) 
$@
;;
*)
usage
exit 1
;;
esac
