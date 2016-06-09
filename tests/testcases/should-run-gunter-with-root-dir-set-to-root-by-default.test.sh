:mock:template gunter <<GUNTER
#!/bin/bash

echo gunter args: "\${@}" >& 2
GUNTER


:mock:true treetrunks
:mock:true guntalina

tests:ensure :bottlebreaker
tests:assert-stderr "gunter args: -d / -l /dev/stdout"
