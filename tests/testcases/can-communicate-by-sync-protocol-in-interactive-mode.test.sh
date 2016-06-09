:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure mkfifo "bottlebreaker.stdin"

tests:run-background bottlebreaker \
    tests:pipe :bottlebreaker --sync '<' "bottlebreaker.stdin"

exec {bottlebreaker_stdin}<>"bottlebreaker.stdin"

:in() {
    tests:ensure cat '>&' "$bottlebreaker_stdin"
}

:out() {
    tests:ensure grep -qF "$2" "$(tests:get-background-$1 $bottlebreaker)"
}

:in <<START
ORGALORG:1 HELLO
ORGALORG:1 NODE [a@b.com]
ORGALORG:1 START
START

:out stdout 'ORGALORG:1 SYNC gunter'
:out stderr '0 files changed'

tests:not :out stdout 'ORGALORG:1 SYNC treetrunks'
tests:not :out stderr '0 files removed'

:in <<SYNC
ORGALORG:1 SYNC [a@b.com] gunter
SYNC

:out stdout 'ORGALORG:1 SYNC treetrunks'
:out stderr '0 files removed'

:in <<SYNC
ORGALORG:1 SYNC [a@b.com] treetrunks
SYNC
