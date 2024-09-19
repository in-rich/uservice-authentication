test:
	bash -c "set -m; GOOGLE_APPLICATION_CREDENTIALS='$(CURDIR)/firebase-key.json' bash '$(CURDIR)/scripts/test.sh'"

run:
	bash -c "set -m; GOOGLE_APPLICATION_CREDENTIALS='$(CURDIR)/firebase-key.json' bash '$(CURDIR)/scripts/run.sh'"

db-console:
	docker exec -it uservice-authentication-postgres-authentication-1 \
		bash -c "PGPASSWORD=postgres psql -U postgres -d postgres"

PHONY: test run
