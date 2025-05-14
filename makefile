live/templ:
	templ generate --watch \
		--proxy="http://localhost:8080" \
		--open-browser=false \
		-v

live/server:
	air --build.delay=100 \
		--build.include_ext=go \
		--misc.clean_on_exit=true

live/tailwind:
	./tailwindcss --input="./assets/input.css" \
		--output="./assets/output.css" \
		--minify \
		--watch

live/assets:
	air --build.cmd="templ generate \
		--notify-proxy" \
		--build.bin=true \
		--build.delay=100 \
		--build.include_dir=assets \
		--build.include_ext=css

live:
	make -j4 live/templ live/server live/tailwind live/assets

build:
	go build -tags=production .
