export:
	@# see https://inkscape.org/en/doc/inkscape-man.html
	rm -f avatar-*.{png,jpg}
	for i in `seq 0 7`; do \
		inkscape avatars.svg --without-gui --export-id=avatar-$$i --export-png=avatar-$$i-tmp.png && \
		gm convert avatar-$$i-tmp.png -shave 3x3 avatar-$$i.jpg && \
		rm avatar-$$i-tmp.png; \
	done

.PHONY: export