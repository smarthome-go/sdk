version := 0.14.0

test:
	go test -v -race .

version:
	python3 update_version.py

gh-release:
	gh release create v$(version) -F ./CHANGELOG.md -t 'SDK v$(version)'
