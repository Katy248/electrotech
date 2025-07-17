#!/bin/bash

curl -X POST localhost:1488/api/user/register -d '{"email": "petrovanton247@gmail.com", "password": "password", "first_name" :"Anton", "surname":"Petrov", "last_name" :"V", "phone_number":"123"}' | jq

curl -X POST localhost:1488/api/user/login -d '{"email": "petrovanton247@gmail.com", "password": "password"}' | jq
