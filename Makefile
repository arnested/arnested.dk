.PHONY: all clean help test FORCE

help: ## Display a list of the public targets.
# Find lines that starts with a word-character, contains a colon and then a
# double hash (underscores are not word-characters, so this excludes private
# targets), then strip the hash and print.
	@grep -E -h "^\w.*:.*##" $(MAKEFILE_LIST) | sed -e 's/\(.*\):.*## *\(.*\)/\1|\2/' | column -s '|' -t | sort

all: $(shell find . -type f -name *.html) ## Tidy all HTML files.

%.html: FORCE
	tidy -qim $@

clean:

test:

FORCE:
