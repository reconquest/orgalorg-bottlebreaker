:mock:print-args gunter
:mock:print-args guntalina
:mock:print-args treetrunks

tests:ensure :bottlebreaker -a -b /var/backup --sync
tests:assert-stderr-re "gunter args: .*-b /var/backup.*"
