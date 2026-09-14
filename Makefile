BINARY := curse
PREFIX := /usr/local/bin

.PHONY: build install uninstall clean run

build:
	go build -o $(BINARY) .

install: $(BINARY)
	install -m 755 $(BINARY) $(PREFIX)/$(BINARY)

$(BINARY):
	go build -o $(BINARY) .

uninstall:
	rm -f $(PREFIX)/$(BINARY)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
