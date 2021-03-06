tests:clone orgalorg-bottlebreaker bin/

:bottlebreaker() {
    go-test:run :bottlebreaker:run "${@}"
}

:bottlebreaker:run() {
    PATH=./bin orgalorg-bottlebreaker "${@}"
}

:mock:true() {
    tests:ensure ln -s "$(which true)" bin/$1
}

:mock:template() {
    tests:put bin/$1

    chmod +x bin/$1
}

:mock:print-args() {
    :mock:template "$1" <<MOCK
#!/bin/bash

echo $1 args: "\${@}" >& 2
MOCK
}

:allow-system-command() {
    local path="$(which "$1")"

    if [ ! "$path" ]; then
        printf "[mock] can't allow system command '%s' for use, "\
            "because it's not found" "$1"
        exit 1
    fi

    ln -sf "$path" bin/$1
}
