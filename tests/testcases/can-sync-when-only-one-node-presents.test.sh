:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure :bottlebreaker -n --sync <<SYNC
ORGALORG:1 HELLO
ORGALORG:1 NODE [a@b.com]
ORGALORG:1 START
ORGALORG:1 SYNC [a@b.com] gunter
ORGALORG:1 SYNC [a@b.com] treetrunks
SYNC

tests:assert-stdout "ORGALORG:1 SYNC gunter"
tests:assert-stdout "ORGALORG:1 SYNC treetrunks"
