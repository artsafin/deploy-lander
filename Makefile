IMAGE_NAME=deploylander

all: build

build:
	docker build -t $(IMAGE_NAME) .

run: build
	docker run --rm -it $(IMAGE_NAME)
