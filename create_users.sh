#! /bin/bash

for p in payloads/admin_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/admins < $p; done
for p in payloads/parent_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/parents < $p; done
for p in payloads/teacher_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/teachers < $p; done
#for p in payloads/student_*; do http --auth-type=jwt --auth="$JWT_TOKEN" $APIURL/students < $p; done
