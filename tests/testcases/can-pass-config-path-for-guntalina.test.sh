:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure :bottlebreaker -a -c /etc/guntalina.another.conf --sync
tests:assert-stderr-re "guntalina args: .*-c /etc/guntalina.another.conf.*"
