GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	rm -rf ./golangci-lint 


HONY: all
all: ;

.PHONY: pg
pg:
	docker run --rm \
		--name=metrics-db \
		-v $(abspath ./db/init/):/docker-entrypoint-initdb.d \
		-v $(abspath ./db/data/):/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD="P@ssw0rd" \
		-d \
		-p 5432:5432 \
		postgres:15.3

.PHONY: stop-pg
stop-pg:
	docker stop metrics

.PHONY: clean-data
clean-data:
	sudo rm -rf ./db/data/