tests:not tests:ensure :bottlebreaker --check-deps
tests:assert-stderr "can't find: 'gunter'"
tests:assert-stderr 'executable file not found in $PATH'

tests:ensure ln -s "$(which true)" ./bin/gunter

tests:not tests:ensure :bottlebreaker --check-deps
tests:assert-stderr "can't find: 'guntalina'"
tests:assert-stderr 'executable file not found in $PATH'

tests:ensure ln -s "$(which true)" ./bin/guntalina

tests:not tests:ensure :bottlebreaker --check-deps
tests:assert-stderr "can't find: 'treetrunks'"
tests:assert-stderr 'executable file not found in $PATH'

tests:ensure ln -s "$(which true)" ./bin/treetrunks

tests:ensure :bottlebreaker --check-deps
tests:assert-stderr-empty
