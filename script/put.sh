curl -X PUT localhost:8080/hello/$1 -d \
'{"dateOfBirth": "'$2'"}' -vvv
