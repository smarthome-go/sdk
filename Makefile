version := 0.13.0

test:
	go test -v -race .

version:
	python3 update_version.py

release: test
	gh release create v$(version) -F ./CHANGELOG.md -t 'SDK v$(version)'
