tests:not tests:ensure :bottlebreaker --check-deps
tests:assert-stderr "unexpected error while checking dependency: 'gunter'"
tests:assert-stderr 'executable file not found in $PATH'

:mock:true gunter

tests:not tests:ensure :bottlebreaker --check-deps
tests:assert-stderr "unexpected error while checking dependency: 'guntalina'"
tests:assert-stderr 'executable file not found in $PATH'

:mock:true guntalina

tests:not tests:ensure :bottlebreaker --check-deps
tests:assert-stderr "unexpected error while checking dependency: 'treetrunks'"
tests:assert-stderr 'executable file not found in $PATH'

:mock:true treetrunks

tests:ensure :bottlebreaker --check-deps
tests:assert-stderr-empty
