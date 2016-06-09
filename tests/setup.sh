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

:copy-system-command() {
    ln -sf "$(which $1)" bin/$1
}
