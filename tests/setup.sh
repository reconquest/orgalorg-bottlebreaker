tests:clone orgalorg-bottlebreaker bin/

:bottlebreaker() {
    go-test:run :bottlebreaker:run "${@}"
}

:bottlebreaker:run() {
    PATH=./bin orgalorg-bottlebreaker "${@}"
}
