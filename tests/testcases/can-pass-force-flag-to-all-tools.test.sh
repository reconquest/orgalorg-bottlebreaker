:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure :bottlebreaker -a -f --sync
tests:assert-stderr-re "guntalina args: .*-f.*"
