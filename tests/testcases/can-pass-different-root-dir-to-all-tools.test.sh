:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure :bottlebreaker -a -r /different/root --sync
tests:assert-stderr-re "gunter args: .*-d /different/root.*"
tests:assert-stderr-re "treetrunks args: .* /different/root"
