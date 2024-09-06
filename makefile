init:
	dagger init --sdk=go --source=./dagger
	
test:
	@dagger call test --source=.

build:
	@dagger call build --source=.

publish:
	@dagger call publish --source=.