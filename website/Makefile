.PHONY: init build dev

all: build

build: init
	bundle exec middleman build

dev: init
	bundle exec middleman server

init:
ifeq ($(shell gem list \^bundler\$$ -i), false)
	gem install bundler --no-document
endif
	bundle check || bundle
