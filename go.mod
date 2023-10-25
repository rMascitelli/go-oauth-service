module github.com/rMascitelli/go-oauth-service/oauth

go 1.21.3

replace github.com/rMascitelli/go-oauth-service/db_connector => ../db_connector

require github.com/rMascitelli/go-oauth-service/db_connector v0.0.0-20231025005322-b7ac6132e06d

require github.com/lib/pq v1.10.9
