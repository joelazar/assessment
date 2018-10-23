#!/bin/sh
#
# Prerequisites for this script: sh, curl and go
#

start_application() {
  echo Start application
  ./node-todo &>/dev/null &
  bgpid=$!
  sleep 0.5
}

stop_application() {
  echo Stop application
  kill $bgpid
  wait $bgpid &>/dev/null
  sleep 0.5
}

echo Build go application
go get -t
go build

start_application

curl http://127.0.0.1:8000/todos -X POST -H "Content-Type: application/json" -d "{\"text\":\"new task\",\"priority\": 4, \"done\":false}" &>/dev/null
curl http://127.0.0.1:8000/todos -X POST -H "Content-Type: application/json" -d "{\"text\":\"second new task\",\"priority\": 2, \"done\":false}" &>/dev/null
output_before_restart="$(curl -s http://127.0.0.1:8000/todos)"

stop_application
start_application

output_after_restart="$(curl -s http://127.0.0.1:8000/todos)"

stop_application

if [ "$output_before_restart" == "$output_after_restart" ]; then
  echo Test succeed!
else
  echo Test failed!
fi

rm todo-database.json
