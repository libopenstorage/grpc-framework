TOP=..
OUTPUT=$(TOP)/docs
WEBSITE=./content
include $(TOP)/make/vars.mk

all: copydocs build copyswagger

copydocs:
	rm -rf $(WEBSITE)
	cp -r template $(WEBSITE)
	cp -r $(TOP)/apis $(WEBSITE)
	echo "API Version is $(VERSION)" > $(WEBSITE)/version.md

copyswagger:
	cp -r swagger $(OUTPUT)
	./make-swagger-initializer.sh $(OUTPUT)

venv:
	python3 -m venv venv
	bash -c "source venv/bin/activate && \
			pip3 install --upgrade pip && \
			pip3 install -r requirements.txt"
	@echo "Type: 'source venv/bin/active' to get access to mkdocs"

build: venv
	bash -c "source venv/bin/activate && \
			mkdocs build"

serve: venv
	bash -c "source venv/bin/activate && \
			mkdocs serve"

clean:
	rm -rf $(WEBSITE) venv
