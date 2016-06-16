:mock:true gunter
:mock:true guntalina
:mock:true treetrunks

tests:ensure :bottlebreaker -a -v --sync
tests:assert-stderr-re 'running command: \[gunter .*]'
tests:assert-stderr-re 'running command: \[guntalina .*]'
tests:assert-stderr-re 'running command: \[treetrunks .*]'
