#!/bin/bash


for i in {1..4}; do
	json="
	{
		\"name\": \"k8s-$i\",
		\"size\": \"base-2-plus\",
		\"image\": \"ubuntu-20-04-amd64\",
		\"ssh_keys\":[\"a1:83:b2:db:84:c9:48:c5:e9:b8:ee:8e:3a:86:d2:c9\"],
		\"backups\": false
	}";

	curl \
	-X POST \
	-H "Authorization: Bearer $(cat regru_token.txt)" \
	-H "Content-Type: application/json" \
	-d "$json" \
	'https://api.cloudvps.reg.ru/v1/reglets';

done

curl -X GET -H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" 'https://api.cloudvps.reg.ru/v1/reglets' | grep '"ip":' | cut -d':' -f2
