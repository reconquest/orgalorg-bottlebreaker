:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure :bottlebreaker -a -n --sync
tests:assert-stderr-re "gunter args: .*-r.*"
tests:assert-stderr-re "treetrunks args: .*-n.*"
tests:assert-stderr-re "guntalina args: .*-r.*"
