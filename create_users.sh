#! /bin/bash


APIURL="http://127.0.0.1:3000/api/v1"
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6ImFmODJhYWI5LWNjYjctNDBiYS1hZTI1LTllNzBlMDhhZDc4MyIsImV4cCI6MTU0MjA2NDI0OCwicm9sZSI6ImFkbWluIiwicm9sZV9pZCI6ImFmODJhYWI5LWNjYjctNDBiYS1hZTI1LTllNzBlMDhhZDc4MyJ9.bKnUf3AZrjEgUY7GQKUN8jqmBDWfr9N7yNfkixYVOJI"

for p in payloads/admin_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/admins < $p; done
for p in payloads/parent_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/parents < $p; done
for p in payloads/teacher_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/teachers < $p; done
