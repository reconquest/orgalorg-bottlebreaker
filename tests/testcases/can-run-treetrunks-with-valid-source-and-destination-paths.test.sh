:mock:true gunter
:mock:true guntalina

:mock:print-args treetrunks

tests:ensure :bottlebreaker -a --sync
tests:assert-stderr "treetrunks args: $(pwd) /"
