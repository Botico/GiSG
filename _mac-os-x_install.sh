echo 'Installer made with GISC'
if ! [ -x "$(command -v git)" ]; then
	 xcode-select --install>&2
	exit 1
fi

go build || exit 1
