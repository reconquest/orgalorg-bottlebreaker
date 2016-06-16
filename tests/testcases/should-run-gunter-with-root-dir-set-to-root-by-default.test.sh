:mock:print-args gunter
:mock:true treetrunks
:mock:true guntalina

tests:ensure :bottlebreaker -a
tests:assert-stderr "gunter args: -d / -l /dev/stdout"
